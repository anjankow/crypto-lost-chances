package server

import (
	"api/internal/app"
	"errors"
	"net/http"

	"go.uber.org/zap"
)

func calculate(a *app.App, w http.ResponseWriter, r *http.Request) (int, error) {
	method := "POST"
	if r.Method != method {
		return http.StatusMethodNotAllowed, errors.New("incorrect method type: expected: " + method + ", received: " + r.Method)
	}

	if err := r.ParseForm(); err != nil {
		a.Logger.Error(err.Error())
	}
	month := r.PostForm.Get("month")
	amount := r.PostForm.Get("amount")

	a.Logger.Info("calculate request", zap.String("month", month), zap.String("amount", amount))

	return http.StatusOK, nil
}

func healthcheck(w http.ResponseWriter, _ *http.Request) {
	w.Write([]byte("all good here"))
}
