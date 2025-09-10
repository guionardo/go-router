package generator

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestToSnakeCase(t *testing.T) {
	assert.Equal(t, "to_snake_case", toSnakeCase("ToSnakeCase"))
}
