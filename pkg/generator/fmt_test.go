package generator

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

const source = `package main
import "fmt"
func check(n int) bool {
return n%2==0
}

func main(){
fmt.Printf("check(2)=%t\n",check(2))
}`

func TestGoFormat(t *testing.T) {

	out, err := GoFormat([]byte(source))
	assert.NoError(t, err)

	outStr := string(out)
	assert.Contains(t, outStr, "\treturn n%2 == 0", "Expected that the row starts with tab")
	assert.Len(t, out, 128)
}

func TestNewFormatWriter(t *testing.T) {
	tmp := t.TempDir()
	tmpOutFile, err := os.CreateTemp(tmp, "go_fmt_*.go")
	assert.NoError(t, err)

	fw := NewFormatWriter(tmpOutFile, tmpOutFile.Name())
	_, err = fw.Write([]byte(source))
	assert.NoError(t, err)

	assert.NoError(t, fw.Close())

	content, err := os.ReadFile(tmpOutFile.Name())
	assert.NoError(t, err)
	assert.Contains(t, string(content), "\treturn n%2 == 0")
}
