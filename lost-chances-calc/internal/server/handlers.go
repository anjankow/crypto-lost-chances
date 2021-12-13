package server

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"lost-chances-calc/internal/app"
	"net/http"
	"time"

	"go.uber.org/multierr"
)

const (
	contentType = "application/json"
)

type calcRequestBody struct {
	RequestID string    `json:"requestID"`
	MonthYear time.Time `json:"monthYear"`
	Amount    int       `json:"amount"`
}

func (c calcRequestBody) validate() (err error) {
	if c.RequestID == "" {
		err = multierr.Append(err, errors.New("missing request ID"))
	}
	if c.Amount == 0 {
		err = multierr.Append(err, errors.New("missing amount"))
	}
	if c.MonthYear.IsZero() {
		err = multierr.Append(err, errors.New("missing date"))
	}
	return err
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

	if err := body.validate(); err != nil {
		return http.StatusBadRequest, errors.New("validation failed: " + err.Error())
	}

	input := app.CalcInput{
		MonthYear:      body.MonthYear,
		FiatInvestment: float64(body.Amount),
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
