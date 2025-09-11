// Structs for development
// RequestStruct, ResponseStruct

package structs

import (
	"time"

	"github.com/go-playground/validator/v10"
)

// Request Sample Structure
//
//easyjson:json
type RequestStruct struct {
	Id         int           `path:"id"`    // Identificator
	Name       string        `query:"name"` // User name
	APIKey     string        `header:"X-API-KEY"`
	When       time.Time     `query:"when"`
	TTL        time.Duration `header:"X-TTL"`
	Operation  uint          `path:"operation" validate:"required"`
	Enabled    bool          `query:"enabled"`
	Value      float64       `header:"value"`
	NumberByte byte          `header:"number_byte"`
	NumberBig  uint64        `query:"number_big"`
	// Body       *Body         `body:"body"`
	BodyField []byte `body:"body_field"`
}

type Body struct {
	Name string
}

//easyjson:json
type ResponseStruct struct {
	Success bool
}

func (rs *RequestStruct) Validate() error {
	return nil
}

func (rs *RequestStruct) GetValidator() *validator.Validate {
	return validator.New(validator.WithRequiredStructEnabled())
}
