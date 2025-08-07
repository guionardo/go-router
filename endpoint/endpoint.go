package endpoint

import (
	"errors"
	"maps"
	"net/http"
	"reflect"
	"slices"
)

type (
	Endpoint[T any, R any] struct {
		path          string
		useValidation bool
		raw           []byte
		paths         map[string]int
		queries       map[string]int
		headers       map[string]int
		bodyField     string
		bodyParseType bodyParseType
		reqType       reflect.Type
		respType      reflect.Type

		parseBodyFunc    func(*http.Request, *T) error
		parsePathFunc    func(*http.Request, *T) error
		parseQueriesFunc func(*http.Request, *T) error
		parseHeadersFunc func(*http.Request, *T) error
		validateFunc     func(*http.Request, *T) error

		isCustomValidator bool
		isPostParser      bool
		isSetupValidator  bool
		initialized       bool

		responser       Responser[T, R]
		customResponser CustomResponser[T, R]
		handlerFunc     func(http.ResponseWriter, *http.Request, *T) error
		handlerName     string
	}

	bodyParseType uint8
)

var (
	errEmptyTag = errors.New("empty_tag")
)

const (
	ContentJSON = "application/json"
	ContentXML  = "application/xml"
	ContentYAML = "application/yaml"

	BodyNo bodyParseType = iota
	BodyFull
	BodyBytes
	BodyString
	BodyJSON
)

func New[T, R any](path string) (is *Endpoint[T, R]) {
	var t reflect.Type
	is, t, err := getEndpoint[T, R]()
	if err != nil {
		panic(err)
	}
	if is == nil || !is.initialized {
		if is, err = buildStruct[T, R](t, path); err == nil {
			setEndpoint(t, is)
		} else {
			panic(err)
		}
	}
	return is
}

func (s *Endpoint[T, R]) GetPath() string {
	return s.path
}

func (s *Endpoint[T, R]) HandlerName() string {
	return s.reqType.String()
}

func (s *Endpoint[T, R]) PathParams() []string {
	return slices.Collect(maps.Keys(s.paths))
}
