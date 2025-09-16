package generator

import "reflect"

type generatorRef byte

var GeneratorImport = reflect.TypeFor[generatorRef]().PkgPath()
