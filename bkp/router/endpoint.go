package router

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/go-playground/validator"
	"github.com/guionardo/go-router/pkg/path_params"
	"github.com/guionardo/go-router/router/payload"
	"github.com/lestrrat-go/urlenc"
)

type (
	// Endpoint is a struct to handle a http request
	// T is the request type
	// R is the response type
	Endpoint[T any, R any] struct {
		Method  string
		Path    string
		Payload *payload.RequestData[T]
		Handler EndpointHandler[T, R]

		headers      Headers
		pathParams   *path_params.PathParams
		emptyPayload bool
		info         *Info
	}

	EndpointHandler[T any, R any] func(ctx *HandlerContext, payload *T) (response *R, statusCode int, err error)
	Option                        func(endpointer)
	endpointer                    interface {
		getHeaders() Headers
		setInfo(info *Info)
		Handle(w http.ResponseWriter, r *http.Request)
	}
	EmptyPayload  struct{}
	BinaryPayload []byte
	Info          struct {
		Summary     string
		Description string
	}
)

var (
	ErrEmptyPayload = errors.New("empty payload") //nolint:revive
)

func NewEndpoint[T any, R any](method string, path string, handler EndpointHandler[T, R], options ...Option) *Endpoint[T, R] {
	payload, err := payload.New[T](path)
	if err != nil {
		panic(err) // TODO: Verificar se panic está correto aqui
	}

	pathParams, err := path_params.NewPathParams(path)
	if err != nil {
		panic(err)
	}
	if len(pathParams.ParamNames()) == 0 {
		pathParams = nil
	}
	ep := &Endpoint[T, R]{
		Method:  method,
		Path:    path,
		Handler: handler,

		headers:      make(Headers),
		pathParams:   pathParams,
		emptyPayload: fmt.Sprintf("%T", new(T)) == fmt.Sprintf("%T", new(EmptyPayload)),
	}
	for _, option := range options {
		option(ep)
	}
	return ep
}

func (e *Endpoint[T, R]) parseAndValidatePayload(r *http.Request) (*T, error) {
	payload, err := e.parsePayload(r)
	if err != nil {
		return nil, err
	}
	if err := e.pathParams.Match(r.URL); err != nil {
		return nil, err
	}
	if err := e.pathParams.Inject(payload); err != nil {
		return nil, err
	}
	if err := e.validatePayload(payload); err != nil {
		return nil, err
	}
	return payload, nil
}

func (e *Endpoint[T, R]) Handle(w http.ResponseWriter, r *http.Request) {
	headers, err := e.headers.Validate(r)
	if err != nil {
		handleError(w, err)
		return
	}
	payload, err := e.parseAndValidatePayload(r)
	if err != nil {
		handleError(w, err)
		return
	}

	ctx := &HandlerContext{
		Context:        r.Context(),
		Headers:        headers,
		responseHeader: make(map[string]string),
	}
	response, statusCode, err := e.Handler(ctx, payload)
	ctx.doResponse(w, response, statusCode, err)

}

func handleError(w http.ResponseWriter, err error) {
	if er, ok := err.(*Error); ok {
		http.Error(w, er.Error(), er.statusCode)
	} else {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (e *Endpoint[T, R]) parsePayload(r *http.Request) (payload *T, err error) {
	if e.emptyPayload {
		return nil, nil
	}

	payload = new(T)

	if err := e.tryParseBody(r, payload); err == nil {
		return payload, nil
	}

	if err := e.tryParseUrl(r, payload); err == nil {
		return payload, nil
	}

	return nil, errors.New("failed to parse payload")

}

func (e *Endpoint[T, R]) tryParseBody(r *http.Request, requestData *T) error {
	if r.Body == nil {
		return ErrEmptyPayload
	}
	body, err := io.ReadAll(r.Body)
	if err == nil && len(body) == 0 {
		return nil
	}

	if err == nil {
		err = json.Unmarshal(body, requestData)
	}
	return err
}

func (e *Endpoint[T, R]) tryParseUrl(r *http.Request, requestData *T) error {
	query := r.URL.RawQuery
	if err := urlenc.Unmarshal([]byte(query), requestData); err != nil {
		return err
	}
	return nil
}

func (e *Endpoint[T, R]) getHeaders() Headers {
	return e.headers
}

func (e *Endpoint[T, R]) setInfo(info *Info) {
	e.info = info
}

func (e *Endpoint[T, R]) validatePayload(payload *T) error {
	if e.emptyPayload {
		return nil
	}

	return validator.New().Struct(payload)
}
