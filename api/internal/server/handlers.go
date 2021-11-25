package server

import (
	"api/internal/app"
	"errors"
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/websocket"
	"go.uber.org/zap"
)

type Results struct {
	Cryptocurrency string
	Income         float32
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type worker struct {
	conn *websocket.Conn
	app  *app.App
}

func (w worker) sendProgressUpdate() error {
	for progressUpdate := 5; ; progressUpdate += 5 {
		time.Sleep(100 * time.Millisecond)
		if err := w.conn.WriteMessage(websocket.TextMessage, []byte(fmt.Sprint(progressUpdate))); err != nil {
			w.app.Logger.Error("writing the progress update failed", zap.Error(err))
			continue
		}
		if progressUpdate >= 100 {
			break
		}
	}
	return nil
}

func progress(a *app.App, w http.ResponseWriter, r *http.Request) (int, error) {
	upgrader.CheckOrigin = func(r *http.Request) bool {
		return true
	}

	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		a.Logger.Warn("upgrade failed", zap.Error(err))
		return http.StatusBadRequest, err
	}

	workerInstance := worker{
		conn: ws,
		app:  a,
	}
	if err = workerInstance.sendProgressUpdate(); err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, nil

}

func calculate(a *app.App, w http.ResponseWriter, r *http.Request) (int, error) {
	method := "POST"
	if r.Method != method {
		return http.StatusMethodNotAllowed, errors.New("incorrect method type: expected: " + method + ", received: " + r.Method)
	}

	if err := r.ParseForm(); err != nil {
		return http.StatusBadRequest, errors.New("can't parse the form: " + err.Error())
	}
	monthStr := r.PostForm.Get("month")
	amountStr := r.PostForm.Get("amount")

	a.Logger.Info("calculate request", zap.String("month", monthStr), zap.String("amount", amountStr))

	if monthStr == "" || amountStr == "" {
		return http.StatusBadRequest, errors.New("parameters 'month' and 'amount' are required")
	}

	amount, err := strconv.Atoi(amountStr)
	if err != nil {
		return http.StatusBadRequest, errors.New("can't parse amount to int: " + err.Error())
	}

	// MAGIC //

	tmpl, err := template.ParseFiles("static/result.html")
	if err != nil {
		return http.StatusInternalServerError, errors.New("can't create the template: " + err.Error())
	}

	// wsAddress := config.GetDomainAddr()
	// a.Logger.Debug("web socket address: " + wsAddress)

	results := Results{Cryptocurrency: "ADA", Income: float32(amount * 2)}
	if err = tmpl.Execute(w, results); err != nil {
		return http.StatusInternalServerError, errors.New("can't execute the template: " + err.Error())
	}

	return http.StatusOK, nil
}

func healthcheck(w http.ResponseWriter, _ *http.Request) {
	w.Write([]byte("all good here"))
}
