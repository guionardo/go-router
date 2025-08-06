package headers

import (
	"fmt"
	"net/http"
	"reflect"
	"strings"
)

type (
	Headers[T any] struct {
		headers map[string]*Header
	}
	Header struct {
		Name       string
		FieldOrder int
		Required   bool
	}
)

func New[T any]() (*Headers[T], error) {
	t := reflect.TypeFor[T]()
	if t.Kind() != reflect.Struct {
		return nil, fmt.Errorf("type %s should be an struct. Got %s", t.Name(), t.Kind().String())
	}

	headers := make(map[string]*Header)

	for i := range t.NumField() {
		field := t.Field(i)
		if !field.IsExported() {
			continue
		}
		headerTag := field.Tag.Get("header")
		if len(headerTag) == 0 {
			continue
		}
		ft := field.Type.Kind()

		if ft != reflect.String {
			return nil, fmt.Errorf("field %s from type %s should be string. Got %s", field.Name, t.Name(), ft.String())
		}
		required := false
		w := strings.Split(headerTag, ",")
		if len(w) > 1 && strings.Contains(headerTag, ",required") {
			required = true
			headerTag = w[0]
		}
		if !required {
			validateTag := field.Tag.Get("validate")
			if strings.Contains(validateTag, "required") {
				required = true
			}
		}
		header := &Header{
			Required:   required,
			FieldOrder: i,
		}
		headers[headerTag] = header
	}
	return &Headers[T]{
		headers: headers,
	}, nil
}

func (h *Headers[T]) Populate(r *http.Request, v *T) (err error) {
	if v == nil {
		return fmt.Errorf("can not unmarshal into a nil value")
	}

	rv := reflect.ValueOf(v)

	ps := rv.Elem()
	for name, header := range h.headers {
		value := r.Header.Get(name)
		if len(value) == 0 {
			if header.Required {
				return fmt.Errorf("missing required header %s from request", name)
			}
			continue
		}
		field := ps.Field(header.FieldOrder)
		if field.IsValid() && field.CanSet() {
			field.SetString(value)
		}
	}
	return nil
}
