package server

import (
	"api/internal/app"
	"api/internal/server/middleware"
	"context"
	"errors"
	"fmt"
	"html/template"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/gorilla/websocket"
	"go.uber.org/zap"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type worker struct {
	conn *websocket.Conn
	app  *app.App
}

type ResultsTemplate struct {
	Cryptocurrency string
	Income         float32
	RequestID      string
}

func (w worker) sendProgressUpdate(ctx context.Context, requestID string) error {

	callback := func(progress int) {
		w.app.Logger.Debug("progress update", zap.Int("progress", progress))
		if err := w.conn.WriteMessage(websocket.TextMessage, []byte(fmt.Sprint(progress))); err != nil {
			w.app.Logger.Error("writing the progress update failed", zap.Error(err))
		}
	}

	w.app.ListenProgress(ctx, requestID, callback)

	w.app.Logger.Debug("end of progress updates")

	return nil
}

func progress(a *app.App, w http.ResponseWriter, r *http.Request) (status int, err error) {

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

	status = http.StatusInternalServerError

	// get request ID
	parsed, err := url.ParseQuery(r.URL.RawQuery)
	if err != nil {
		return
	}
	requestID := parsed["id"][0]
	a.Logger.Debug("displaying progress updates for request " + requestID)

	if err = workerInstance.sendProgressUpdate(r.Context(), requestID); err != nil {
		return
	}

	return http.StatusOK, nil

}

func getUserInput(a *app.App, r *http.Request) (input app.UserInput, err error) {
	if err = r.ParseForm(); err != nil {
		err = errors.New("can't parse the form: " + err.Error())
		return
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
	ctx := r.Context()
	requestID := middleware.GetRequestID(ctx)
	a.Logger.Debug("handling calculate request", zap.String("id", requestID))

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
	results, err := a.StartCalculation(r.Context(), userInput)
	if err != nil {
		return http.StatusInternalServerError, errors.New("can't process the request: " + err.Error())
	}

	tmpl, err := template.ParseFiles("static/result.html")
	if err != nil {
		return http.StatusInternalServerError, errors.New("can't create the template: " + err.Error())
	}

	if err = tmpl.Execute(w, makeResultTemplateValues(requestID, results)); err != nil {
		return http.StatusInternalServerError, errors.New("can't execute the template: " + err.Error())
	}

	return http.StatusOK, nil
}

func makeResultTemplateValues(requestID string, results app.Results) ResultsTemplate {
	return ResultsTemplate{
		Cryptocurrency: results.Cryptocurrency,
		Income:         results.Income,
		RequestID:      requestID,
	}
}

func healthcheck(w http.ResponseWriter, _ *http.Request) {
	w.Write([]byte("all good here"))
}
