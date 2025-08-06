package inspect

import (
	"fmt"
	"reflect"

	"github.com/go-playground/validator"
)

// PayloadStruct will 
type PayloadStruct[T any] struct {
	PathFields  []string
	QueryFields []string
	Data        T
	useValidate bool
}

func NewEndpointStruct[T any](path string) *PayloadStruct[T] {
	t := reflect.TypeFor[T]()
	var (
		pathFields  = make([]string, 0)
		queryFields = make([]string, 0)
		useValidate = false
	)
	for i := range t.NumField() {
		field := t.Field(i)
		if !field.IsExported() {
			continue
		}
		pathTag := field.Tag.Get("path")
		queryTag := field.Tag.Get("urlenc")
		jsonTag := field.Tag.Get("json")
		if len(pathTag) > 0 {
			if len(queryTag) == 0 && len(jsonTag) == 0 {
				panic(fmt.Sprintf("add a tag \"urlenc:%s\" to field %s.%s", pathTag, t.String(), field.Name))
			}
		}
		if len(pathTag) > 0 {
			pathFields = append(pathFields, pathTag)
		}
		if len(queryTag) > 0 {
			queryFields = append(queryFields, queryTag)
		}
		if len(field.Tag.Get("validate")) > 0 {
			useValidate = true
		}
	}

	return &PayloadStruct[T]{
		PathFields:  pathFields,
		QueryFields: queryFields,
		useValidate: useValidate,
	}
}

func (es *PayloadStruct[T]) Validate(value T) error {
	if !es.useValidate {
		return nil
	}
	v := validator.New()
	return v.Struct(value)
}
