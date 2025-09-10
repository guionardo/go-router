package endpoint

import (
	"fmt"
	"net/http"
	"reflect"
	"runtime"
	"slices"
	"strings"

	"github.com/guionardo/go-router/pkg/path_params"
)

func buildStruct[T, R any](t reflect.Type, path string) (*Endpoint[T, R], error) {
	is := &Endpoint[T, R]{
		reqType:  t,
		respType: reflect.TypeFor[R](),
		path:     path,
	}
	if err := is.buildResponser(); err != nil {
		return nil, err
	}
	var si any = new(T)

	useValidation := false
	var (
		paths        = make(map[string]int)
		queries      = make(map[string]int)
		headers      = make(map[string]int)
		descriptions = make(map[int]string)
	)
	var bodyField string
	bodyParseType := BodyNo
	for i := range t.NumField() {
		field := t.Field(i)
		if !field.IsExported() {
			continue
		}

		if _, _, err := extractTag(field.Tag, "validate"); err == nil {
			useValidation = true
		}
		if pathTag, _, err := extractTag(field.Tag, "path"); err == nil {
			paths[pathTag] = i
		}
		if queryTag, _, err := extractTag(field.Tag, "query"); err == nil {
			queries[queryTag] = i
		}
		if headerTag, _, err := extractTag(field.Tag, "header"); err == nil {
			headers[headerTag] = i
		}
		if _, _, err := extractTag(field.Tag, "body"); err == nil {
			if len(bodyField) > 0 {
				return nil, fmt.Errorf("field '%s' has tag body, but previous field '%s' also have. Should be only one field with the 'body' tag", field.Name, bodyField)
			}
			ft := field.Type.Kind()
			if ft == reflect.Pointer {
				ft = field.Type.Elem().Kind()
			}
			switch ft {
			case reflect.Slice:
				switch field.Type.Elem().Kind() {
				case reflect.Uint8:
					bodyParseType = BodyBytes

					// []byte
				case reflect.Struct:
					// []struct{}
				case reflect.Map:
					// []map[x]x
				}

			case reflect.String:
				// string
				bodyParseType = BodyString

			case reflect.Struct:
				// struct
				bodyParseType = BodyJSON

			case reflect.Map:
				bodyParseType = BodyJSON
				// map
			}
			// o campo body pode ser string, []byte, []struct{}, struct{}

			bodyField = field.Name
		}
		if descTag, extras, err := extractTag(field.Tag, "description"); err == nil {
			descriptions[i] = strings.Join(append([]string{descTag}, extras...), ",")
		}
	}
	pathParams, err := path_params.GetPathParamsNames(path)
	if err == nil {
		err = checkPathParamFields(path, pathParams, paths, t)
	}

	if err != nil {
		return nil, err
	}

	is.useValidation = useValidation
	is.paths = paths
	is.queries = queries
	is.headers = headers

	_, isCustomValidator := si.(customValidator)
	_, isPostParser := si.(postParser)
	_, isSetupValidator := si.(setupValidator)

	is.isCustomValidator = isCustomValidator
	is.isPostParser = isPostParser
	is.isSetupValidator = isSetupValidator
	is.bodyParseType = bodyParseType
	is.bodyField = bodyField

	is.parseBodyFunc = reqFunc(bodyParseType != BodyNo, is.parseBody)
	is.parsePathFunc = reqFunc(len(paths) > 0, is.parsePath)
	is.parseQueriesFunc = reqFunc(len(queries) > 0, is.parseQuery)
	is.parseHeadersFunc = reqFunc(len(headers) > 0, is.parseHeaders)
	is.validateFunc = reqFunc(useValidation, is.validate)

	is.initialized = true
	return is, nil
}

// reqFunc returns function only if the condition is true, or a noop func if false
func reqFunc[T any](condition bool, rf func(*http.Request, *T) error) func(*http.Request, *T) error {
	if condition {
		return rf
	}
	return func(*http.Request, *T) error {
		return nil
	}
}

func checkPathParamFields(path string, pathParams []string, paths map[string]int, t reflect.Type) error {
	for _, p := range pathParams {
		if _, ok := paths[p]; !ok {
			return fmt.Errorf("expected field with 'path' tag in %s type", t.Name())
		}
	}
	for p := range paths {
		if !slices.Contains(pathParams, p) {
			return fmt.Errorf("expected path param '%s' in endpoint path '%s'", p, path)
		}
	}
	return nil
}

func extractTag(tag reflect.StructTag, name string) (value string, extra []string, err error) {
	valTag := tag.Get(name)
	if len(valTag) == 0 {
		return "", []string{}, errEmptyTag
	}
	w := strings.Split(valTag, ",")
	return w[0], w[1:], nil
}

func (is *Endpoint[T, R]) buildResponser() error {
	if is.handlerFunc != nil {
		return nil
	}
	var s any = new(T)

	if sr, srt := s.(Responser[T, R]); srt {
		is.handlerFunc = is.handleSimple
		is.handlerName = runtime.FuncForPC(reflect.ValueOf(sr.Handle).Pointer()).Name()
		is.responser = sr
	} else if cr, crt := s.(CustomResponser[T, R]); crt {
		is.handlerName = runtime.FuncForPC(reflect.ValueOf(cr.Handle).Pointer()).Name()
		is.handlerFunc = is.handleCustom
		is.customResponser = cr
	} else {
		tcr := reflect.TypeFor[CustomResponser[T, R]]()
		tsr := reflect.TypeFor[Responser[T, R]]()
		return fmt.Errorf("type %s should implements interfaces %s or %s", is.reqType.String(), tcr.Name(), tsr.Name())
	}

	if _, ok := s.(Responser[T, R]); !ok {
		it := reflect.TypeFor[Responser[T, R]]()
		return fmt.Errorf("struct '%s' must implement the interface %s", is.reqType.String(), it.Name())
	}

	return nil
}
