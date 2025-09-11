package generator

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestToSnakeCase(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"ToSnakeCase", "to_snake_case"},
		{"toSnakeCase", "to_snake_case"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := toSnakeCase(tt.input)
			assert.Equal(t, tt.expected, got)
		})
	}
}
