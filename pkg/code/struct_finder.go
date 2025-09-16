package code

import (
	"errors"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io/fs"
	"iter"
	"log/slog"
	"path"
	"path/filepath"
	"slices"
	"strings"

	"github.com/guionardo/go-router/pkg/tools"
)

type (
	GoRouterStruct struct {
		SourceFile  string
		Import      string
		PackageName string
		StructName  string
		RouterData  map[string]string
		AnsweredBy  *GoRouterStruct
	}
	GoRouterStructs []*GoRouterStruct
)

// FindStructsFromBaseCode returns a map with [filename][](Structs)
func FindStructsFromBaseCode(rootFolder string) (GoRouterStructs, error) {
	moduleName, err := tools.GetModuleName(path.Join(rootFolder, "go.mod"))
	if err != nil {
		return nil, err
	}
	var allStructs GoRouterStructs
	errWalk := filepath.WalkDir(rootFolder, func(p string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() || !strings.HasSuffix(d.Name(), ".go") {
			return nil
		}
		structs, err := findStructsFromFile(p, rootFolder, moduleName)
		if err != nil {
			return err
		}
		allStructs = append(allStructs, structs...)
		return nil
	})

	if errWalk != nil {
		slog.Error("failed finding structs", slog.String("basecode", rootFolder), slog.Any("error", errWalk))
		return nil, errWalk
	}
	return allStructs, nil

}

func findStructsFromFile(fileName, rootFolder, moduleName string) (structs []*GoRouterStruct, err error) {
	const logHead = "findingStructs "

	logger := slog.With(slog.String("fileName", fileName))

	fset := token.NewFileSet()

	f, err := parser.ParseFile(fset, fileName, nil, parser.ParseComments)

	if err != nil {
		logger.Error(logHead, slog.Any("error", err))
		return nil, err
	}
	var (
		relativeFolder  = path.Dir(strings.TrimPrefix(fileName, rootFolder))
		importName      = path.Join(moduleName, relativeFolder)
		lastComment     string
		lastIdent       string
		lastCommentNode ast.Node
		packageName     string
		goRouterValues  map[string]string
	)

	ast.Inspect(f, func(n ast.Node) bool {
		if n == nil {
			return true
		}
		switch t := n.(type) {
		case *ast.Comment:
			if !strings.HasPrefix(t.Text, `//gorouter:`) {
				return true
			}
			if grv, err := ExtractKeyValues(strings.TrimPrefix(t.Text, `//`)); err == nil {
				lastComment = t.Text
				lastCommentNode = n
				goRouterValues = grv
			} else {
				logger.Error("invalid gorouter comment", slog.String("comment", t.Text), slog.Any("error", err))
			}

		case *ast.Ident:
			if t.IsExported() {
				lastIdent = t.Name
			}
		// case *ast.TypeSpec: // Type spec é o comentário imediatamente antes de um tipo
		// 	if dt := t.Doc.Text(); len(dt) > 0 {
		// 		docText = append(docText, dt)
		// 		fmt.Println(dt)
		// 	}

		case *ast.StructType:
			if lastCommentNode != nil {
				diff := lastCommentNode.End() - n.Pos()
				if diff < 10 && strings.Contains(lastComment, "gorouter:") {
					structs = append(structs, &GoRouterStruct{
						SourceFile:  fileName,
						Import:      importName,
						PackageName: packageName,
						StructName:  lastIdent,
						RouterData:  goRouterValues,
					})
				}
			}

			lastIdent = ""
			lastComment = ""
			lastCommentNode = nil
			goRouterValues = nil
		case *ast.File:
			packageName = t.Name.Name
			logger = logger.With(slog.String("package", packageName))
		}

		return true
	})
	return

}

func ExtractKeyValues(input string) (map[string]string, error) {
	// Check if the string starts with the required prefix
	const prefix = "gorouter:"
	if !strings.HasPrefix(input, prefix) {
		return nil, fmt.Errorf("string does not start with required prefix: %s", prefix)
	}

	// Remove the prefix
	data := strings.TrimPrefix(input, prefix)

	// Split by spaces to get key=value pairs
	pairs := strings.Fields(data)

	result := make(map[string]string)

	for _, pair := range pairs {
		// Split each pair by '=' to separate key and value
		parts := strings.SplitN(pair, "=", 2)
		if len(parts) != 2 {
			return nil, fmt.Errorf("invalid key=value pair: %s", pair)
		}

		key, value := parts[0], parts[1]

		// Remove quotes if present
		if len(value) >= 2 && value[0] == '"' && value[len(value)-1] == '"' {
			value = value[1 : len(value)-1]
		}

		result[key] = value
	}

	return result, nil
}

const (
	keyPath           = "path"
	keyResponseStruct = "responseTo"
)

func (s *GoRouterStruct) String() string {
	var answeredBy string
	if s.AnsweredBy != nil {
		answeredBy = fmt.Sprintf(" - answered by %s.%s", s.AnsweredBy.PackageName, s.AnsweredBy.StructName)
	}
	return fmt.Sprintf("%s %s.%s%s", s.SourceFile, s.PackageName, s.StructName, answeredBy)
}
func (s GoRouterStructs) IsValid() error {

	if len(s) == 0 {
		return fmt.Errorf("no structs")
	}
	var err error
	for _, st := range s {
		var (
			respStruct = st.RouterData[keyResponseStruct]
		)
		if len(respStruct) == 0 {
			if rs := s.getResponseStructByRequest(st.PackageName, st.StructName); rs == nil {
				err = errors.Join(err, fmt.Errorf("struct %s.%s has no response struct with '%s=%s'", st.PackageName, st.StructName, keyResponseStruct, st.StructName))
			} else {
				st.AnsweredBy = rs
			}
			if _, ok := st.RouterData[keyPath]; !ok {
				err = errors.Join(err, fmt.Errorf("struct %s.%s must have a '%s' value", st.PackageName, st.StructName, keyPath))
			}
		} else if s.getStructByName(st.PackageName, st.StructName) == nil {
			err = errors.Join(err, fmt.Errorf("struct %s.%s responds to an unexistent struct %s.%s", st.PackageName, st.StructName, st.PackageName, respStruct))
		}

	}
	return err

}

func (s GoRouterStructs) GetRequestResponsePairStructs() iter.Seq2[*GoRouterStruct, *GoRouterStruct] {
	return func(yield func(req, resp *GoRouterStruct) bool) {
		for _, st := range s {
			if st.AnsweredBy != nil {
				if !yield(st, st.AnsweredBy) {
					return
				}
			}
		}
	}
}

func (s GoRouterStructs) getResponseStructByRequest(packageName, structName string) *GoRouterStruct {
	i := slices.IndexFunc(s, func(st *GoRouterStruct) bool {
		return st.PackageName == packageName && slices.Contains(strings.Split(st.RouterData[keyResponseStruct], ","), structName)
	})
	if i < 0 {
		return nil
	}
	return s[i]
}

func (s GoRouterStructs) getStructByName(packageName, structName string) *GoRouterStruct {
	structNames := strings.Split(structName, ",")
	i := slices.IndexFunc(s, func(st *GoRouterStruct) bool {
		return st.PackageName == packageName && slices.Contains(structNames, st.StructName)
	})
	if i < 0 {
		return nil
	}
	return s[i]
}
