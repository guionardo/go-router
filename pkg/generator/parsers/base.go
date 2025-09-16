package parsers

import (
	"iter"
	"reflect"

	"github.com/guionardo/go-router/pkg/attributor"
	"github.com/guionardo/go/pkg/set"
)

type (
	Base[T any] struct {
		t       reflect.Type
		Imports set.Set[string]
		fields  map[int]reflect.StructField
		tag     string
	}
)

func (b *Base[T]) readFields(tag string) {
	b.tag = tag
	b.fields = make(map[int]reflect.StructField)
	t := reflect.TypeFor[T]()
	if t.Kind() == reflect.Pointer {
		t = t.Elem()
	}
	b.t = t
	for i := range b.t.NumField() {
		field := b.t.Field(i)
		if field.IsExported() && len(field.Tag.Get(tag)) > 0 {
			b.fields[i] = field
		}
	}
}

func (b *Base[T]) addImport(i string) {
	if len(b.Imports) == 0 {
		b.Imports = set.New[string]()
	}
	b.Imports.Add(i)
}

func (b *Base[T]) fillImports(inputTypeIsString bool) {
	if len(b.Imports) == 0 {
		b.Imports = set.New[string]()
	}
	for _, field := range b.fields {
		k := field.Type
		if k.Kind() == reflect.Pointer {
			k = k.Elem()
		}
		if inputTypeIsString {
			a := attributor.NewParserFromType(k)
			for _, i := range a.Imports() {
				b.addImport(i)
			}
		}
	}
}

func (b *Base[T]) Fields() iter.Seq2[reflect.StructField, string] {
	return func(yield func(reflect.StructField, string) bool) {
		for _, field := range b.fields {
			if !yield(field, field.Tag.Get(b.tag)) {
				return
			}
		}
	}
}
