package attributor

import (
	"fmt"
	"reflect"

	"github.com/guionardo/go-router/pkg/tools"
	"github.com/guionardo/go/pkg/set"
)

type (
	Parser interface {
		Code(receptor, fieldName, value string) string
		Imports() []string
	}
	ConcreteParser struct {
		imports            set.Set[string]
		convFunc           string
		castTempl          string
		invalidTypeMessage string
		attrFunc           string
	}
	StringParser struct {
	}
)

var parsers = make(map[string]Parser)

func NewParser[T any]() Parser {
	t := reflect.TypeFor[T]()
	return NewParserFromType(t)
}

func NewParserFromType(t reflect.Type) Parser {
	if t.Kind() == reflect.String {
		return &StringParser{}
	}
	typeName := fmt.Sprintf("%s.%s", t.PkgPath(), t.Name())
	if p, ok := parsers[typeName]; ok {
		return p
	}

	p := &ConcreteParser{
		castTempl: "%s",
	}
	switch t.Kind() {
	case reflect.Int:
		p.imports.Add(tools.ToolsImport)
		p.attrFunc = `func () error { return tools.ParseInt("%s",%s,&%s.%s)}`

	case reflect.Uint:
		p.imports.Add(tools.ToolsImport)
		p.attrFunc = `func () error { return tools.ParseUInt("%s",%s,&%s.%s)}`

	case reflect.Int8:
		p.imports.Add(tools.ToolsImport)
		p.attrFunc = `func () error { return tools.ParseInt8("%s",%s,&%s.%s)}`

	case reflect.Uint8:
		p.imports.Add(tools.ToolsImport)
		p.attrFunc = `func () error { return tools.ParseUInt8("%s",%s,&%s.%s)}`

	case reflect.Int16:
		p.imports.Add(tools.ToolsImport)
		p.attrFunc = `func () error { return tools.ParseInt16("%s",%s,&%s.%s)}`

	case reflect.Uint16:
		p.imports.Add(tools.ToolsImport)
		p.attrFunc = `func () error { return tools.ParseUInt16("%s",%s,&%s.%s)}`

	case reflect.Int32:
		p.imports.Add(tools.ToolsImport)
		p.attrFunc = `func () error { return tools.ParseInt32("%s",%s,&%s.%s)}`

	case reflect.Uint32:
		p.imports.Add(tools.ToolsImport)
		p.attrFunc = `func () error { return tools.ParseUInt32("%s",%s,&%s.%s)}`

	case reflect.Int64:
		if t.PkgPath() == "time" {
			if t.Name() == "Duration" {
				p.imports.Add(tools.ToolsImport)
				p.attrFunc = `func () error { return tools.ParseDuration("%s",%s,&%s.%s)}`

			}
		} else {
			p.imports.Add(tools.ToolsImport)
			p.attrFunc = `func () error { return tools.ParseInt64("%s",%s,&%s.%s)}`

		}
	case reflect.Uint64:
		p.imports.Add(tools.ToolsImport)
		p.attrFunc = `func () error { return tools.ParseUInt64("%s",%s,&%s.%s)}`

	case reflect.Bool:
		p.imports.Add(tools.ToolsImport)
		p.attrFunc = `func () error { return tools.ParseBool("%s",%s,&%s.%s)}`

	case reflect.Float32:
		p.imports.Add(tools.ToolsImport)
		p.attrFunc = `func () error { return tools.ParseFloat32("%s",%s,&%s.%s)}`

	case reflect.Float64:
		p.imports.Add(tools.ToolsImport)
		p.attrFunc = `func () error { return tools.ParseFloat64("%s",%s,&%s.%s)}`

	case reflect.Struct:
		if t.PkgPath() == "time" {
			if t.Name() == "Time" {
				p.imports.Add(tools.ToolsImport)
				p.attrFunc = `func () error { return tools.ParseTime("%s",%s,&%s.%s)}`

			}
		}

	}
	if len(p.attrFunc) == 0 {
		p.invalidTypeMessage = fmt.Sprintf("// unparseable type %s", typeName)
	}

	parsers[typeName] = p
	return p
}

func (pi *ConcreteParser) Code(receptor, fieldName, value string) string {
	if len(pi.invalidTypeMessage) > 0 {
		return pi.invalidTypeMessage
	}
	convFunc := fmt.Sprintf(pi.convFunc, value)
	valueVar := fmt.Sprintf(pi.castTempl, "value")
	if len(pi.attrFunc) > 0 {
		return fmt.Sprintf(pi.attrFunc, fieldName, value, receptor, fieldName)
	}
	return fmt.Sprintf(`if value,err:=%s; err!=nil { return err } else { %s.%s = %s }`,
		convFunc, receptor, fieldName, valueVar)

}

func (pi *ConcreteParser) Imports() []string {
	imports := make([]string, 0, len(pi.imports))
	for i := range pi.imports.Iter() {
		imports = append(imports, i)
	}
	return imports
}

func (sp *StringParser) Code(receptor, fieldName, value string) string {
	return fmt.Sprintf(`func () error { %s.%s = %s
	return nil }`, receptor, fieldName, value)
}
func (sp *StringParser) Imports() []string {
	return []string{}
}
