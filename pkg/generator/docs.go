package generator

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"strings"
)

type (
	DocReader struct {
		fileName string
		structs  map[string]structDoc
	}
	structDoc struct {
		name    string
		comment string
		fields  []fieldDoc
	}
	fieldDoc struct {
		name    string
		comment string
	}
)

func NewDocReader(fileName string) (*DocReader, error) {
	dr := &DocReader{
		fileName: fileName,
		structs:  make(map[string]structDoc),
	}
	return dr, dr.Read()
}

func (r *DocReader) Read() error {
	fset := token.NewFileSet()

	f, err := parser.ParseFile(fset, r.fileName, nil, parser.ParseComments)
	if err != nil {
		return err
	}
	var lastComment, lastIdent string
	var lastCommentAsGroup bool
	docText := make([]string, 0, 4)
	ast.Inspect(f, func(n ast.Node) bool {
		if n == nil {
			return true
		}
		switch t := n.(type) {
		case *ast.Comment:
			if !lastCommentAsGroup {
				lastComment = lastComment + t.Text
			}
		case *ast.CommentGroup:
			lastComment = strings.TrimSpace(t.Text())
			lastCommentAsGroup = true
		case *ast.Ident:
			if t.IsExported() {
				lastIdent = t.Name
			}
		case *ast.TypeSpec:
			if dt := t.Doc.Text(); len(dt) > 0 {
				docText = append(docText, dt)
				fmt.Println(dt)
			}

		case *ast.StructType:
			sd := structDoc{
				name:    lastIdent,
				comment: lastComment,
				fields:  make([]fieldDoc, 0, len(t.Fields.List)),
			}

			for _, field := range t.Fields.List {
				sd.fields = append(sd.fields, fieldDoc{name: field.Names[0].Name, comment: strings.TrimSpace(field.Comment.Text())})
			}
			r.structs[lastIdent] = sd
			lastIdent = ""
			lastComment = ""
			lastCommentAsGroup = false
		}
		return true
	})
	return nil
}

func (r *DocReader) Print() {
	fmt.Printf("Filename: %s\n", r.fileName)
	for sn, sr := range r.structs {
		fmt.Printf("\t%s %s\n", sn, sr.comment)
		for _, f := range sr.fields {
			fmt.Printf("\t\t%s %s\n", f.name, f.comment)
		}
	}
}

func ReadDoc(fileName string) error {

	fset := token.NewFileSet()

	f, err := parser.ParseFile(fset, fileName, nil, parser.ParseComments)
	if err != nil {
		panic(err)
	}
	var lastComment, lastIdent string

	docText := make([]string, 0, 4)
	ast.Inspect(f, func(n ast.Node) bool {
		if n == nil {
			return true
		}
		switch t := n.(type) {
		case *ast.TypeSpec:
			if dt := t.Doc.Text(); len(dt) > 0 {
				docText = append(docText, dt)
				fmt.Println(dt)
			}

		case *ast.StructType:
			fmt.Printf("Struct %s %s\n", lastIdent, lastComment)
			for i, field := range t.Fields.List {
				fmt.Printf(" #%d %s //%s\n", i, field.Names[0], field.Comment.Text())
				// fmt.Println(field.Doc.Text())
				// fmt.Println(field.Comment.Text())
			}
		case *ast.GenDecl:

			if t.Doc != nil {
				fmt.Printf("gendecl: %+v\n", t.Doc.Text())
			}

		case *ast.Comment:
			lastComment = t.Text

		case *ast.Ident:
			if t.IsExported() {
				lastIdent = t.Name
			}
		case *ast.ImportSpec, *ast.BasicLit, *ast.File, *ast.CommentGroup, *ast.Field, *ast.FieldList:
			// skip
		default:
			fmt.Printf("%T: %+v\n", t, t)
		}
		return true
	})
	return nil
}
