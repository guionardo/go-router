package tools_test

import (
	"testing"

	"github.com/guionardo/go-router/pkg/tools"
	"github.com/stretchr/testify/assert"
)

func TestParseNumbers(t *testing.T) {
	type TS struct {
		Vint   int
		Vint8  int8
		Vint16 int16
		Vint32 int32
		Vint64 int64

		Uint   uint
		Uint8  uint8
		Uint16 uint16
		Uint32 uint32
		Uint64 uint64

		VFloat32 float32
		VFloat64 float64

		Bool bool
	}
	ts := TS{}
	t.Run("int", func(t *testing.T) {
		err := tools.ParseInt("VInt", "10", &ts.Vint)
		assert.NoError(t, err)
		assert.Equal(t, 10, ts.Vint)
	})
	t.Run("int error", func(t *testing.T) {
		err := tools.ParseInt("VInt", "x", &ts.Vint)
		assert.Error(t, err)
		assert.Equal(t, "VInt - strconv.ParseInt: parsing \"x\": invalid syntax", err.Error())
		assert.Equal(t, 0, ts.Vint)
	})

	t.Run("int8", func(t *testing.T) {
		err := tools.ParseInt8("VInt8", "12", &ts.Vint8)
		assert.NoError(t, err)
		assert.Equal(t, int8(12), ts.Vint8)
	})

	t.Run("int16", func(t *testing.T) {
		err := tools.ParseInt16("VInt16", "14", &ts.Vint16)
		assert.NoError(t, err)
		assert.Equal(t, int16(14), ts.Vint16)
	})
	t.Run("int32", func(t *testing.T) {
		err := tools.ParseInt32("VInt32", "15", &ts.Vint32)
		assert.NoError(t, err)
		assert.Equal(t, int32(15), ts.Vint32)
	})

	t.Run("int64", func(t *testing.T) {
		err := tools.ParseInt64("VInt16", "20", &ts.Vint64)
		assert.NoError(t, err)
		assert.Equal(t, int64(20), ts.Vint64)
	})

	t.Run("uint", func(t *testing.T) {
		err := tools.ParseUInt("VInt", "10", &ts.Uint)
		assert.NoError(t, err)
		assert.Equal(t, uint(10), ts.Uint)
	})

	t.Run("uint8", func(t *testing.T) {
		err := tools.ParseUInt8("VInt8", "12", &ts.Uint8)
		assert.NoError(t, err)
		assert.Equal(t, uint8(12), ts.Uint8)
	})

	t.Run("uint16", func(t *testing.T) {
		err := tools.ParseUInt16("VInt16", "14", &ts.Uint16)
		assert.NoError(t, err)
		assert.Equal(t, uint16(14), ts.Uint16)
	})
	t.Run("uint32", func(t *testing.T) {
		err := tools.ParseUInt32("VInt32", "15", &ts.Uint32)
		assert.NoError(t, err)
		assert.Equal(t, uint32(15), ts.Uint32)
	})

	t.Run("uint64", func(t *testing.T) {
		err := tools.ParseUInt64("VInt16", "20", &ts.Uint64)
		assert.NoError(t, err)
		assert.Equal(t, uint64(20), ts.Uint64)
	})

	t.Run("float32", func(t *testing.T) {
		err := tools.ParseFloat32("VFloat32", "12.34", &ts.VFloat32)
		assert.NoError(t, err)
		assert.Equal(t, float32(12.34), ts.VFloat32)
	})
	t.Run("float64", func(t *testing.T) {
		err := tools.ParseFloat64("VFloat64", "12.34", &ts.VFloat64)
		assert.NoError(t, err)
		assert.Equal(t, float64(12.34), ts.VFloat64)
	})

	t.Run("bool_true", func(t *testing.T) {
		err := tools.ParseBool("Bool", "true", &ts.Bool)
		assert.NoError(t, err)
		assert.True(t, ts.Bool)
	})
	t.Run("bool_false", func(t *testing.T) {
		err := tools.ParseBool("Bool", "false", &ts.Bool)
		assert.NoError(t, err)
		assert.False(t, ts.Bool)
	})
}
