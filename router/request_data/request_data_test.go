package request_data

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	t.Run("get_ping", func(t *testing.T) {
		pl, err := New[string]("/ping", WithIgnoreBody)
		if !assert.NoError(t, err) {
			return
		}
		if !assert.NotNil(t, pl) {
			return
		}

		req := httptest.NewRequest(http.MethodGet, "http://localhost/ping", nil)
		out, err := pl.Handle(req)
		assert.NoError(t, err)
		assert.Empty(t, out)
	})
	t.Run("get_ping_with_path", func(t *testing.T) {
		type input struct {
			Id int `path:"id" json:"id"`
		}
		pl, err := New[input]("/ping/:id", WithIgnoreBody)
		if !assert.NoError(t, err) {
			return
		}
		if !assert.NotNil(t, pl) {
			return
		}
		req := httptest.NewRequest(http.MethodGet, "http://localhost/ping/10", nil)
		out, err := pl.Handle(req)
		assert.NoError(t, err)
		assert.Equal(t, 10, out.Id)

	})
	t.Run("get_ping_with_path_and_query", func(t *testing.T) {
		type input struct {
			Id      int    `path:"id" json:"id"`
			Message string `urlenc:"message"`
		}
		pl, err := New[input]("/ping/:id", WithIgnoreBody)
		if !assert.NoError(t, err) {
			return
		}
		if !assert.NotNil(t, pl) {
			return
		}
		req := httptest.NewRequest(http.MethodGet, "http://localhost/ping/20?message=ABCDEFG", nil)
		out, err := pl.Handle(req)
		assert.NoError(t, err)
		assert.Equal(t, 20, out.Id)
		assert.Equal(t, "ABCDEFG", out.Message)

	})
	t.Run("post_data_with_path_and_query", func(t *testing.T) {
		type input struct {
			Id      int    `path:"id" json:"id"`
			Message string `urlenc:"message"`
			Numbers []int  `json:"numbers"`
		}
		pl, err := New[input]("/ping/:id")
		if !assert.NoError(t, err) {
			return
		}
		if !assert.NotNil(t, pl) {
			return
		}

		payload := input{Numbers: []int{2, 4, 8, 0, 23, 231}}
		content, _ := json.Marshal(payload)

		req := httptest.NewRequest(http.MethodPost, "http://localhost/ping/30?"+url.QueryEscape("message=Some Text"), bytes.NewBuffer(content))
		out, err := pl.Handle(req)
		assert.NoError(t, err)
		assert.Equal(t, 30, out.Id)
		assert.Equal(t, "Some Text", out.Message)
		assert.Equal(t, payload.Numbers, out.Numbers)
	})
	t.Run("get_ping_with_path_validating", func(t *testing.T) {
		type input struct {
			Id int `path:"id" json:"id" validate:"gte=1,lte=10"`
		}
		pl, err := New[input]("/ping/:id", WithIgnoreBody)
		if !assert.NoError(t, err) {
			return
		}
		if !assert.NotNil(t, pl) {
			return
		}
		req := httptest.NewRequest(http.MethodGet, "http://localhost/ping/0", nil)
		_, err = pl.Handle(req)
		assert.Error(t, err)

	})

}
