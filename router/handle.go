package router

import (
	"encoding/json"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/guionardo/go-router/pkg/tools"
)

type (
	HasStatusCode interface {
		StatusCode() int
	}

	BadRequestError struct {
		Message    string `json:"message"`
		StatusCode int    `json:"status_code"`
	}
)

func Handle(w http.ResponseWriter, response any, err error) error {
	var statusCode int
	switch e := err.(type) {
	case *tools.ParseError, validator.FieldError:
		statusCode = http.StatusBadRequest
		response = BadRequestError{
			Message:    e.Error(),
			StatusCode: statusCode,
		}
		w.WriteHeader(statusCode)
		return encode(w, response)
	}
	if e, ok := err.(HasStatusCode); ok {
		statusCode = e.StatusCode()
	} else if r, ok := response.(HasStatusCode); ok {
		statusCode = r.StatusCode()
	} else if err == nil {
		statusCode = http.StatusOK
	} else {
		statusCode = http.StatusBadGateway
	}

	w.WriteHeader(statusCode)
	if response != nil {
		return encode(w, response)
	}
	return nil
}

func encode(w http.ResponseWriter, body any) error {
	var (
		content []byte
		err     error
	)
	if m, ok := body.(json.Marshaler); ok {
		content, err = m.MarshalJSON()
	} else {
		content, err = json.Marshal(body)
	}
	if err == nil {
		_, err = w.Write(content)
	}
	return err
}
