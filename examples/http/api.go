package main

import (
	"net/http"

	"github.com/guionardo/go-router/pkg/inspect"
	"github.com/guionardo/go-router/router"
)

type (
	PingRequest struct {
	}
	PingResponse struct {
		Message string
	}
	UserRequest struct {
		Id   int    `path:"id"`
		Auth string `header:"auth" validate:"required"`
	}
	UserResponse struct {
		OldId int `json:"old_id"`
		NewId int `json:"new_id"`
	}
)

func (pr *PingRequest) Handle(r *http.Request, payload *PingRequest) (response *PingResponse, status int, err error) {
	return nil, 200, nil
}

func createServer() *http.ServeMux {
	mux := http.NewServeMux()
	return mux
}

func main() {
	mux := createServer()
	r := router.New(router.Title("PING API"), router.Version("0.1.0"))

	// ePing := router.NewEndpoint(http.MethodGet, "/ping", func(ctx *router.HandlerContext, payload *router.EmptyPayload) (response *PingResponse, statusCode int, err error) {
	// 	return &PingResponse{"PONG"}, http.StatusOK, nil
	// })
	ePing, err := inspect.New[PingRequest, PingResponse]("/ping")
	if err != nil {
		panic(err)
	}
	r.Add(http.MethodGet, ePing)

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
