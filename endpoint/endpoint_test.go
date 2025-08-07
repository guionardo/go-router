package endpoint

import (
	"testing"

	reflections "github.com/guionardo/go-router/pkg/reflect"
	"github.com/stretchr/testify/assert"
)

func TestIsStruct(t *testing.T) {
	t.Run("given_a_string_should_return_false", func(t *testing.T) {
		assert.False(t, reflections.IsStruct[string]())
	})
	t.Run("given_a_struct_should_return_true", func(t *testing.T) {
		type sample struct {
			name string
		}
		assert.True(t, reflections.IsStruct[sample]())
		assert.True(t, reflections.IsStruct[*sample]())
	})
}

func TestIsArrayOfStruct(t *testing.T) {
	t.Run("given_a_string_should_return_false", func(t *testing.T) {
		assert.False(t, reflections.IsArrayOfStruct[string]())
	})
	t.Run("given_a_bytes_should_return_false", func(t *testing.T) {
		assert.False(t, reflections.IsArrayOfStruct[[]byte]())
	})
	t.Run("given_a_slice_array_should_return_true", func(t *testing.T) {
		assert.True(t, reflections.IsArrayOfStruct[[]struct{ name string }]())
	})
}
