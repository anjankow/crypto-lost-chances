package server

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"lost-chances-calc/internal/app"
	"net/http"
	"time"
)

const (
	contentType = "application/json"
)

type calcRequestBody struct {
	RequestID string    `json:"requestID"`
	MonthYear time.Time `json:"monthYear"`
	Amount    int       `json:"amount"`
}

func calculate(a *app.App, w http.ResponseWriter, r *http.Request) (int, error) {
	ctx := r.Context()
	a.Logger.Debug("handling calculate request")

	if r.Header.Get("Content-Type") != contentType {
		return http.StatusUnsupportedMediaType, errors.New("incorrect content-type: expected: " + contentType + ", received: " + r.Header.Get("Content-Type"))
	}

	bodyBytes, err := ioutil.ReadAll(r.Body)
	r.Body.Close()
	if err != nil {
		return http.StatusInternalServerError, errors.New("can't read the request body: " + err.Error())
	}

	var body calcRequestBody
	if err = json.Unmarshal(bodyBytes, &body); err != nil {
		return http.StatusBadRequest, errors.New("invalid body: " + err.Error())
	}

	input := app.CalcInput{
		MonthYear: body.MonthYear,
		Amount:    float64(body.Amount),
	}

	lostChance, err := a.Calculate(ctx, body.RequestID, input)
	if err != nil {
		return http.StatusInternalServerError, errors.New("failed to start the calculation: " + err.Error())
	}

	resp, err := json.Marshal(lostChance)
	if err != nil {
		return http.StatusInternalServerError, errors.New("failed to marshal the response: " + err.Error())
	}

	if _, err := w.Write(resp); err != nil {
		return http.StatusInternalServerError, errors.New("failed to write the response: " + err.Error())
	}

	return http.StatusOK, nil

}

func healthcheck(w http.ResponseWriter, _ *http.Request) {
	w.Write([]byte("all good here"))
}
