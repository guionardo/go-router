package inspect

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsStruct(t *testing.T) {
	t.Run("given_a_string_should_return_false", func(t *testing.T) {
		assert.False(t, IsStruct[string]())
	})
	t.Run("given_a_struct_should_return_true", func(t *testing.T) {
		type sample struct {
			name string
		}
		assert.True(t, IsStruct[sample]())
		assert.True(t, IsStruct[*sample]())
	})
}

func TestIsArrayOfStruct(t *testing.T) {
	t.Run("given_a_string_should_return_false", func(t *testing.T) {
		assert.False(t, IsArrayOfStruct[string]())
	})
	t.Run("given_a_bytes_should_return_false", func(t *testing.T) {
		assert.False(t, IsArrayOfStruct[[]byte]())
	})
	t.Run("given_a_slice_array_should_return_true", func(t *testing.T) {
		assert.True(t, IsArrayOfStruct[[]struct{ name string }]())
	})
}
