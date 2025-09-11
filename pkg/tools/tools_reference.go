package tools

import "reflect"

type toolsRef byte

var ToolsImport = reflect.TypeFor[toolsRef]().PkgPath()
