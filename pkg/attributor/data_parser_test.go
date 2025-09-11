package attributor_test

import (
	"testing"
	"time"

	"github.com/guionardo/go-router/pkg/attributor"
	"github.com/stretchr/testify/assert"
)

func TestNewParserNumbers(t *testing.T) {
	t.Run("int", func(t *testing.T) {
		p := attributor.NewParser[int]()
		assert.Equal(t, []string{"github.com/guionardo/go-router/pkg/tools"}, p.Imports())
		assert.Equal(t, "func () error { return tools.ParseInt(\"Id\",\"1\",&h.Id)}",
			p.Code("h", "Id", "\"1\""))
		p2 := attributor.NewParser[int]()
		assert.Same(t, p2, p)
	})
	t.Run("uint", func(t *testing.T) {
		p := attributor.NewParser[uint]()
		assert.Equal(t, []string{"github.com/guionardo/go-router/pkg/tools"}, p.Imports())
		assert.Equal(t, "func () error { return tools.ParseUInt(\"Id\",\"1\",&h.Id)}",
			p.Code("h", "Id", "\"1\""))
	})
	t.Run("int8", func(t *testing.T) {
		p := attributor.NewParser[int8]()
		assert.Equal(t, []string{"github.com/guionardo/go-router/pkg/tools"}, p.Imports())
		assert.Equal(t, "func () error { return tools.ParseInt8(\"Id\",\"1\",&h.Id)}",
			p.Code("h", "Id", "\"1\""))
	})
	t.Run("uint8", func(t *testing.T) {
		p := attributor.NewParser[uint8]()
		assert.Equal(t, []string{"github.com/guionardo/go-router/pkg/tools"}, p.Imports())
		assert.Equal(t, "func () error { return tools.ParseUInt8(\"Id\",\"1\",&h.Id)}",
			p.Code("h", "Id", "\"1\""))
	})
	t.Run("int16", func(t *testing.T) {
		p := attributor.NewParser[int16]()
		assert.Equal(t, []string{"github.com/guionardo/go-router/pkg/tools"}, p.Imports())
		assert.Equal(t, "func () error { return tools.ParseInt16(\"Id\",\"1\",&h.Id)}",
			p.Code("h", "Id", "\"1\""))
	})
	t.Run("uint16", func(t *testing.T) {
		p := attributor.NewParser[uint16]()
		assert.Equal(t, []string{"github.com/guionardo/go-router/pkg/tools"}, p.Imports())
		assert.Equal(t, "func () error { return tools.ParseUInt16(\"Id\",\"1\",&h.Id)}",
			p.Code("h", "Id", "\"1\""))
	})
	t.Run("int32", func(t *testing.T) {
		p := attributor.NewParser[int32]()
		assert.Equal(t, []string{"github.com/guionardo/go-router/pkg/tools"}, p.Imports())
		assert.Equal(t, "func () error { return tools.ParseInt32(\"Id\",\"1\",&h.Id)}",
			p.Code("h", "Id", "\"1\""))
	})
	t.Run("uint32", func(t *testing.T) {
		p := attributor.NewParser[uint32]()
		assert.Equal(t, []string{"github.com/guionardo/go-router/pkg/tools"}, p.Imports())
		assert.Equal(t, "func () error { return tools.ParseUInt32(\"Id\",\"1\",&h.Id)}",
			p.Code("h", "Id", "\"1\""))
	})
	t.Run("int64", func(t *testing.T) {
		p := attributor.NewParser[int64]()
		assert.Equal(t, []string{"github.com/guionardo/go-router/pkg/tools"}, p.Imports())
		assert.Equal(t, "func () error { return tools.ParseInt64(\"Id\",\"1\",&h.Id)}",
			p.Code("h", "Id", "\"1\""))
	})
	t.Run("uint64", func(t *testing.T) {
		p := attributor.NewParser[uint64]()
		assert.Equal(t, []string{"github.com/guionardo/go-router/pkg/tools"}, p.Imports())
		assert.Equal(t, "func () error { return tools.ParseUInt64(\"Id\",\"1\",&h.Id)}",
			p.Code("h", "Id", "\"1\""))
	})
	t.Run("float32", func(t *testing.T) {
		p := attributor.NewParser[float32]()
		assert.Equal(t, []string{"github.com/guionardo/go-router/pkg/tools"}, p.Imports())
		assert.Equal(t, "func () error { return tools.ParseFloat32(\"Id\",\"1\",&h.Id)}",
			p.Code("h", "Id", "\"1\""))
	})
	t.Run("float64", func(t *testing.T) {
		p := attributor.NewParser[float64]()
		assert.Equal(t, []string{"github.com/guionardo/go-router/pkg/tools"}, p.Imports())
		assert.Equal(t, "func () error { return tools.ParseFloat64(\"Id\",\"1\",&h.Id)}",
			p.Code("h", "Id", "\"1\""))
	})

}

func TestNonNumbers(t *testing.T) {
	t.Run("bool", func(t *testing.T) {
		p := attributor.NewParser[bool]()
		assert.Equal(t, []string{"github.com/guionardo/go-router/pkg/tools"}, p.Imports())
		assert.Equal(t, "func () error { return tools.ParseBool(\"Enabled\",\"1\",&h.Enabled)}",
			p.Code("h", "Enabled", "\"1\""))
	})
	t.Run("time_Time", func(t *testing.T) {
		p := attributor.NewParser[time.Time]()
		assert.Equal(t, []string{"github.com/guionardo/go-router/pkg/tools"}, p.Imports())
		assert.Equal(t, "func () error { return tools.ParseTime(\"When\",\"2020-01-01 00:00:00\",&h.When)}",
			p.Code("h", "When", `"2020-01-01 00:00:00"`))
	})
	t.Run("time_Duration", func(t *testing.T) {
		p := attributor.NewParser[time.Duration]()
		assert.Equal(t, []string{"github.com/guionardo/go-router/pkg/tools"}, p.Imports())
		assert.Equal(t, "func () error { return tools.ParseDuration(\"Interval\",\"1h00\",&h.Interval)}",
			p.Code("h", "Interval", `"1h00"`))
	})
	t.Run("string", func(t *testing.T) {
		p := attributor.NewParser[string]()
		assert.Empty(t, p.Imports())
		assert.Equal(t, "func () error { h.Name = \"Guionardo\"\n\treturn nil }",
			p.Code("h", "Name", `"Guionardo"`))

	})
}
