package request_data

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"reflect"

	"github.com/guionardo/go-router/pkg/inspect"
	"github.com/guionardo/go-router/pkg/path_params"
	"github.com/lestrrat-go/urlenc"
)

type (
	// RequestData is a struct to read
	RequestData[T any] struct {
		pathData        *path_params.PathParams
		payloadType     PayloadType
		bodyHandler     func(*http.Request) (T, error) // handler for reading body of request
		endpointStruct  *inspect.PayloadStruct[T]
		ignoreBody      bool
		expectedHeaders []string
	}

	PayloadType       byte
	RequestDataOption func(requester)

	requester interface {
		setIgnoreBody(bool)
	}
)

const (
	PayloadTypeBytes       PayloadType = iota // raw []byte
	PayloadTypeString                         // string
	PayloadTypeStruct                         // struct
	PayloadTypeStructArray                    // []struct
	PayloadTypeReader                         // io.Reader
	PayloadTypeError
)

func New[T any](path string, options ...RequestDataOption) (*RequestData[T], error) {
	payloadType, err := getPayloadType[T]()
	if err != nil {
		return nil, err
	}
	pathParams, err := path_params.NewPathParams(path)
	if err != nil {
		return nil, err
	}

	pl := &RequestData[T]{
		payloadType: payloadType,
		pathData:    pathParams,
	}
	switch payloadType {
	case PayloadTypeBytes:
		pl.bodyHandler = pl.handleBytes
	case PayloadTypeString:
		pl.bodyHandler = pl.handleString
	case PayloadTypeReader:
		pl.bodyHandler = pl.handleReader
	case PayloadTypeStructArray:
		pl.bodyHandler = pl.handleStructArray
	case PayloadTypeStruct:
		pl.endpointStruct = inspect.NewEndpointStruct[T](path)
		//TODO: Implementar match entre os campos do path e os disponíveis na struct
		// TODO: Implementar verificação nos campos da struct pela tag header
		pl.bodyHandler = pl.handleStruct
	}
	for _, option := range options {
		option(pl)
	}
	return pl, nil
}

func (p *RequestData[T]) Handle(r *http.Request) (output T, err error) {
	return p.bodyHandler(r)
}

func (p *RequestData[T]) handleBytes(r *http.Request) (out T, err error) {
	var body any
	body, err = io.ReadAll(r.Body)
	if err == nil {
		r.Body.Close()
		return body.(T), err
	}
	return out, err
}

func (p *RequestData[T]) handleString(r *http.Request) (out T, err error) {
	var body any
	body, err = io.ReadAll(r.Body)
	if err == nil {
		r.Body.Close()
		var str any = string(body.([]byte))
		return str.(T), nil
	}
	return out, err
}

func (p *RequestData[T]) handleReader(r *http.Request) (out T, err error) {
	if r.Body != nil {
		var reader any = r.Body
		return reader.(T), nil
	}
	return out, fmt.Errorf("request body is empty")
}

func (p *RequestData[T]) handleStructArray(r *http.Request) (out T, err error) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return out, err
	}

	r.Body.Close()
	err = json.Unmarshal(body, out)
	return out, err
}

func (p *RequestData[T]) handleStruct(r *http.Request) (out T, err error) {
	if !p.ignoreBody {
		if err = json.NewDecoder(r.Body).Decode(&out); err != nil {
			return out, err
		}
	}
	// Parse path and query data
	if err = p.pathData.MatchAndInject(r.URL, &out); err != nil {
		return out, err
	}
	path, err := url.PathUnescape(r.URL.RawQuery)
	if err != nil {
		return out, err
	}
	if err = urlenc.Unmarshal([]byte(path), &out); err != nil {
		return out, err
	}
	// for _,header:=range p.expectedHeaders{
	// 	// TODO: Utilizar header
	// }

	if err := p.endpointStruct.Validate(out); err != nil {
		return out, err
	}

	return out, nil

}

func getPayloadType[T any]() (PayloadType, error) {
	if inspect.IsArrayOfByte[T]() {
		return PayloadTypeBytes, nil
	}
	if inspect.IsString[T]() {
		return PayloadTypeString, nil
	}
	if inspect.IsStruct[T]() {
		return PayloadTypeStruct, nil
	}
	if inspect.IsArrayOfStruct[T]() {
		return PayloadTypeStructArray, nil
	}
	if inspect.IsReader[T]() {
		return PayloadTypeReader, nil
	}
	return PayloadTypeError, fmt.Errorf("invalid payload type: %s", reflect.TypeFor[T]().String())
}

func (p *RequestData[T]) setIgnoreBody(ignoreBody bool) {
	p.ignoreBody = ignoreBody
}

func WithIgnoreBody(p requester) {
	p.setIgnoreBody(true)
}

func (p *RequestData[T]) analiseStructFields() {
	t := reflect.TypeFor[T]()
	if t.Kind() == reflect.Pointer {
		t = t.Elem()
	}
	for i := range t.NumField() {
		field := t.Field(i)
		if tagHeader, ok := field.Tag.Lookup("header"); ok {
			// header should be only the name of the header
			p.expectedHeaders = append(p.expectedHeaders, tagHeader)
		}
	}
}
