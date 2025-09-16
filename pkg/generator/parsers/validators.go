package parsers

import "reflect"

type Validators[T any] struct {
	Base[T]
	HasValidations bool
}

func NewValidators[T any]() *Validators[T] {
	p := &Validators[T]{}
	p.readFields("validate")
	if len(p.fields) > 0 {
		if !reflect.TypeFor[T]().Implements(reflect.TypeFor[Validator]()) &&
			!reflect.TypeFor[T]().Implements(reflect.TypeFor[ValidatorGetter]()) {
			p.Imports.Add("github.com/go-playground/validator/v10")
		}
		p.HasValidations = true
	}
	p.fillImports(true)

	return p
}
