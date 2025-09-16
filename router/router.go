package router

import (
	"fmt"
	"log/slog"
	"net/http"
	"strings"

	"github.com/guionardo/go-router/pkg/logging"
)

type (
	Router struct {
		info      *RouterInfo
		endpoints map[string]map[string]Handler
	}
	RouterOption func(*Router)
)

func New(options ...RouterOption) *Router {
	r := &Router{
		info:      &RouterInfo{},
		endpoints: make(map[string]map[string]Handler),
	}
	for _, option := range options {
		option(r)
	}
	return r
}

func (r *Router) Add(method string, path string, handler Handler) *Router {
	method = strings.ToUpper(method)
	msh, ok := r.endpoints[method]
	if !ok {
		msh = make(map[string]Handler)
		r.endpoints[method] = msh
	}
	msh[path] = handler

	return r
}

func (r *Router) Get(path string, handler Handler) *Router {
	r.Add(http.MethodGet, path, handler)
	return r
}

func (r *Router) Post(path string, handler Handler) *Router {
	r.Add(http.MethodPost, path, handler)
	return r
}

func (r *Router) SetupHTTP(h *http.ServeMux) {
	for method, handlers := range r.endpoints {
		for path, handler := range handlers {
			logAddHandler(method, path, handler)
			h.HandleFunc(fmt.Sprintf("%s %s", method, path), handler.Handle)
		}
	}
}

func logAddHandler(method string, path string, handler Handler) {
	logging.Get().Debug("Route",
		slog.String("method", method),
		slog.String("path", path),
		// slog.String("handler", handler.HandlerName()))
	)
}
