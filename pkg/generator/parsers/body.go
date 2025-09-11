package parsers

import (
	"fmt"
	"log"
	"reflect"
	"strings"

	reflections "github.com/guionardo/go-router/pkg/reflect"
)

type Body[T any] struct {
	Base[T]
	bodyField      *reflect.StructField
	bodyFieldIndex int
	bodyType       reflect.Type
	requestType    *reflections.Type[T]
	bodyFunc       string
}

func NewBody[T any]() *Body[T] {
	p := &Body[T]{
		bodyFieldIndex: -1,
		requestType:    reflections.New[T](),
	}

	p.readFields("body")
	if p.hasBodyField() {
		p.parseBodyFunc()
		return p
	}

	p.readFields("json")
	if len(p.fields) > 0 {
		p.bodyType = p.t
		p.bodyFunc = p.parseBodyRequestFunc()
	}
	return p
}

func (b *Body[T]) hasBodyField() bool {
	var (
		bodyField *reflect.StructField
		index     int
	)
	for i, field := range b.fields {
		if index <= i {
			bodyField = &field
		}
	}
	if bodyField == nil {
		return false
	}
	if len(b.fields) > 1 {
		log.Printf("[WARN] struct %s.%s has more than one `body` field. Will use the field %s", b.t.PkgPath(), b.t.Name(), bodyField.Name)
	}
	b.bodyField = bodyField
	b.bodyFieldIndex = index
	b.bodyType = bodyField.Type
	b.fillImports(true)
	return true
}

func (b *Body[T]) parseBodyFunc() {
	if b.bodyField != nil {
		b.bodyFunc = b.parseBodyFieldFunc()
		return
	}
	b.bodyFunc = b.parseBodyRequestFunc()
}

func (b *Body[T]) ParseBodyFunc() string {
	return b.bodyFunc
}

func (b *Body[T]) parseBodyFieldFunc() string {
	bt := b.bodyField.Type
	// isPointer := false
	if bt.Kind() == reflect.Pointer {
		// isPointer = true
		bt = bt.Elem()
	}

	if bt.Kind() == reflect.String {
		b.addImport("io")
		return fmt.Sprintf(`func (h *%s) ParseBody(r *http.Request) (err error) {
body,err:=io.ReadAll(r.Body)
if err==nil{ h.%s = string(body) }
return err
}
`, b.t.Name(), b.bodyField.Name)
	} else if bt.Kind() == reflect.Slice {
		bt = bt.Elem()
		if bt.Kind() == reflect.Uint8 {
			b.addImport("io")
			return fmt.Sprintf(`func (h *%s) ParseBody(r *http.Request) (err error) {
body,err:=io.ReadAll(r.Body)
if err==nil{h.%s = body}
return err
}
`, b.t.Name(), b.bodyField.Name)
		}
	}

	// TODO: Implementar tratamento se o tipo for string ou bytes
	bodyType := reflections.NewFromType[T](b.bodyType)
	bodyTypeName := bodyType.Type.Name()
	if bodyType.PackageName != b.requestType.PackageName {
		bodyTypeName = b.requestType.PackageName + "." + bodyTypeName
	}
	var pointerStr string
	if b.bodyField.Type.Kind() == reflect.Pointer {
		pointerStr = "&"
	}
	b.addImport("encoding/json")
	return fmt.Sprintf(`func (h *%s) ParseBody(r *http.Request) (err error) {
	var body %s
	if err = json.NewDecoder(r.Body).Decode(&body);err==nil { h.%s = %sbody	}
	return err
	}`, b.t.Name(), bodyTypeName, b.bodyField.Name, pointerStr)
}

func (b *Body[T]) parseBodyRequestFunc() string {
	b.addImport("encoding/json")
	attribs := make([]string, 0, len(b.fields))
	for _, f := range b.fields {
		attribs = append(attribs, fmt.Sprintf("h.%s = body.%s", f.Name, f.Name))
	}
	return fmt.Sprintf(`func (h *%s) ParseBody(r *http.Request) (err error) {
	var body %s
	if err = json.NewDecoder(r.Body).Decode(&body);err!=nil {return err}
%s
	return err
	}`, b.t.Name(), b.t.Name(), strings.Join(attribs, "\n"))
}

// b.tag = tag
// 	b.fields = make(map[int]reflect.StructField)
// 	t := reflect.TypeFor[T]()
// 	if t.Kind() == reflect.Pointer {
// 		t = t.Elem()
// 	}
// 	b.t = t
// 	for i := range b.t.NumField() {
// 		field := b.t.Field(i)
// 		if field.IsExported() && len(field.Tag.Get(tag)) > 0 {
// 			b.fields[i] = field
// 		}
// 	}

/*
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
*/
