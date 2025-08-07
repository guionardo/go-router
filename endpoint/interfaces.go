package endpoint

import (
	"net/http"

	"github.com/go-playground/validator"
)

type (
	HandlerStruct interface {
		GetPath() string
		Handle(http.ResponseWriter, *http.Request)
		HandlerName() string
		PathParams() []string
	}
	Responser[T any, R any] interface {
		Handle(r *http.Request, payload *T) (response *R, status int, err error)
	}
	CustomResponser[T any, R any] interface {
		Handle(w http.ResponseWriter, r *http.Request, payload *T) error
	}
	customValidator interface {
		Validate() error
	}
	postParser interface {
		PostParse(error) error
	}
	setupValidator interface {
		NewValidator() *validator.Validate
	}

	responseError interface {
		Error() string
		StatusCode() int
	}
)
