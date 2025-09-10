package generator

import (
	"fmt"
	"reflect"

	"github.com/mailru/easyjson"
)

type (
	Generator[T Endpointer] struct {
		sourceFile    string
		pathFields    []int
		queryFields   []int
		jsonFields    []int
		bodyType      reflect.Type
		hasValidation bool
	}
	Endpointer interface {
		easyjson.MarshalerUnmarshaler
	}
)

func NewGenerator[T Endpointer](sourceFile string) (*Generator[T], error) {
	t := reflect.TypeFor[T]()
	var (
		pathFields,
		queryFields,
		jsonFields []int
		tagPath, tagQuery, tagJson, tagBody string
		bodyType                            reflect.Type
		hasBody, hasValidation              bool
	)
	for fieldN := range t.NumField() {
		field := t.Field(fieldN)
		if !field.IsExported() {
			continue
		}
		if tagPath = field.Tag.Get("path"); len(tagPath) > 0 {
			pathFields = append(pathFields, fieldN)
		}
		if tagQuery = field.Tag.Get("query"); len(tagQuery) > 0 {
			queryFields = append(queryFields, fieldN)
		}
		if tagJson = field.Tag.Get("json"); len(tagJson) > 0 && tagJson != "-" {
			jsonFields = append(jsonFields, fieldN)
		}
		if tagBody = field.Tag.Get("body"); len(tagBody) > 0 {
			if bodyType != nil {
				return nil, fmt.Errorf("body field was defined more than one time at struct %s", t.String())
			}
			fType := field.Type
			if fType.Kind() == reflect.Pointer || fType.Kind() == reflect.Array {
				fType = fType.Elem()
			}
			if fType.Kind() != reflect.Struct {
				return nil, fmt.Errorf("body field must be an struct or array of struct, got %s", field.Type.String())
			}
			bodyType = field.Type
			hasBody = true
		}
		if len(tagQuery)+len(tagJson)+len(tagPath) == 0 {
			hasBody = true
		}
		hasValidation = hasValidation || len(field.Tag.Get("validate")) > 0
	}
	if bodyType == nil && hasBody {
		bodyType = t
	}
	return &Generator[T]{
		sourceFile:    sourceFile,
		pathFields:    pathFields,
		queryFields:   queryFields,
		jsonFields:    jsonFields,
		bodyType:      bodyType,
		hasValidation: hasValidation,
	}, nil
}
