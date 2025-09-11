package outputs_test

import (
	"os"
	"path"
	"testing"

	"github.com/guionardo/go-router/pkg/outputs"
	"github.com/stretchr/testify/assert"
)

func TestSignFile(t *testing.T) {
	// Create sample file
	fileName := path.Join(t.TempDir(), "sample.go")
	if !assert.NoError(t, os.WriteFile(fileName, []byte(`package testing

	// TEST SIGNATURE`), 0644)) {
		return
	}

	signed, err := outputs.IsFileSigned(fileName)
	assert.False(t, signed)
	assert.NoError(t, err)

	// Sign file
	if !assert.NoError(t, outputs.SignFile(fileName)) {
		return
	}

	signed, err = outputs.IsFileSigned(fileName)
	assert.NoError(t, err)
	assert.True(t, signed)

	// Add content to file
	f, err := os.OpenFile(fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if !assert.NoError(t, err) {
		return
	}
	f.WriteString("\n// APPENDED DATA")
	f.Close()
	signed, err = outputs.IsFileSigned(fileName)
	assert.True(t, signed)
	assert.Error(t, err)

}
