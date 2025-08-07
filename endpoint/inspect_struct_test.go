package endpoint

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type (
	TestStruct struct {
		APIKey string    `header:"apikey" validate:"required"`
		ID     int       `path:"id"`
		Option int       `query:"option"`
		Date   time.Time `query:"date"`
	}
	RespStruct struct {
		Message string `json:"message"`
	}
)

func (ts *TestStruct) Validate() error {
	if ts.ID == 0 {
		return fmt.Errorf("ID should be positive")
	}
	return nil
}

func TestNew(t *testing.T) {
	t.Run("invalid_struct_dont_implements_expected_interfaces_should_fail", func(t *testing.T) {
		type InvalidStruct struct {
			None int
		}
		assert.Panics(t, func() {
			_ = New[InvalidStruct, RespStruct]("/ping")
		})
	})

	t.Run("get_request_should_succeed", func(t *testing.T) {
		inspectStruct := New[TestStruct, RespStruct]("/test/:id")
		if !assert.NotNil(t, inspectStruct) {
			return
		}
		req := httptest.NewRequest(http.MethodGet, "/test/1234?option=9874343&date=2025-07-30", nil)
		req.Header.Add("apikey", "API KEY")
		s, err := inspectStruct.Parse(req)
		assert.NoError(t, err)
		assert.Equal(t, 1234, s.ID)
		assert.Equal(t, "API KEY", s.APIKey)
		assert.Equal(t, 9874343, s.Option)
		assert.Equal(t, time.Date(2025, 7, 30, 0, 0, 0, 0, time.UTC), s.Date)
	})
	t.Run("get_request_without_required_field_should_fail", func(t *testing.T) {
		inspectStruct := New[TestStruct, RespStruct]("/test/:id")
		if !assert.NotNil(t, inspectStruct) {
			return
		}
		req := httptest.NewRequest(http.MethodGet, "/test/1234?option=9874343&date=2025-07-30", nil)
		_, err := inspectStruct.Parse(req)
		assert.Error(t, err)
	})

	t.Run("post_request_with_bytes_body", func(t *testing.T) {
		type TestBodyStruct struct {
			ID   int    `path:"id"`
			Body []byte `body:"body"`
		}
		inspectStruct := New[TestBodyStruct, RespStruct]("/test/:id")
		if !assert.NotNil(t, inspectStruct) {
			return
		}
		req := httptest.NewRequest(http.MethodGet, "/test/22", bytes.NewBuffer([]byte{1, 2, 3, 4, 5, 6}))
		is, err := inspectStruct.Parse(req)
		assert.NoError(t, err)
		assert.Equal(t, 22, is.ID)
		assert.Equal(t, []byte{1, 2, 3, 4, 5, 6}, is.Body)

	})
	t.Run("post_request_with_string_body", func(t *testing.T) {
		type TestBodyStruct struct {
			ID   int    `path:"id"`
			Body string `body:"body"`
		}
		inspectStruct := New[TestBodyStruct, RespStruct]("/test/:id")
		if !assert.NotNil(t, inspectStruct) {
			return
		}
		req := httptest.NewRequest(http.MethodGet, "/test/22", bytes.NewBufferString("SOME CONTENT\nHELLO"))
		is, err := inspectStruct.Parse(req)
		assert.NoError(t, err)
		assert.Equal(t, 22, is.ID)
		assert.Equal(t, "SOME CONTENT\nHELLO", is.Body)

	})

	t.Run("post_request_with_struct_body", func(t *testing.T) {
		type SubBody struct {
			Name string `json:"name"`
			Age  int    `json:"age"`
		}
		type TestBodyStruct struct {
			ID   int     `path:"id"`
			Body SubBody `body:"body"`
		}
		var inspectStruct *Endpoint[TestBodyStruct, RespStruct]
		assert.NotPanics(t, func() {
			inspectStruct = New[TestBodyStruct, RespStruct]("/test/:id")
		})

		if !assert.NotNil(t, inspectStruct) {
			return
		}
		req := httptest.NewRequest(http.MethodGet, "/test/25", bytes.NewBufferString("{\"name\":\"Guionardo\",\"age\":48}"))
		is, err := inspectStruct.Parse(req)
		assert.NoError(t, err)
		assert.Equal(t, 25, is.ID)
		assert.Equal(t, SubBody{"Guionardo", 48}, is.Body)

	})
	t.Run("post_request_with_map_body", func(t *testing.T) {
		type TestBodyStruct struct {
			ID   int            `path:"id"`
			Body map[string]any `body:"body"`
		}
		inspectStruct := New[TestBodyStruct, RespStruct]("/test/:id")
		if !assert.NotNil(t, inspectStruct) {
			return
		}
		req := httptest.NewRequest(http.MethodGet, "/test/25", bytes.NewBufferString("{\"name\":\"Guionardo\",\"age\":48}"))
		is, err := inspectStruct.Parse(req)
		assert.NoError(t, err)
		assert.Equal(t, 25, is.ID)
		assert.Equal(t, map[string]any{"name": "Guionardo", "age": 48}, is.Body)

	})

}
