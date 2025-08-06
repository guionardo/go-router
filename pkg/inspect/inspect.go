package inspect

import (
	"io"
	"reflect"
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
