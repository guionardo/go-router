package attributor

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

type Sample struct {
	NumericId int
	Order     uint
}

func Test_attributor_Create(t *testing.T) {
	a := New[Sample]()
	assert.NotNil(t, a)
	fields := make([]reflect.StructField, 0)
	ta := reflect.TypeFor[Sample]()

	for i := range ta.NumField() {
		field := ta.Field(i)
		if field.IsExported() {
			fields = append(fields, field)
		}
	}
	assert.NotEmpty(t, fields)

	t.Run("numeric_id", func(t *testing.T) {
		at := a.Create(fields[0], "h", `"1"`)
		assert.Contains(t, at, "h.NumericId = value")
	})
	t.Run("order", func(t *testing.T) {
		at := a.Create(fields[1], "h", `"1"`)
		assert.Contains(t, at, "h.Order = value")
	})

}
