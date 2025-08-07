package payloads

import "net/http"

type (
	PingRequest struct {
	}
	PingResponse struct {
		Message string
	}
)

func (pr *PingRequest) Handle(r *http.Request, payload *PingRequest) (response *PingResponse, status int, err error) {
	return &PingResponse{"PONG"}, 200, nil
}
