package main

import (
	"log/slog"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/guionardo/go-router/endpoint"
	"github.com/guionardo/go-router/examples/http/payloads"
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
	g := gin.Default()

	r := router.New(router.Title("PING API"), router.Version("0.1.0"))

	ePing := endpoint.New[payloads.PingRequest, payloads.PingResponse]("/ping")
	eUser := endpoint.New[payloads.UserRequest, payloads.UserResponse]("/user/{id}")
	eProduct := endpoint.New[payloads.ProductRequest, payloads.ProductResponse]("/prod/{id}")

	r.Get(ePing, eUser, eProduct)

	r.SetupHTTP(mux)
	r.SetupGin(g)
	http.ListenAndServe(":8080", g)
}
