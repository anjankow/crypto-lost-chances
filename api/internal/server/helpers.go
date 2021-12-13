package server

import (
	"errors"
	"net/http"
	"net/url"

	"github.com/gorilla/websocket"
)

func upgradeConnection(w http.ResponseWriter, r *http.Request) (conn *websocket.Conn, err error) {
	upgrader.CheckOrigin = func(r *http.Request) bool {
		return true
	}

	conn, err = upgrader.Upgrade(w, r, nil)
	if err != nil {
		err = errors.New("upgrade failed: " + err.Error())
		return
	}

	return
}

func getRequestIdFromQuery(r *http.Request) (requestID string, err error) {

	parsed, err := url.ParseQuery(r.URL.RawQuery)
	if err != nil {
		return
	}

	requestIDlist, ok := parsed["id"]
	if !ok || len(requestIDlist) == 0 || requestIDlist[0] == "" {
		return requestID, errors.New("missing request ID")
	}

	requestID = requestIDlist[0]
	return
}
