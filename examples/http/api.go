package main

import (
	"log/slog"
	"net/http"
	"os"

	"github.com/guionardo/go-router/examples/http/payloads"
	"github.com/guionardo/go-router/pkg/inspect"
	"github.com/guionardo/go-router/router"
)

func createServer() *http.ServeMux {
	mux := http.NewServeMux()
	return mux
}

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	slog.SetDefault(logger)
	mux := createServer()
	r := router.New(router.Title("PING API"), router.Version("0.1.0"))

	ePing, err := inspect.New[payloads.PingRequest, payloads.PingResponse]("/ping")
	if err != nil {
		panic(err)
	}
	eUser, err := inspect.New[payloads.UserRequest, payloads.UserResponse]("/user/{id}")
	if err != nil {
		panic(err)
	}

	eProduct, err := inspect.New[payloads.ProductRequest, payloads.ProductResponse]("/prod/{id}")
	if err != nil {
		panic(err)
	}
	r.Get(ePing, eUser, eProduct)

	// eUser := router.NewEndpoint(http.MethodGet, "/user/:id", func(ctx *router.HandlerContext, payload *UserRequest) (response *UserResponse, statusCode int, err error) {
	// 	return &UserResponse{
	// 		OldId: payload.Id,
	// 		NewId: -payload.Id,
	// 	}, http.StatusOK, nil
	// })
	// r.Add(eUser)
	r.SetupHTTP(mux)
	http.ListenAndServe(":8080", mux)
}
