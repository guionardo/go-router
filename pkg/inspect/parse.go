package inspect

import (
	"errors"
	"net/http"
)

func (s *InspectStruct[T, R]) Parse(r *http.Request) (payload *T, err error) {
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

func (s *InspectStruct[T, R]) customValidate(err error, payload *T) error {
	var is any = payload
	if cv, ok := is.(customValidator); ok {
		return errors.Join(err, cv.Validate())
	}
	return err

}
func (s *InspectStruct[T, R]) postParse(err error, payload *T) error {
	var is any = payload
	if pp, ok := is.(postParser); ok {
		return pp.PostParse(err)
	}
	return err
}
