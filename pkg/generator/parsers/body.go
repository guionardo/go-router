package parsers

import (
	"log"
	"reflect"
)

type Body[T any] struct {
	Base[T]
	bodyField      reflect.StructField
	bodyFieldIndex int
	bodyType       reflect.Type
}

func NewBody[T any]() *Body[T] {
	p := &Body[T]{bodyFieldIndex: -1}

	p.readFields("body")
	if p.hasBodyField() {
		return p
	}

	p.readFields("json")
	if len(p.fields) > 0 {
		p.bodyType = p.t
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
	b.bodyField = *bodyField
	b.bodyFieldIndex = index
	b.bodyType = bodyField.Type
	b.fillImports(true)
	return true
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
