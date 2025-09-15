package outputs

import (
	"fmt"
	"net/http"
	"path"
	"strings"

	reflections "github.com/guionardo/go-router/pkg/reflect"
	"github.com/guionardo/go-router/pkg/tools"
)

type (
	Outputs[T, R any] struct {
		origin string // file where the struct is defined
		err    error

		PackageName      string
		ParseRequestFile string // file that contains the generated method ParseRequest(r *http.Request) error
		ProcessFile      string // file that contains the method Process(r *http.Request, payload T) error
	}
	Processor[T, R any] interface {
		Process(r *http.Request) (T, error)
	}
)

func New[T, R any]() *Outputs[T, R] {
	o := &Outputs[T, R]{}
	tType := reflections.New[T]()
	if tType.Error != nil {
		o.err = tType.Error
		return o
	}

	if !tType.IsStruct {
		o.err = fmt.Errorf("%s.%s should be a struct", tType.Type.PkgPath(), tType.Type.Name())
		return o
	}
	o.origin = tType.SourceFile

	cuttedOrigin, _ := strings.CutSuffix(o.origin, ".go")

	structNameToFile := tools.ToSnakeCase(tType.Type.Name())
	o.ProcessFile = fmt.Sprintf("%s_%s_process.go", cuttedOrigin, structNameToFile)
	o.ParseRequestFile = fmt.Sprintf("%s_%s_parser.go", cuttedOrigin, structNameToFile)
	o.PackageName = path.Base(tType.PackageName)
	return o
}
