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
		assert.Equal(t, []string{"strconv"}, p.Imports())
		assert.Equal(t, `if value,err:=strconv.ParseInt("1",0,strconv.IntSize); err!=nil { return err } else { h.Id = value }`,
			p.Code("h", "Id", "\"1\""))
		p2 := attributor.NewParser[int]()
		assert.Same(t, p2, p)
	})
	t.Run("uint", func(t *testing.T) {
		p := attributor.NewParser[uint]()
		assert.Equal(t, []string{"strconv"}, p.Imports())
		assert.Equal(t, `if value,err:=strconv.ParseUint("1",0,strconv.IntSize); err!=nil { return err } else { h.Id = value }`,
			p.Code("h", "Id", "\"1\""))
	})
	t.Run("int8", func(t *testing.T) {
		p := attributor.NewParser[int8]()
		assert.Equal(t, []string{"strconv"}, p.Imports())
		assert.Equal(t, `if value,err:=strconv.ParseInt("1",0,8); err!=nil { return err } else { h.Id = int8(value) }`,
			p.Code("h", "Id", "\"1\""))
	})
	t.Run("uint8", func(t *testing.T) {
		p := attributor.NewParser[uint8]()
		assert.Equal(t, []string{"strconv"}, p.Imports())
		assert.Equal(t, `if value,err:=strconv.ParseUint("1",0,8); err!=nil { return err } else { h.Id = uint8(value) }`,
			p.Code("h", "Id", "\"1\""))
	})
	t.Run("int16", func(t *testing.T) {
		p := attributor.NewParser[int16]()
		assert.Equal(t, []string{"strconv"}, p.Imports())
		assert.Equal(t, `if value,err:=strconv.ParseInt("1",0,16); err!=nil { return err } else { h.Id = int16(value) }`,
			p.Code("h", "Id", "\"1\""))
	})
	t.Run("uint16", func(t *testing.T) {
		p := attributor.NewParser[uint16]()
		assert.Equal(t, []string{"strconv"}, p.Imports())
		assert.Equal(t, `if value,err:=strconv.ParseUint("1",0,16); err!=nil { return err } else { h.Id = uint16(value) }`,
			p.Code("h", "Id", "\"1\""))
	})
	t.Run("int32", func(t *testing.T) {
		p := attributor.NewParser[int32]()
		assert.Equal(t, []string{"strconv"}, p.Imports())
		assert.Equal(t, `if value,err:=strconv.ParseInt("1",0,32); err!=nil { return err } else { h.Id = int32(value) }`,
			p.Code("h", "Id", "\"1\""))
	})
	t.Run("uint32", func(t *testing.T) {
		p := attributor.NewParser[uint32]()
		assert.Equal(t, []string{"strconv"}, p.Imports())
		assert.Equal(t, `if value,err:=strconv.ParseUint("1",0,32); err!=nil { return err } else { h.Id = uint32(value) }`,
			p.Code("h", "Id", "\"1\""))
	})
	t.Run("int64", func(t *testing.T) {
		p := attributor.NewParser[int64]()
		assert.Equal(t, []string{"strconv"}, p.Imports())
		assert.Equal(t, `if value,err:=strconv.ParseInt("1",0,64); err!=nil { return err } else { h.Id = value }`,
			p.Code("h", "Id", "\"1\""))
	})
	t.Run("uint64", func(t *testing.T) {
		p := attributor.NewParser[uint64]()
		assert.Equal(t, []string{"strconv"}, p.Imports())
		assert.Equal(t, `if value,err:=strconv.ParseUint("1",0,64); err!=nil { return err } else { h.Id = value }`,
			p.Code("h", "Id", "\"1\""))
	})
	t.Run("float32", func(t *testing.T) {
		p := attributor.NewParser[float32]()
		assert.Equal(t, []string{"strconv"}, p.Imports())
		assert.Equal(t, `if value,err:=strconv.ParseFloat("1",32); err!=nil { return err } else { h.Id = float32(value) }`,
			p.Code("h", "Id", "\"1\""))
	})
	t.Run("float64", func(t *testing.T) {
		p := attributor.NewParser[float64]()
		assert.Equal(t, []string{"strconv"}, p.Imports())
		assert.Equal(t, `if value,err:=strconv.ParseFloat("1",64); err!=nil { return err } else { h.Id = value }`,
			p.Code("h", "Id", "\"1\""))
	})

}

func TestNonNumbers(t *testing.T) {
	t.Run("bool", func(t *testing.T) {
		p := attributor.NewParser[bool]()
		assert.Equal(t, []string{"strconv"}, p.Imports())
		assert.Equal(t, `if value,err:=strconv.ParseBool("1"); err!=nil { return err } else { h.Enabled = value }`,
			p.Code("h", "Enabled", "\"1\""))
	})
	t.Run("time_Time", func(t *testing.T) {
		p := attributor.NewParser[time.Time]()
		assert.Equal(t, []string{"github.com/guionardo/go-router/pkg/tools"}, p.Imports())
		assert.Equal(t, `if value,err:=tools.ParseTime("2020-01-01 00:00:00"); err!=nil { return err } else { h.When = value }`,
			p.Code("h", "When", `"2020-01-01 00:00:00"`))
	})
	t.Run("time_Duration", func(t *testing.T) {
		p := attributor.NewParser[time.Duration]()
		assert.Equal(t, []string{"time"}, p.Imports())
		assert.Equal(t, `if value,err:=time.ParseDuration("1h00"); err!=nil { return err } else { h.Interval = value }`,
			p.Code("h", "Interval", `"1h00"`))
	})
	t.Run("string", func(t *testing.T) {
		p := attributor.NewParser[string]()
		assert.Empty(t, p.Imports())
		assert.Equal(t, `h.Name = "Guionardo"`,
			p.Code("h", "Name", `"Guionardo"`))

	})
}
