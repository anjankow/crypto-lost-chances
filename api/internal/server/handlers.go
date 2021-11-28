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

type UserInput struct {
	MonthYear time.Time
	Amount    int
}

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

func getUserInput(a *app.App, r *http.Request) (input UserInput, err error) {
	if err := r.ParseForm(); err != nil {
		errors.New("can't parse the form: " + err.Error())
	}
	dateStr := r.PostForm.Get("month")
	amountStr := r.PostForm.Get("amount")

	if dateStr == "" || amountStr == "" {
		err = errors.New("parameters 'month' and 'amount' are required")
		return
	}
	a.Logger.Info("user input", zap.String("date", dateStr), zap.String("amount", amountStr))

	layout := "2006-01-02"
	date, err := time.Parse(layout, dateStr+"-01")
	if err != nil {
		err = errors.New("can't parse 'month' parameter: " + err.Error())
		return
	}

	amount, err := strconv.Atoi(amountStr)
	if err != nil {
		err = errors.New("can't parse amount to int: " + err.Error())
		return
	}

	input.Amount = amount
	input.MonthYear = date

	return
}

func handleCalculate(a *app.App, w http.ResponseWriter, r *http.Request) (int, error) {
	method := "POST"
	if r.Method != method {
		return http.StatusMethodNotAllowed, errors.New("incorrect method type: expected: " + method + ", received: " + r.Method)
	}

	userInput, err := getUserInput(a, r)
	if err != nil {
		return http.StatusBadRequest, err
	}
	a.Logger.Info("calculate request", zap.String("month", userInput.MonthYear.Month().String()), zap.Int("month", userInput.MonthYear.Year()), zap.Int("amount", userInput.Amount))

	// MAGIC //

	tmpl, err := template.ParseFiles("static/result.html")
	if err != nil {
		return http.StatusInternalServerError, errors.New("can't create the template: " + err.Error())
	}

	results := Results{Cryptocurrency: "ADA", Income: float32(userInput.Amount * 2)}
	if err = tmpl.Execute(w, results); err != nil {
		return http.StatusInternalServerError, errors.New("can't execute the template: " + err.Error())
	}

	return http.StatusOK, nil
}

func healthcheck(w http.ResponseWriter, _ *http.Request) {
	w.Write([]byte("all good here"))
}
