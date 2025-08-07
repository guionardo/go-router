package router

import (
	"fmt"
	"log/slog"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/guionardo/go-router/endpoint"
	"github.com/guionardo/go-router/pkg/logging"
	"github.com/guionardo/go-router/pkg/path_params"
)

type (
	Router struct {
		info      *RouterInfo
		endpoints map[string]map[string]endpoint.HandlerStruct
	}
	RouterOption func(*Router)
)

func New(options ...RouterOption) *Router {
	r := &Router{
		info:      &RouterInfo{},
		endpoints: make(map[string]map[string]endpoint.HandlerStruct),
	}
	for _, option := range options {
		option(r)
	}
	return r
}

func (r *Router) Add(method string, ep endpoint.HandlerStruct) *Router {
	method = strings.ToUpper(method)
	msh, ok := r.endpoints[method]
	if !ok {
		msh = make(map[string]endpoint.HandlerStruct)
		r.endpoints[method] = msh
	}
	msh[ep.GetPath()] = ep

	return r
}

func (r *Router) Get(endpoints ...endpoint.HandlerStruct) *Router {
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
			// TODO: Fix path to use :params
			path, err := path_params.InvalidatePathToGin(path)
			if err == nil {
				h.Handle(method, path, func(c *gin.Context) {
					// Fix params into http.Request from gin.Context
					for _, pp := range handler.PathParams() {
						c.Request.SetPathValue(pp, c.Param(pp))
					}

					handler.Handle(c.Writer, c.Request)
				})
			}
		}
	}
}
