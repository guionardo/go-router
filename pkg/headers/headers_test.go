package headers

import (
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHeaders(t *testing.T) {
	t.Run("given_not_struct_should_return_error", func(t *testing.T) {
		_, err := New[string]()
		assert.Error(t, err)
	})
	t.Run("given_struct_with_no_string_fields_should_return_error", func(t *testing.T) {
		type sample struct {
			Name int `header:"name"`
		}
		_, err := New[sample]()
		assert.Error(t, err)
	})
	t.Run("given_struct_with_string_fields_should_return_header", func(t *testing.T) {
		type sample struct {
			Name string `header:"name,required"`
			Auth string `header:"auth"`
		}
		header, err := New[sample]()
		assert.NoError(t, err)
		assert.NotNil(t, header)
		assert.Len(t, header.headers, 2)

		s := sample{}
		req := httptest.NewRequest("GET", "/ping", nil)
		req.Header.Add("name", "test")

		err = header.Populate(req, &s)
		assert.NoError(t, err)
		assert.Equal(t, "test", s.Name)
	})
	t.Run("given_struct_with_required_string_fields_should_return_error_if_missing_header", func(t *testing.T) {
		type sample struct {
			Name string `header:"name,required"`
		}
		header, err := New[sample]()
		assert.NoError(t, err)

		s := sample{}
		req := httptest.NewRequest("GET", "/ping", nil)
		err = header.Populate(req, &s)
		assert.Error(t, err)

	})
}

func BenchmarkPopulate(b *testing.B) {
	type sample struct {
		Name string `header:"name,required"`
		Auth string `header:"auth"`
	}
	header, _ := New[sample]()
	req := httptest.NewRequest("GET", "/ping", nil)
	s := sample{}
	for b.Loop() {
		_ = header.Populate(req, &s)
	}
}
