package attributor

import (
	"fmt"
	"reflect"

	"github.com/guionardo/go/pkg/set"
)

type (
	Attributor interface {
		Create(field reflect.StructField, receptor string, value string, args ...any) string
		Imports() []string
	}
	attributor[T any] struct {
		receptor string
		parsers  map[string]fieldParser
	}
	fieldParser struct {
		field  reflect.StructField
		parser Parser
	}
)

func New[T any]() Attributor {

	t := reflect.TypeFor[T]()
	if t.Kind() == reflect.Pointer {
		t = t.Elem()
	}
	attr, ok := attributors[t]
	if !ok {
		parsers := make(map[string]fieldParser)
		for i := range t.NumField() {
			field := t.Field(i)
			if field.IsExported() {
				parsers[field.Name] = fieldParser{
					field,
					NewParserFromType(field.Type),
				}
			}
		}
		a := &attributor[T]{
			receptor: "h", // for handler
			parsers:  parsers,
		}

		attr = a
		attributors[t] = a
	}
	return attr
}

var attributors = make(map[reflect.Type]Attributor)

func (a *attributor[T]) Create(field reflect.StructField, receptor string, value string, args ...any) string {
	parser, ok := a.parsers[field.Name]
	if ok {
		fmtVal := fmt.Sprintf(value, args...)
		return parser.parser.Code(a.receptor, field.Name, fmtVal)
	}
	t := reflect.TypeFor[T]()
	return fmt.Sprintf("// Create field %s.%s.%s failed - no parser", t.PkgPath(), t.Name(), field.Name)
}

func (a *attributor[T]) Imports() []string {
	imports := set.Set[string]{}
	for _, parser := range a.parsers {
		imports.AddMultiple(parser.parser.Imports()...)

	}
	return imports.ToArray()
}
