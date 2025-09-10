package tools

import (
	"strconv"
	"time"
)

type (
	Ints interface {
		int | int8 | int16 | int32 | int64
	}
)

func ParseInt(fieldName, s string, v *int) error {
	va, err := strconv.ParseInt(s, 0, strconv.IntSize)
	*v = int(va)
	return NewParseError(fieldName, err)
}
func ParseInt8(fieldName, s string, v *int8) error {
	va, err := strconv.ParseInt(s, 0, 8)
	*v = int8(va)
	return NewParseError(fieldName, err)
}

func ParseInt16(fieldName, s string, v *int16) error {
	va, err := strconv.ParseInt(s, 0, 16)
	*v = int16(va)
	return NewParseError(fieldName, err)
}

func ParseInt32(fieldName, s string, v *int32) error {
	va, err := strconv.ParseInt(s, 0, 32)
	*v = int32(va)
	return NewParseError(fieldName, err)
}
func ParseInt64(fieldName, s string, v *int64) error {
	va, err := strconv.ParseInt(s, 0, 64)
	*v = int64(va)
	return NewParseError(fieldName, err)
}

func ParseUInt(fieldName, s string, v *uint) error {
	va, err := strconv.ParseUint(s, 0, strconv.IntSize)
	*v = uint(va)
	return NewParseError(fieldName, err)
}
func ParseUInt8(fieldName, s string, v *uint8) error {
	va, err := strconv.ParseUint(s, 0, 8)
	*v = uint8(va)
	return NewParseError(fieldName, err)
}

func ParseUInt16(fieldName, s string, v *uint16) error {
	va, err := strconv.ParseUint(s, 0, 16)
	*v = uint16(va)
	return NewParseError(fieldName, err)
}

func ParseUInt32(fieldName, s string, v *uint32) error {
	va, err := strconv.ParseUint(s, 0, 32)
	*v = uint32(va)
	return NewParseError(fieldName, err)
}
func ParseUInt64(fieldName, s string, v *uint64) error {
	va, err := strconv.ParseUint(s, 0, 64)
	*v = uint64(va)
	return NewParseError(fieldName, err)
}

func ParseFloat32(fieldName, s string, v *float32) error {
	va, err := strconv.ParseFloat(s, 32)
	*v = float32(va)
	return NewParseError(fieldName, err)
}

func ParseFloat64(fieldName, s string, v *float64) error {
	va, err := strconv.ParseFloat(s, 64)
	*v = va
	return NewParseError(fieldName, err)
}

func ParseBool(fieldName, s string, v *bool) error {
	va, err := strconv.ParseBool(s)
	*v = va
	return NewParseError(fieldName, err)
}

func ParseTime(fieldName, s string, v *time.Time) error {
	va, err := ParseTimeLayouts(s)
	*v = va
	return NewParseError(fieldName, err)
}

func ParseDuration(fieldName, s string, v *time.Duration) error {
	va, err := time.ParseDuration(s)
	*v = va
	return NewParseError(fieldName, err)
}
