package generator

import (
	"context"
	"fmt"
	"io"
)

type GoFileBuilder struct {
	ctx    context.Context
	w      io.Writer
	cancel context.CancelFunc
	err    error
}

func NewGoFileBuilder(w io.Writer) *GoFileBuilder {
	ctx, cancel := context.WithCancel(context.Background())
	return &GoFileBuilder{
		ctx:    ctx,
		w:      w,
		cancel: cancel,
	}
}

func (g *GoFileBuilder) write(writeFuncs ...func() error) {
	for _, wf := range writeFuncs {
		select {
		case <-g.ctx.Done():
			return
		default:
			if err := wf(); err != nil {
				g.cancel()
				g.err = err
			}
		}
	}
}

func (g *GoFileBuilder) wf(format string, args ...any) func() error {
	return func() error {
		return g.ws(format, args...)
	}
}

func (g *GoFileBuilder) ws(format string, args ...any) error {
	row := []byte(fmt.Sprintf(format, args...) + "\n")
	n, err := g.w.Write(row)
	if err == nil && len(row) != n {
		err = fmt.Errorf("expected writing %d bytes, but %d was written. Source: %s", len(row), n, string(row))
	}
	return err
}
