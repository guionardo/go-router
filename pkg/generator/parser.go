package generator

import (
	"iter"
	"reflect"
)

func getFieldsByTag(t reflect.Type, tag string) iter.Seq2[reflect.StructField, string] {
	return func(yield func(reflect.StructField, string) bool) {
		for i := range t.NumField() {
			field := t.Field(i)
			if !field.IsExported() {
				continue
			}
			if tagValue := field.Tag.Get(tag); len(tagValue) > 0 {
				if !yield(field, tagValue) {
					return
				}
			}
		}
	}
}
