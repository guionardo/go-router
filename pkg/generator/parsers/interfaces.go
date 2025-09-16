package parsers

import (
	"github.com/go-playground/validator/v10"
)

type (
	Validator interface {
		Validate() error
	}
	ValidatorGetter interface {
		GetValidator() *validator.Validate
	}
)
