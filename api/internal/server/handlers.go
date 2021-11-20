package server

import (
	"api/internal/app"
	"errors"
	"io/ioutil"
	"net/http"
)

func calculate(a *app.App, w http.ResponseWriter, r *http.Request) (int, error) {
	method := "POST"
	if r.Method != method {
		return http.StatusMethodNotAllowed, errors.New("incorrect method type: expected: " + method + ", received: " + r.Method)
	}

	contentType := "application/json"
	if r.Header.Get("Content-Type") != contentType {
		return http.StatusUnsupportedMediaType, errors.New("incorrect content-type: expected: " + contentType + ", received: " + r.Header.Get("Content-Type"))
	}

	_, err := ioutil.ReadAll(r.Body)
	r.Body.Close()
	if err != nil {
		return http.StatusInternalServerError, errors.New("can't read the request body: " + err.Error())
	}

	// var frame app.Frame
	// if err = json.Unmarshal(body, &frame); err != nil {
	// 	return http.StatusBadRequest, errors.New("invalid body: " + err.Error())
	// }

	// if err = a.HandleFrame(frame); err != nil {
	// 	return http.StatusInternalServerError, errors.New("handler error: " + err.Error())
	// }

	return http.StatusOK, nil
}

func healthcheck(w http.ResponseWriter, _ *http.Request) {
	w.Write([]byte("all good here"))
}
