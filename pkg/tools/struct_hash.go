package tools

import (
	"fmt"
	"hash/crc64"
	"reflect"
	"strings"
)

var Building bool

func TypeHash(t reflect.Type) uint64 {
	sb := strings.Builder{}
	for i := range t.NumField() {
		field := t.Field(i)
		if !field.IsExported() {
			continue
		}
		sb.WriteString(field.Name)
		sb.WriteString(field.Type.Name())
		sb.WriteString(string(field.Tag))
	}

	return crc64.Checksum([]byte(sb.String()), crc64.MakeTable(crc64.ECMA))

}

func ValidateHash[T any](expectedHash uint64) {
	if hash := TypeHash(reflect.TypeFor[T]()); hash != expectedHash {
		if Building {
			fmt.Printf("Regenerating code for type %s", reflect.TypeFor[T]().Name())
		} else {
			panic(fmt.Sprintf("Unexpected modification on struct: Run generate code for type %s", reflect.TypeFor[T]().Name()))
		}
	}
}
