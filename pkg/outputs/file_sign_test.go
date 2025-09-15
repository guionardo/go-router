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
	t.Run("0_create_sample", func(t *testing.T) {
		assert.NoError(t, os.WriteFile(fileName, []byte(`package testing

		// TEST SIGNATURE`), 0644))
	})

	t.Run("1_test_unsigne_file_should_return_false", func(t *testing.T) {
		signed, err := outputs.IsFileSigned(fileName)
		assert.False(t, signed)
		assert.NoError(t, err)
	})

	t.Run("2_sign_file_should_succeed", func(t *testing.T) {
		assert.NoError(t, outputs.SignFile(fileName))
	})

	t.Run("3_test_signed_file_should_return_true", func(t *testing.T) {
		signed, err := outputs.IsFileSigned(fileName)
		assert.NoError(t, err)
		assert.True(t, signed)
	})

	t.Run("4_add_content_to_file_should_return_false_on_sign", func(t *testing.T) {
		f, err := os.OpenFile(fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if !assert.NoError(t, err) {
			return
		}
		_, _ = f.WriteString("\n// APPENDED DATA")
		_ = f.Close()
		signed, err := outputs.IsFileSigned(fileName)
		assert.True(t, signed)
		assert.Error(t, err)

	})

}
