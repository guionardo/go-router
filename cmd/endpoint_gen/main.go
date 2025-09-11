package main

import (
	"github.com/guionardo/go-router/cmd/endpoint_gen/structs"
	"github.com/guionardo/go-router/pkg/generator"
)

func main() {
	g := generator.New[*structs.RequestStruct, *structs.ResponseStruct]()

	if err := g.Generate(); err != nil {
		panic(err)
	}

}
