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
	"runtime"
	"strings"

	pathtools "github.com/guionardo/go/pkg/path_tools"
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

var fsRoot string

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

	root, err := getProjectRootFolder()
	if err == nil {
		t.ProjectRootFolder = root
		t.ModuleName, err = getModuleName(path.Join(root, "go.mod"))
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
func getProjectRootFolder() (root string, err error) {
	cwd, err := os.Getwd()
	if err != nil {
		return "", err
	}
	current := cwd
	for len(current) > len(fsRoot) {
		if pathtools.FileExists(path.Join(current, "go.mod")) {
			return current, nil
		}
		current = path.Dir(current)
	}
	return "", fmt.Errorf("could not find go.mod from current path: %s", cwd)
}

func getModuleName(goModFile string) (string, error) {
	content, err := os.ReadFile(goModFile)
	if err != nil {
		return "", err
	}
	for row := range strings.SplitSeq(string(content), "\n") {
		if module, found := strings.CutPrefix(row, "module "); found {
			return strings.TrimSpace(module), nil
		}
	}
	return "", fmt.Errorf("couldn't find module param in %s", goModFile)
}

func init() {
	if runtime.GOOS == "windows" {
		fsRoot = "C:\\"
	} else {
		fsRoot = "/"
	}
}
