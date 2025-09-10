package main

import (
	"fmt"

	"github.com/guionardo/go-router/cmd/endpoint_gen/structs"
	"github.com/guionardo/go-router/pkg/generator"
)

func main() {
	g := generator.New[*structs.RequestStruct, *structs.ResponseStruct]()
	w, err := g.GetWriter()
	if err != nil {
		panic(err)
	}
	fw := generator.NewFormatWriter(w)
	if err = g.Generate(fw); err != nil {
		panic(err)
	}
	fw.Close()

	fmt.Println(fw.(*generator.FormatWritter).String())
}
