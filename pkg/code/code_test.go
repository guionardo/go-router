package code_test

import (
	"os"
	"path"
	"testing"

	"github.com/guionardo/go-router/pkg/code"
	pathtools "github.com/guionardo/go/pkg/path_tools"
	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	tmp := t.TempDir()
	codeFileTest := path.Join(tmp, "code_file_test.go")
	os.Remove(codeFileTest)
	c := code.New("code_test", codeFileTest)
	assert.NotNil(t, c)

	c.AddHeader("test for generated files", "just for fun")
	c.AddImports("fmt")
	c.AddFunc("", "main()", "fmt.Println(\"Hello World!\")", "A simple hello world application")
	c.AddType("funkyStruct", "struct {\nname string \nage int\n}")
	c.AddFunc("f *funkyStruct", "DoSomething()", "fmt.Println(\"Doing something\")", "let's do it")
	err := c.Write()
	assert.NoError(t, err)
	assert.True(t, pathtools.FileExists(codeFileTest))
}
