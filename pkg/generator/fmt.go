package generator

import (
	"fmt"
	"io"
	"os"
	"os/exec"
)

type FormatWritter struct {
	data []byte
	w    io.WriteCloser
}

func NewFormatWriter(w io.WriteCloser) io.WriteCloser {
	return &FormatWritter{
		data: make([]byte, 0),
		w:    w,
	}
}

func (fw *FormatWritter) Write(data []byte) (int, error) {
	fw.data = append(fw.data, data...)
	return len(data), nil
}

func (fw *FormatWritter) Close() error {
	out, err := GoFormat(fw.data)
	if err != nil {
		return err
	}
	n, err := fw.w.Write(out)
	if err == nil && n != len(out) {
		err = fmt.Errorf("expected writing %d bytes but %d as written", len(out), n)
	}
	if err == nil {
		err = fw.w.Close()
		if err == nil {
			fw.data = out
		}
	}
	return err
}

func (fw *FormatWritter) String() string {
	return string(fw.data)
}

func GoFormat(data []byte) (out []byte, err error) {
	f, err := os.CreateTemp("", "go_fmt_*.go")
	if err != nil {
		return nil, err
	}
	defer func() {
		f.Close()
		os.Remove(f.Name())
	}()
	n, err := f.Write(data)
	if err == nil && n != len(data) {
		err = fmt.Errorf("error writing temporary file. Expected %d, but %d was written", len(data), n)
	}
	if err != nil {
		return nil, err
	}
	f.Close()
	cmd := exec.Command("go", "fmt", f.Name())
	output, err := cmd.CombinedOutput()
	if err == nil {
		out, err = os.ReadFile(f.Name())
	} else {
		fmt.Printf("error %s", string(output))
	}
	return out, err
}
