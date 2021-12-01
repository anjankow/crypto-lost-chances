package server

import (
	"lost-chances-calc/internal/app"
	"net/http"
)

func healthcheck(w http.ResponseWriter, _ *http.Request) {
	w.Write([]byte("all good here"))
}

func calculate(a *app.App, w http.ResponseWriter, r *http.Request) (int, error) {

	return 0, nil
}
