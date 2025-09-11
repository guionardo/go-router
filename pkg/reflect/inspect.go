package reflections

import (
	"fmt"
	"io"
	"reflect"
	"strconv"
	"time"
)

func IsStruct[T any]() bool {
	tp := reflect.TypeFor[T]()
	kind := tp.Kind()
	if kind == reflect.Pointer {
		kind = tp.Elem().Kind()
	}
	return kind == reflect.Struct
}

func IsArray[T any]() bool {
	kind := reflect.TypeFor[T]().Kind()
	return kind == reflect.Array || kind == reflect.Slice
}

func IsArrayOfStruct[T any]() bool {
	if !IsArray[T]() {
		return false
	}
	return reflect.TypeFor[T]().Elem().Kind() == reflect.Struct
}

func IsArrayOfByte[T any]() bool {
	if !IsArray[T]() {
		return false
	}
	return reflect.TypeFor[T]().Elem().Kind() == reflect.Uint8
}

func IsString[T any]() bool {
	return reflect.TypeFor[T]().Kind() == reflect.String
}

func IsReader[T any]() bool {
	return reflect.TypeFor[T]().Implements(reflect.TypeFor[io.Reader]())
}

func NewValue(value any, field reflect.Value) (reflect.Value, error) {
	kind := field.Kind()
	switch v := value.(type) {
	case string:
		switch kind {
		case reflect.String:
			return reflect.ValueOf(v), nil
		case reflect.Bool:
			bv, err := strconv.ParseBool(v)
			return reflect.ValueOf(bv), err
		case reflect.Int:
			iv, err := strconv.Atoi(v)
			return reflect.ValueOf(iv), err
		case reflect.Float64:
			fv, err := strconv.ParseFloat(v, 64)
			return reflect.ValueOf(fv), err
		case reflect.Float32:
			fv, err := strconv.ParseFloat(v, 32)
			return reflect.ValueOf(float32(fv)), err
		case reflect.Struct:
			if _, ok := field.Interface().(time.Time); ok {
				tv, err := StrToTime(v)
				return reflect.ValueOf(tv), err
			}
		}

	case bool:
		switch kind {
		case reflect.Bool:
			return reflect.ValueOf(v), nil
		case reflect.String:
			return reflect.ValueOf(BoolValue(v, "true", "false")), nil
		case reflect.Int:
			return reflect.ValueOf(BoolValue(v, 1, 0)), nil
		}
	case []byte:
		switch kind {
		case reflect.Array:
			return reflect.ValueOf(v), nil
		}
	}

	return reflect.ValueOf(nil), fmt.Errorf("unexpected convertion %s to %s", reflect.TypeOf(value).Name(), kind.String())
}
