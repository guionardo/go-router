package endpoint

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"net/http"
	"reflect"

	"github.com/go-playground/validator"
	reflections "github.com/guionardo/go-router/pkg/reflect"
	"gopkg.in/yaml.v3"
)

func (s *Endpoint[T, R]) Parse(r *http.Request) (payload *T, err error) {
	payload = new(T)
	for _, f := range []func(*http.Request, *T) error{
		s.parseBodyFunc,
		s.parsePathFunc,
		s.parseQueriesFunc,
		s.parseHeadersFunc,
		s.validateFunc} {
		if errF := f(r, payload); errF != nil {
			err = errors.Join(err, errF)
		}
	}

	err = s.customValidate(err, payload)
	err = s.postParse(err, payload)
	return payload, err
}

func (s *Endpoint[T, R]) customValidate(err error, payload *T) error {
	var is any = payload
	if cv, ok := is.(customValidator); ok {
		return errors.Join(err, cv.Validate())
	}
	return err

}
func (s *Endpoint[T, R]) postParse(err error, payload *T) error {
	var is any = payload
	if pp, ok := is.(postParser); ok {
		return pp.PostParse(err)
	}
	return err
}

func (s *Endpoint[T, R]) parseBody(r *http.Request, payload *T) (err error) {
	switch s.bodyParseType {
	case BodyNo:
		return nil
	case BodyBytes:
		return s.parseBodyBytes(r, payload)
	case BodyString:
		return s.parseBodyString(r, payload)
	case BodyJSON:
		return s.parseBodyJSON(r, payload)
	}

	switch r.Header.Get("Content-Type") {
	case ContentJSON:
		err = json.NewDecoder(r.Body).Decode(payload)
	case ContentYAML:
		err = yaml.NewDecoder(r.Body).Decode(payload)
	case ContentXML:
		err = xml.NewDecoder(r.Body).Decode(payload)
	default:
		body, err := io.ReadAll(r.Body)
		if err == nil {
			s.raw = body
		}
	}
	if err != nil {
		s.raw = nil
	}
	return err
}

func (s *Endpoint[T, R]) parseBodyBytes(r *http.Request, payload *T) (err error) {
	v := reflect.ValueOf(payload).Elem()
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

func (s *Endpoint[T, R]) parseBodyString(r *http.Request, payload *T) (err error) {
	v := reflect.ValueOf(payload).Elem()
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

func (s *Endpoint[T, R]) parseBodyJSON(r *http.Request, payload *T) (err error) {
	v := reflect.ValueOf(payload).Elem()
	field := v.FieldByName(s.bodyField)
	if field.IsValid() && field.CanSet() {
		s := reflect.New(field.Type()).Interface()
		if err = json.NewDecoder(r.Body).Decode(s); err == nil {
			vs := reflect.ValueOf(s).Elem()
			field.Set(vs)

			// var body []byte
			// if body, err = io.ReadAll(r.Body); err == nil {
			// 	r.Body.Close()
			// 	if err = json.Unmarshal(body, s); err == nil {
			// 		vs := reflect.ValueOf(s).Elem()
			// 		field.Set(vs)
			// 	}

		} else {
			err = fmt.Errorf("field '%s' is not valid or cannot be set", field.Type().Name())
		}
	}

	return err
}

func (s *Endpoint[T, R]) parsePath(r *http.Request, payload *T) (err error) {
	for pathParam, index := range s.paths {
		if value := r.PathValue(pathParam); len(value) > 0 {
			err = s.inject(payload, index, value)
		} else {
			err = fmt.Errorf("expected value for '%s' path parameter", pathParam)
		}
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *Endpoint[T, R]) parseQuery(r *http.Request, payload *T) (err error) {
	for key, index := range s.queries {
		if value := r.URL.Query().Get(key); len(value) > 0 {
			if err = s.inject(payload, index, value); err != nil {
				return err
			}
		}
	}
	return nil
}

func (s *Endpoint[T, R]) parseHeaders(r *http.Request, payload *T) (err error) {
	for key, index := range s.headers {
		if value := r.Header.Get(key); len(value) > 0 {
			if err = s.inject(payload, index, value); err != nil {
				return err
			}
		}
	}
	return nil
}

func (s *Endpoint[T, R]) inject(payload *T, index int, value any) error {
	rv := reflect.ValueOf(payload).Elem() // payload is always a pointer
	f := rv.Field(index)
	nv, err := reflections.NewValue(value, f)
	if err == nil {
		f.Set(nv)
	}
	return err
}

func (s *Endpoint[T, R]) validate(r *http.Request, payload *T) error {
	var (
		si  any = payload
		val *validator.Validate
	)
	if s.isSetupValidator {
		val = si.(setupValidator).NewValidator()
	} else {
		val = validator.New()
	}

	if val == nil {
		return fmt.Errorf("struct '%s' implements the method NewValidator() *validator.Validate but returned nil", s.reqType.Name())
	}
	return val.StructCtx(r.Context(), payload)
}
