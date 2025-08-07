package payloads

import (
	"fmt"
	"net/http"
)

type (
	ProductRequest struct {
		Id int `path:"id"`
	}
	ProductResponse struct {
		Message string `json:"message"`
	}
)

func (pr *ProductRequest) Handle(r *http.Request, payload *ProductRequest) (*ProductResponse, int, error) {
	return nil, http.StatusBadGateway, fmt.Errorf("Forced error code: %d", payload.Id)
}
