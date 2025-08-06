package inspect

import "net/http"

type (
	HandlerStruct interface {
		GetPath() string
		Handle(http.ResponseWriter, *http.Request)
	}
)
