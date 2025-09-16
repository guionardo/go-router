package payloads

import "net/http"

type (
	UserRequest struct {
		Id   int    `path:"id"`
		Auth string `header:"auth" validate:"required"`
	}
	UserResponse struct {
		OldId int `json:"old_id"`
		NewId int `json:"new_id"`
	}
)

func (pr *UserRequest) Handle(r *http.Request, payload *UserRequest) (response *UserResponse, status int, err error) {
	return &UserResponse{OldId: payload.Id, NewId: payload.Id + 1}, 200, nil
}
