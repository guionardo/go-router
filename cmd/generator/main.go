package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path"
	"strings"

	"github.com/guionardo/go-router/pkg/code"
	"github.com/guionardo/go-router/pkg/generator"
	"github.com/guionardo/go-router/pkg/tools"
	pathtools "github.com/guionardo/go/pkg/path_tools"
)

var (
	projectRoot string
)

func main() {
	projectRoot, _ = os.Getwd()
	flag.StringVar(&projectRoot, "project_root", projectRoot, "project root folder (containing a valid go.mod file)")

	flag.Parse()
	if len(projectRoot) == 0 {
		flag.Usage()
		return
	}
	projectRoot, err := tools.GetProjectRootFolder(projectRoot)
	if err != nil {
		fmt.Println(err.Error())
		flag.Usage()
		return
	}
	moduleName, err := tools.GetModuleName(path.Join(projectRoot, "go.mod"))
	if err != nil {
		fmt.Printf("Invalid go.mod: %s\n", err.Error())
		return
	}
	fmt.Println("Generating go-router methods for structs")
	fmt.Printf("Project root: %s\n", projectRoot)
	fmt.Printf("Module name : %s\n", moduleName)
	structs, err := code.FindStructsFromBaseCode(projectRoot)
	if err != nil {
		fmt.Printf("Error reading structs: %s\n", err.Error())
		return
	}

	if err = structs.IsValid(); err != nil {
		fmt.Printf("%s\n", err.Error())
		return
	}

	for _, s := range structs {
		fmt.Println(s)
	}

	outputFolder := path.Join(projectRoot, "go_router_temp")

	if err := pathtools.CreatePath(outputFolder); err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	defer os.RemoveAll(outputFolder)

	c := code.New("main", path.Join(outputFolder, "main.go"))
	c.AddImports(generator.GeneratorImport)
	generators := []string{}
	for req, resp := range structs.GetRequestResponsePairStructs() {
		c.AddImports(req.Import)
		funcGen := fmt.Sprintf("gen_%s_%s_%s() error", req.PackageName, req.StructName, resp.StructName)
		body := fmt.Sprintf(`return generator.New[*%s.%s,*%s.%s]().Generate()`, req.PackageName, req.StructName, resp.PackageName, resp.StructName)

		c.AddFunc("", funcGen, body, "")
		generators = append(generators, fmt.Sprintf("gen_%s_%s_%s", req.PackageName, req.StructName, resp.StructName))
	}
	mainBody := fmt.Sprintf(`generators:=[]func()error{%s}
	for _,g:=range generators{
	if err:=g();err!=nil{
	panic(err)
	}
	}`, strings.Join(generators, ","))
	c.AddFunc("", "main()", mainBody, "")
	if err := c.Write(); err != nil {
		panic(err)
	}

	cmd := exec.Command("go", "run", c.FileName)
	output, err := cmd.CombinedOutput()
	fmt.Println(string(output))
	if err != nil {
		panic(err)
	}

}
