package inspect

import "net/http"

type (
	HandlerStruct interface {
		GetPath() string
		Handle(http.ResponseWriter, *http.Request)
		HandlerName() string
	}
	Responser[T any, R any] interface {
		Handle(r *http.Request, payload *T) (response *R, status int, err error)
	}
	CustomResponser[T any, R any] interface {
		Handle(w http.ResponseWriter, r *http.Request, payload *T) error
	}
)
