package router

import (
	"fmt"
	"log/slog"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/guionardo/go-router/pkg/inspect"
	"github.com/guionardo/go-router/pkg/logging"
)

type (
	Router struct {
		info      *RouterInfo
		endpoints map[string]map[string]inspect.HandlerStruct
	}
	RouterOption func(*Router)
)

func New(options ...RouterOption) *Router {
	r := &Router{
		info:      &RouterInfo{},
		endpoints: make(map[string]map[string]inspect.HandlerStruct),
	}
	for _, option := range options {
		option(r)
	}
	return r
}

func (r *Router) Add(method string, endpoint inspect.HandlerStruct) *Router {
	method = strings.ToUpper(method)
	msh, ok := r.endpoints[method]
	if !ok {
		msh = make(map[string]inspect.HandlerStruct)
		r.endpoints[method] = msh
	}
	msh[endpoint.GetPath()] = endpoint

	return r
}

func (r *Router) Get(endpoints ...inspect.HandlerStruct) *Router {
	for _, endpoint := range endpoints {
		r.Add(http.MethodGet, endpoint)
	}
	return r
}

func (r *Router) SetupHTTP(h *http.ServeMux) {
	for method, handlers := range r.endpoints {
		for path, handler := range handlers {
			logging.Get().Debug("Route", slog.String("method", method), slog.String("path", path), slog.String("handler", handler.HandlerName()))
			h.HandleFunc(fmt.Sprintf("%s %s", method, path), handler.Handle)
		}
	}
}

func (r *Router) SetupGin(h *gin.Engine) {
	for method, handlers := range r.endpoints {
		for path, handler := range handlers {
			logging.Get().Debug("Route", slog.String("method", method), slog.String("path", path), slog.String("handler", handler.HandlerName()))
			h.Handle(method, path, func(c *gin.Context) {
				handler.Handle(c.Writer, c.Request)
			})
		}
	}
}
