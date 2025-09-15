package reflections

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path"
	"reflect"
	"strings"

	"github.com/guionardo/go-router/pkg/tools"
)

type Type[T any] struct {
	PackageName       string
	Pointer           bool
	ProjectRootFolder string
	ModuleName        string
	ModuleFolder      string
	SourceFile        string
	Type              reflect.Type
	IsStruct          bool
	Error             error
}

func New[T any]() *Type[T] {
	tp := reflect.TypeFor[T]()
	return NewFromType[T](tp)
}

func NewFromType[T any](tp reflect.Type) *Type[T] {
	pointer := false
	if tp.Kind() == reflect.Pointer {
		pointer = true
		tp = tp.Elem()
	}
	t := &Type[T]{
		PackageName: path.Base(tp.PkgPath()),
		Pointer:     pointer,
		Type:        tp,
		IsStruct:    tp.Kind() == reflect.Struct,
	}

	root, err := tools.GetProjectRootFolder()
	if err == nil {
		t.ProjectRootFolder = root
		t.ModuleName, err = tools.GetModuleName(path.Join(root, "go.mod"))
		if err == nil {
			sourceFolder, found := strings.CutPrefix(tp.PkgPath(), t.ModuleName)
			if !found {
				err = fmt.Errorf("source folder not found")
			} else {
				t.ModuleFolder = path.Join(root, sourceFolder)
				t.SourceFile, err = findSourceFileForTypeInFolder(t.ModuleFolder, tp.Name())
				if err == nil && len(t.SourceFile) == 0 {
					err = fmt.Errorf("source file not found")
				}
			}
		}
	}
	t.Error = err
	return t
}

func (t Type[T]) FindContentOnFiles(source string) (fileName string) {
	if t.Error != nil {
		return ""
	}

	files, err := os.ReadDir(t.ModuleFolder)
	if err != nil {
		return ""
	}
	sourceContent := []byte(source)
	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".go") {
			content, err := os.ReadFile(path.Join(t.ModuleFolder, file.Name()))
			if err == nil && bytes.Contains(content, sourceContent) {
				return path.Join(t.ModuleFolder, file.Name())
			}
		}
	}
	return ""
}

func findSourceFileForTypeInFolder(sourceFolder, structName string) (sourceFile string, err error) {
	fi, err := os.ReadDir(sourceFolder)
	if err != nil {
		return
	}
	for _, f := range fi {
		if !f.IsDir() {
			sourceFile = path.Join(sourceFolder, f.Name())
			if err = findSourceInFileForType(sourceFile, structName); err == nil {
				return sourceFile, nil
			}
		}
	}
	return "", fmt.Errorf("folder %s doesn not contains file with struct %s", sourceFolder, structName)
}

func findSourceInFileForType(fileName string, structName string) error {
	fset := token.NewFileSet()

	f, err := parser.ParseFile(fset, fileName, nil, 0)
	if err != nil {
		return err
	}

	var (
		lastIdent string
		found     bool
	)
	ast.Inspect(f, func(n ast.Node) bool {
		if n == nil {
			return true
		}
		switch t := n.(type) {
		case *ast.Ident:
			if t.IsExported() {
				lastIdent = t.Name
			}

		case *ast.StructType:
			if lastIdent == structName {
				found = true
				return false
			}

			lastIdent = ""
		}
		return true
	})
	if !found {
		return fmt.Errorf("file %s doesn't contains struct %s", fileName, structName)
	}
	return nil

}
