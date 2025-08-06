package router

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/guionardo/go-router/pkg/config"
)

type Error struct {
	err        error
	statusCode int

	Message string `json:"message,omitempty"`
	Status  int    `json:"status_code,omitempty"`
}

func (e *Error) Error() string {
	return fmt.Sprintf("status: %d - %s, error: %s", e.statusCode, http.StatusText(e.statusCode), e.err.Error())
}

func (e *Error) StatusCode() int {
	return e.statusCode
}

func NewError(statusCode int, err error) *Error {
	if err == nil {
		err = errors.New("nil error")
	}

	if config.DevelopmentMode {
		return &Error{err: err, statusCode: statusCode, Message: err.Error(), Status: statusCode}
	}
	return &Error{err: err, statusCode: statusCode}
}

func NewErrorF(statusCode int, format string, args ...any) *Error {
	return NewError(statusCode, fmt.Errorf(format, args...))
}
