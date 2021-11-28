package server

import (
	"api/internal/app"
	"api/internal/config"
	"api/internal/server/middleware"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"

	"go.uber.org/zap"
)

type server struct {
	logger     *zap.Logger
	app        *app.App
	httpServer *http.Server
	addr       string
}

type appHandler struct {
	app    *app.App
	Handle AppHandleFunc
}

type AppHandleFunc func(*app.App, http.ResponseWriter, *http.Request) (int, error)

func (ser server) registerHandlers(router *mux.Router) {
	router.PathPrefix("/calculate").Methods("POST").Handler(middleware.AddRequestID(appHandler{app: ser.app, Handle: handleCalculate}))
	router.PathPrefix("/progress").Handler(appHandler{app: ser.app, Handle: progress})

	router.HandleFunc("/health", healthcheck)

	fileServer := http.FileServer(http.Dir("./static"))
	router.Handle("/", fileServer)

}

func NewServer(logger *zap.Logger, a *app.App) server {

	addr := config.GetPort()
	logger.Info(fmt.Sprint("listening on address: ", addr))

	return server{
		logger: logger,
		app:    a,
		addr:   addr,
	}
}

func (ser server) Run() error {
	router := mux.NewRouter()
	ser.registerHandlers(router)

	ser.httpServer = &http.Server{
		Handler:  router,
		ErrorLog: zap.NewStdLog(ser.logger),
		Addr:     ser.addr,
	}

	return ser.httpServer.ListenAndServe()
}

func (appHndl appHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	appHndl.app.Logger.Debug("request received", zap.String("method", r.Method), zap.String("url", r.URL.Path), zap.String("content-type", r.Header.Get("Content-Type")))

	status, err := appHndl.Handle(appHndl.app, w, r)

	if err != nil {
		appHndl.app.Logger.Warn("request failed", zap.Error(err))
		w.Write([]byte(fmt.Sprintln(err)))

		switch status {
		case http.StatusNotFound:
			http.NotFound(w, r)
		case http.StatusInternalServerError:
			http.Error(w, http.StatusText(status), status)
		default:
			http.Error(w, http.StatusText(status), status)
		}
		return
	}
}
