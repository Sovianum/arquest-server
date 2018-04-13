package dao

import "net/http"

func NewDBErr(code int, msg string) DBError {
	return &dbError{
		msg:  msg,
		code: code,
	}
}

func NewCrashDBErr(err error) DBError {
	if err == nil {
		return nil
	}
	return &dbError{
		msg:  err.Error(),
		code: http.StatusInternalServerError,
	}
}

type DBError interface {
	error
	Code() int
}

type dbError struct {
	msg  string
	code int
}

func (err *dbError) Code() int {
	return err.code
}

func (err *dbError) Error() string {
	return err.msg
}
