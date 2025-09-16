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
