package sceneries

import (
	"net/http"
	"time"
)

type (
	PostSimple struct {
		APIKey string         `header:"apikey" validate:"required"`
		ID     int            `path:"id"`
		Option int            `query:"option"`
		Date   time.Time      `query:"date"`
		Body   PostSimpleBody `body:"body"`
	}
	PostSimpleBody struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}
	PostSimpleResponse struct {
		Length int `json:"length"`
	}
)

func (gs *PostSimple) Handle(r *http.Request, payload *PostSimple) (*PostSimpleResponse, int, error) {
	time.Sleep(time.Millisecond * 100) // spend some time on an important process
	return &PostSimpleResponse{payload.Body.Age}, http.StatusAccepted, nil
}
