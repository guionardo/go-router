package router

import "reflect"

type routerRef byte

var RouterImport = reflect.TypeFor[routerRef]().PkgPath()
