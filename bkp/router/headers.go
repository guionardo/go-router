package router

import (
	"net/http"
	"strings"
)

type (
	Header struct {
		Name      string
		Required  bool
		Value     *string
		Validator HeaderValidatorFunc
	}
	Headers             map[string]*Header
	HeaderValidatorFunc func(headerValue *string) error
)

func (h Headers) Validate(r *http.Request) (map[string]string, error) {
	headers := make(map[string]string)
	for headerName, header := range h {
		headerValue := r.Header.Get(headerName)
		if header.Required && len(headerValue) == 0 {
			return nil, NewErrorF(http.StatusBadRequest, "header %s is required", headerName)
		}
		if header.Validator != nil {
			if err := header.Validator(&headerValue); err != nil {
				return nil, NewError(http.StatusBadRequest, err)
			}
		}
		headers[headerName] = headerValue
	}
	return headers, nil
}

func WithHeader(headerName string, required bool) Option {
	return func(e endpointer) {
		e.getHeaders()[strings.ToLower(headerName)] = &Header{Name: headerName, Required: true}
	}
}

func WithHeaderFunc(headerName string, validator HeaderValidatorFunc) Option {
	return func(e endpointer) {
		e.getHeaders()[strings.ToLower(headerName)] = &Header{Name: headerName, Validator: validator}
	}
}
