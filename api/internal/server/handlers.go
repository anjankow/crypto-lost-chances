package server

import (
	"api/internal/app"
	"errors"
	"html/template"
	"net/http"
	"strconv"

	"go.uber.org/zap"
)

type Results struct {
	Cryptocurrency string
	Income         float32
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

	tmpl, err := template.ParseFiles("static/result.gohtml")
	if err != nil {
		return http.StatusInternalServerError, errors.New("can't create the template: " + err.Error())
	}

	results := Results{Cryptocurrency: "ADA", Income: float32(amount * 2)}
	if err = tmpl.Execute(w, results); err != nil {
		return http.StatusInternalServerError, errors.New("can't execute the template: " + err.Error())
	}

	return http.StatusOK, nil
}

func healthcheck(w http.ResponseWriter, _ *http.Request) {
	w.Write([]byte("all good here"))
}
