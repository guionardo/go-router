package sceneries

import (
	"net/http"
	"time"
)

type (
	// GetSimple only receives an required integer Id from query string
	// GET /?id=1
	GetSimple struct {
		Id int `query:"id" validate:"required"`
	}
	GetSimpleResponse struct {
		Success bool `json:"success"`
	}
)

func (gs *GetSimple) Handle(r *http.Request, payload *GetSimple) (*GetSimpleResponse, int, error) {
	time.Sleep(time.Millisecond * 100) // spend some time on an important process
	return &GetSimpleResponse{true}, http.StatusAccepted, nil
}
