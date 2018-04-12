package mylog

import (
	golog "github.com/op/go-logging"
	"io"
	"net/http"
)

const (
	requestStartLogTemplate   = `Started handling request to url %v with method %v`
	requestSuccessLogTemplate = `Request to url %v with method %v handled successfully`
	requestBodyTemplate       = "Request to url %v with method %v has body %v"
	responseBodyTemplate      = "Request to url %v with method %v has response with body %v"
	requestErrorTemplate      = `Failed on URL %v with error \"%v\"`
	logFormat                 = `%{color}%{time:15:04:05.000} %{shortfunc} ▶ %{level:.4s} %{id:03x}%{color:reset} %{message}`
)

func NewLogger(writer io.Writer) *Logger {
	format := golog.MustStringFormatter(logFormat)
	backend := golog.NewLogBackend(writer, "", 0)
	backendFormatter := golog.NewBackendFormatter(backend, format)

	backendLeveled := golog.AddModuleLevel(backendFormatter)
	backendLeveled.SetLevel(golog.INFO, "")

	logger := golog.MustGetLogger("main")

	logger.SetBackend(backendLeveled)

	return &Logger{*logger}
}

type Logger struct {
	golog.Logger
}

func (logger *Logger) LogRequestBody(r *http.Request, body string) {
	logger.Infof(requestBodyTemplate, r.URL.Path, r.Method, body)
}

func (logger *Logger) LogResponseBody(r *http.Request, body string) {
	logger.Infof(responseBodyTemplate, r.URL.Path, r.Method, body)
}

func (logger *Logger) LogRequestStart(r *http.Request) {
	logger.Infof(requestStartLogTemplate, r.URL.Path, r.Method)
}

func (logger *Logger) LogRequestSuccess(r *http.Request) {
	logger.Infof(requestSuccessLogTemplate, r.URL.Path, r.Method)
}

func (logger *Logger) LogRequestError(r *http.Request, err error) {
	logger.Errorf(requestErrorTemplate, r.URL.Path, err.Error())
}
