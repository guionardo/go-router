package inspect

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"reflect"
	"runtime"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/go-playground/validator"
	"github.com/guionardo/go-router/pkg/config"
	"github.com/guionardo/go-router/pkg/logging"
	"github.com/guionardo/go-router/pkg/path_params"
	"gopkg.in/yaml.v3"
)

type (
	InspectStruct[T any, R any] struct {
		path          string
		pathParams    *path_params.PathParams
		useValidation bool
		payload       *T
		raw           []byte
		paths         map[string]int
		queries       map[string]int
		headers       map[string]int
		bodyField     string
		bodyParseType bodyParseType
		typeName      string

		parseBodyFunc    func(*http.Request) error
		parsePathFunc    func(*http.Request) error
		parseQueriesFunc func(*http.Request) error
		parseHeadersFunc func(*http.Request) error
		validateFunc     func(*http.Request) error

		isCustomValidator bool
		isPostParser      bool
		isSetupValidator  bool
		initialized       bool

		responser       Responser[T, R]
		customResponser CustomResponser[T, R]
		handlerFunc     func(http.ResponseWriter, *http.Request, *T)
		handlerName     string
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

func New[T, R any](path string) (is *InspectStruct[T, R], err error) {
	var t reflect.Type
	if is, t, err = inspectStructGet[T, R](); err != nil {
		return
	}
	if is == nil || !is.initialized {
		is, err = buildStruct[T, R](t, path)
		if err == nil {
			poolSet(t, is)
		}
	}
	return is, err
}

func (s *InspectStruct[T, R]) GetPath() string {
	return s.path
}
func (s *InspectStruct[T, R]) Parse(r *http.Request) (*T, error) {
	s.payload = new(T)
	var err error
	for _, f := range []func(*http.Request) error{
		s.parseBodyFunc,
		s.parsePathFunc,
		s.parseQueriesFunc,
		s.parseHeadersFunc,
		s.validateFunc} {
		if err = f(r); err != nil {
			break
		}
	}

	err = s.customValidate()
	err = s.postParse(err)
	return s.payload, err
}

func (s *InspectStruct[T, R]) customValidate() error {
	var is any = s.payload
	if cv, ok := is.(customValidator); ok {
		return cv.Validate()
	}
	return nil

}
func (s *InspectStruct[T, R]) postParse(err error) error {
	var is any = s.payload
	if pp, ok := is.(postParser); ok {
		return pp.PostParse(err)
	}
	return err
}

func reqFunc(condition bool, rf func(*http.Request) error) func(*http.Request) error {
	if !condition {
		return func(*http.Request) error {
			return nil
		}
	}
	return rf
}

func buildStruct[T, R any](t reflect.Type, path string) (*InspectStruct[T, R], error) {
	is := &InspectStruct[T, R]{
		typeName: t.Name(),
		path:     path,
	}
	var si any = new(T)

	if cr, isCustomResponser := si.(CustomResponser[T, R]); isCustomResponser {
		is.handlerName = runtime.FuncForPC(reflect.ValueOf(cr.Handle).Pointer()).Name()
		is.handlerFunc = is.handleCustom
	} else if sr, isSimpleResponser := si.(Responser[T, R]); isSimpleResponser {
		is.handlerName = runtime.FuncForPC(reflect.ValueOf(sr.Handle).Pointer()).Name()
		is.handlerFunc = is.handleSimple
	} else {
		tcr := reflect.TypeFor[CustomResponser[T, R]]()
		tsr := reflect.TypeFor[Responser[T, R]]()
		return nil, fmt.Errorf("type %s should implements interfaces %s or %s", t.Name(), tcr.Name(), tsr.Name())
	}

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
	pathParams, err := path_params.NewPathParams(path)
	if err == nil {
		err = checkPathParamFields(path, paths, pathParams, t)
	}
	if err != nil {
		return nil, err
	}
	is.pathParams = pathParams
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

func checkPathParamFields(path string, paths map[string]int, pathParams *path_params.PathParams, t reflect.Type) error {
	for _, pathParam := range pathParams.ParamNames() {
		if _, ok := paths[pathParam]; !ok {
			return fmt.Errorf("expected field with 'path' tag in %s type", t.Name())
		}
	}
	for fieldParam := range paths {
		if !slices.Contains(pathParams.ParamNames(), fieldParam) {
			return fmt.Errorf("expected path param '%s' in endpoint path '%s'", fieldParam, path)
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

func (s *InspectStruct[T, R]) parseBody(r *http.Request) (err error) {
	switch s.bodyParseType {
	case BodyNo:
		return nil
	case BodyBytes:
		return s.parseBodyBytes(r)
	case BodyString:
		return s.parseBodyString(r)
	case BodyJSON:
		return s.parseBodyJSON(r)
	}

	switch r.Header.Get("Content-Type") {
	case ContentJSON:
		err = json.NewDecoder(r.Body).Decode(s.payload)
	case ContentYAML:
		err = yaml.NewDecoder(r.Body).Decode(s.payload)
	case ContentXML:
		err = xml.NewDecoder(r.Body).Decode(s.payload)
	default:
		body, err := io.ReadAll(r.Body)
		if err == nil {
			s.raw = body
		}
	}
	if err != nil {
		s.raw = nil
		s.payload = nil
	}
	return err
}

func (s *InspectStruct[T, R]) parseBodyBytes(r *http.Request) (err error) {
	v := reflect.ValueOf(s.payload).Elem()
	field := v.FieldByName(s.bodyField)
	if field.IsValid() && field.CanSet() {
		var body []byte
		if body, err = io.ReadAll(r.Body); err == nil {
			r.Body.Close()
			field.SetBytes(body)
		} else {
			err = fmt.Errorf("field '%s' is not valid or cannot be set", field.Type().Name())
		}
	}

	return err
}

func (s *InspectStruct[T, R]) parseBodyString(r *http.Request) (err error) {
	v := reflect.ValueOf(s.payload).Elem()
	field := v.FieldByName(s.bodyField)
	if field.IsValid() && field.CanSet() {
		var body []byte
		if body, err = io.ReadAll(r.Body); err == nil {
			r.Body.Close()
			field.SetString(string(body))
		} else {
			err = fmt.Errorf("field '%s' is not valid or cannot be set", field.Type().Name())
		}
	}

	return err
}

func (s *InspectStruct[T, R]) parseBodyJSON(r *http.Request) (err error) {
	v := reflect.ValueOf(s.payload).Elem()
	field := v.FieldByName(s.bodyField)
	if field.IsValid() && field.CanSet() {
		s := reflect.New(field.Type()).Interface()
		var body []byte
		if body, err = io.ReadAll(r.Body); err == nil {
			r.Body.Close()
			if err = json.Unmarshal(body, s); err == nil {
				vs := reflect.ValueOf(s).Elem()
				field.Set(vs)
			}

		} else {
			err = fmt.Errorf("field '%s' is not valid or cannot be set", field.Type().Name())
		}
	}

	return err
}

func (s *InspectStruct[T, R]) parsePath(r *http.Request) (err error) {
	if err := s.pathParams.Match(r.URL); err != nil {
		return err
	}
	for key, value := range s.pathParams.GetAll() {
		index := s.paths[key]
		if err := s.inject(index, value); err != nil {
			return err
		}
	}
	return nil
}

func (s *InspectStruct[T, R]) parseQuery(r *http.Request) (err error) {
	for key, index := range s.queries {
		if value := r.URL.Query().Get(key); len(value) > 0 {
			if err = s.inject(index, value); err != nil {
				return err
			}
		}
	}
	return nil
}

func (s *InspectStruct[T, R]) parseHeaders(r *http.Request) (err error) {
	for key, index := range s.headers {
		if value := r.Header.Get(key); len(value) > 0 {
			if err = s.inject(index, value); err != nil {
				return err
			}
		}
	}
	return nil
}

func (s *InspectStruct[T, R]) inject(index int, value any) error {
	rv := reflect.ValueOf(s.payload).Elem() // payload is always a pointer
	f := rv.Field(index)
	nv, err := newValue(value, f)
	if err == nil {
		f.Set(nv)
	}
	return err
}

func newValue(value any, field reflect.Value) (reflect.Value, error) {
	kind := field.Kind()
	switch v := value.(type) {
	case string:
		switch kind {
		case reflect.String:
			return reflect.ValueOf(v), nil
		case reflect.Bool:
			bv, err := strconv.ParseBool(v)
			return reflect.ValueOf(bv), err
		case reflect.Int:
			iv, err := strconv.Atoi(v)
			return reflect.ValueOf(iv), err
		case reflect.Float64:
			fv, err := strconv.ParseFloat(v, 64)
			return reflect.ValueOf(fv), err
		case reflect.Float32:
			fv, err := strconv.ParseFloat(v, 32)
			return reflect.ValueOf(float32(fv)), err
		case reflect.Struct:
			if _, ok := field.Interface().(time.Time); ok {
				tv, err := strToTime(v)
				return reflect.ValueOf(tv), err
			}
		}

	case bool:
		switch kind {
		case reflect.Bool:
			return reflect.ValueOf(v), nil
		case reflect.String:
			return reflect.ValueOf(boolValue(v, "true", "false")), nil
		case reflect.Int:
			return reflect.ValueOf(boolValue(v, 1, 0)), nil
		}
	case []byte:
		switch kind {
		case reflect.Array:
			return reflect.ValueOf(v), nil
		}
	}

	return reflect.ValueOf(nil), fmt.Errorf("unexpected convertion %s to %s", reflect.TypeOf(value).Name(), kind.String())
}

func (s *InspectStruct[T, R]) validate(r *http.Request) error {
	var (
		si  any = s.payload
		val *validator.Validate
	)
	if s.isSetupValidator {
		val = si.(setupValidator).NewValidator()
	} else {
		val = validator.New()
	}

	if val == nil {
		return fmt.Errorf("struct '%s' implements the method NewValidator() *validator.Validate but returned nil", s.typeName)
	}
	return val.StructCtx(r.Context(), s.payload)
}

func (s *InspectStruct[T, R]) Handle(w http.ResponseWriter, r *http.Request) {
	s.handlerFunc(w, r, s.payload)
}

func (s *InspectStruct[T, R]) handleCustom(w http.ResponseWriter, r *http.Request, payload *T) {
	if err := s.customResponser.Handle(w, r, s.payload); err != nil {
		logging.Get().Warn("",
			slog.String("path", r.URL.Path),
			slog.String("method", s.handlerName),
			slog.Any("error", err),
		)
	}
}
func (s *InspectStruct[T, R]) handleSimple(w http.ResponseWriter, r *http.Request, payload *T) {
	// Simple handler
	response, status, err := s.responser.Handle(r, s.payload)
	if responseErr, ok := err.(responseError); ok && status <= 0 {
		status = responseErr.StatusCode()
	}

	if status <= 0 {
		status = http.StatusBadGateway
	}

	w.WriteHeader(status)
	if config.DevelopmentMode {
		w.Header().Add("X-Handler-Name", s.handlerName)
	}

	var responseBody []byte
	if err != nil {
		responseBody, _ = json.Marshal(err)
		w.Write(responseBody)
	} else if response != nil {
		responseBody, _ = json.Marshal(response)
		w.Write(responseBody)
	}
}
