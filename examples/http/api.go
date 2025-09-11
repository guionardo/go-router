package main

import (
	"log/slog"
	"net/http"
	"os"

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
	// TODO: Implement example

	r.SetupHTTP(mux)

	http.ListenAndServe(":8080", mux)
}
