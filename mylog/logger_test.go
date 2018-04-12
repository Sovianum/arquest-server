package mylog

import (
	"bytes"
	"fmt"
	"github.com/op/go-logging"
	"io"
	"net/http"
	"testing"
)

const (
	method = "METHOD"
	url    = "/URL"
)

func TestLogger_LogRequestStart(t *testing.T) {
	var writer bytes.Buffer
	logger := getLogger(&writer)
	req, _ := http.NewRequest(
		method,
		url,
		nil,
	)

	logger.LogRequestStart(req)
	msg := string(writer.Bytes())
	expected := fmt.Sprintf(
		requestStartLogTemplate,
		req.URL.Path,
		req.Method,
	)

	if msg[:len(msg)-1] != expected { //msg[:len(msg) - 1] removes last \n symbol
		t.Errorf("Expected \"%v\"; got \"%v\"", expected, msg)
	}
}

func TestLogger_LogRequestSuccess(t *testing.T) {
	var writer bytes.Buffer
	logger := getLogger(&writer)
	req, _ := http.NewRequest(
		method,
		url,
		nil,
	)

	logger.LogRequestSuccess(req)
	msg := string(writer.Bytes())
	expected := fmt.Sprintf(
		requestSuccessLogTemplate,
		req.URL.Path,
		req.Method,
	)

	if msg[:len(msg)-1] != expected { //msg[:len(msg) - 1] removes last \n symbol
		t.Errorf("Expected \"%v\"; got \"%v\"", expected, msg)
	}
}

func TestLogger_LogRequestError(t *testing.T) {
	var writer bytes.Buffer
	logger := getLogger(&writer)
	req, _ := http.NewRequest(
		method,
		url,
		nil,
	)

	errorMsg := "msg"
	err := fmt.Errorf(errorMsg)

	logger.LogRequestError(req, err)
	msg := string(writer.Bytes())
	expected := fmt.Sprintf(
		requestErrorTemplate,
		req.URL.Path,
		err.Error(),
	)

	if msg[:len(msg)-1] != expected { //msg[:len(msg) - 1] removes last \n symbol
		t.Errorf("Expected \"%v\"; got \"%v\"", expected, msg)
	}
}

func getLogger(writer io.Writer) *Logger {
	format := logging.MustStringFormatter(
		`%{message}`,
	)
	backend := logging.NewLogBackend(writer, "", 0)
	backendFormatter := logging.NewBackendFormatter(backend, format)

	backendLeveled := logging.AddModuleLevel(backendFormatter)
	backendLeveled.SetLevel(logging.INFO, "")

	logger := logging.MustGetLogger("main")

	logger.SetBackend(backendLeveled)

	return &Logger{*logger}
}
