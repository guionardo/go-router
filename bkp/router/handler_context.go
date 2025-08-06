package router

import (
	"context"
	"encoding/json"
	"net/http"
)

type HandlerContext struct {
	Context        context.Context
	Headers        map[string]string
	responseHeader map[string]string
	statusCode     int
}

func (h *HandlerContext) SetHeader(name, value string) *HandlerContext {
	h.responseHeader[name] = value
	return h
}

func (h *HandlerContext) SetStatusCode(statusCode int) *HandlerContext {
	h.statusCode = statusCode
	return h
}

func (h *HandlerContext) doResponse(w http.ResponseWriter, responseBody any, statusCode int, err error) {
	if err == nil {
		if h.statusCode > 0 {
			statusCode = h.statusCode
		}
		w.WriteHeader(statusCode)
		if responseBody != nil {
			h.responseHeader["Content-Type"] = "application/json"
		}
		for name, value := range h.responseHeader {
			w.Header().Set(name, value)
		}
		if responseBody != nil {
			encoder := json.NewEncoder(w)
			_ = encoder.Encode(responseBody)
		}
		return
	}
	if _, ok := err.(*Error); !ok {
		err = NewError(statusCode, err)
	}
	handleError(w, err)
}
