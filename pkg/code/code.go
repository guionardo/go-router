package code

import (
	"fmt"
	"io"
	"os"
	"slices"
	"sort"

	"github.com/guionardo/go/pkg/set"
)

type (
	// Code is an abstraction to generate a source file
	Code struct {
		FileName    string
		PackageName string
		Header      []string
		imports     set.Set[string]
		funcs       []*Func
		types       map[string]string
		raw         []string

		w io.WriteCloser
	}
	Func struct {
		Receiver string
		Name     string
		Body     string
		Comment  string
	}
)

func New(packageName, fileName string) *Code {
	return &Code{
		FileName:    fileName,
		PackageName: packageName,
		Header:      make([]string, 0),
		imports:     set.New[string](),
		funcs:       make([]*Func, 0),
		types:       make(map[string]string),
		raw:         make([]string, 0),
	}
}

func (c *Code) AddHeader(headerRows ...string) {
	c.Header = append(c.Header, headerRows...)
}

func (c *Code) AddImports(imports ...string) {
	c.imports.AddMultiple(imports...)
}

func (c *Code) AddType(typeName, typeBody string) {
	c.types[typeName] = typeBody
}

func (c *Code) AddFunc(receiver, name, body, comment string) {
	index := slices.IndexFunc(c.funcs, func(f *Func) bool {
		return f.Receiver == receiver && f.Name == name
	})
	if index == -1 {
		c.funcs = append(c.funcs, &Func{
			Receiver: receiver,
			Name:     name,
			Body:     body,
			Comment:  comment,
		})
	} else {
		c.funcs[index].Body = body
		c.funcs[index].Comment = comment
	}
}

func (c *Code) AddRaw(code string) {
	c.raw = append(c.raw, code)
}

func (c *Code) Write() error {
	f, err := os.CreateTemp("", "*.go")
	if err != nil {
		return err
	}
	tmpFile := f.Name()

	c.w = NewFormatWriter(f, c.FileName)

	err = c.writeHeader()
	if err == nil {
		err = c.writePackage()
	}
	if err == nil {
		err = c.writeImports()
	}
	if err == nil {
		err = c.writeTypes()
	}
	if err == nil {
		err = c.writeFuncs()
	}

	if err == nil {
		err = c.writeRaw()
	}
	errClose := c.w.Close()
	if err == nil {
		if errClose == nil {
			os.Remove(c.FileName)
			err = os.Rename(tmpFile, c.FileName)
		} else {
			err = errClose
		}
	}

	return err
}

func (c *Code) write(format string, args ...any) error {
	_, err := fmt.Fprintf(c.w, format+"\n", args...)
	return err
}

func (c *Code) writeHeader() error {
	for _, header := range c.Header {
		if err := c.write("// %s", header); err != nil {
			return err
		}
	}

	return nil
}

func (c *Code) writePackage() error {
	return c.write("package %s", c.PackageName)
}

func (c *Code) writeImports() (err error) {
	switch len(c.imports) {
	case 0:
		return nil
	case 1:
		for imp := range c.imports.Iter() {
			return c.write(`import "%s"`, imp)
		}
		return nil
	default:
		sorted := c.imports.ToArray()
		sort.Strings(sorted)
		if err = c.write("import ("); err == nil {
			for _, imp := range sorted {
				if err = c.write("\t\"%s\"", imp); err != nil {
					return err
				}
			}
		}
		return c.write(")")
	}
}

func (c *Code) writeTypes() (err error) {
	switch len(c.types) {
	case 0:
		return nil
	case 1:
		for t, b := range c.types {
			return c.write("type %s %s", t, b)
		}
	default:
		keys := []string{}
		for k := range c.types {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		if err = c.write("type ("); err == nil {
			for _, typeName := range keys {
				if err = c.write("\t%s %s", typeName, c.types[typeName]); err != nil {
					break
				}
			}
		}
	}
	return err
}

func (c *Code) writeFuncs() (err error) {
	sort.Slice(c.funcs, func(i int, j int) bool {
		signI := fmt.Sprintf("%s %s", c.funcs[i].Receiver, c.funcs[i].Name)
		signJ := fmt.Sprintf("%s %s", c.funcs[j].Receiver, c.funcs[j].Name)
		return signI < signJ
	})
	for _, f := range c.funcs {
		if err == nil && len(f.Comment) > 0 {
			err = c.write("// %s %s", f.Name, f.Comment)
		}
		if err == nil && len(f.Receiver) == 0 {
			err = c.write("func %s {\n%s\n}\n", f.Name, f.Body)
		} else if err == nil {
			err = c.write("func (%s) %s {\n%s\n}\n", f.Receiver, f.Name, f.Body)
		}
	}
	return nil
}

func (c *Code) writeRaw() error {
	for _, raw := range c.raw {
		if err := c.write(raw); err != nil {
			return err
		}
	}
	return nil
}
