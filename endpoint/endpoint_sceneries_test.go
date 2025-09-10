package endpoint_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/guionardo/go-router/endpoint"
	"github.com/guionardo/go-router/pkg/errors"
	"github.com/guionardo/go-router/pkg/sceneries"
	"github.com/guionardo/go-router/router"
	"github.com/stretchr/testify/assert"
)

func getServer(r *router.Router) *httptest.Server {
	mux := http.NewServeMux()
	r.SetupHTTP(mux)
	return httptest.NewServer(mux)
}
func TestGetSimple(t *testing.T) {
	var e *endpoint.Endpoint[sceneries.GetSimple, sceneries.GetSimpleResponse]
	if !assert.NotPanics(t, func() {
		e = endpoint.New[sceneries.GetSimple, sceneries.GetSimpleResponse]("/")
	}) {
		return
	}
	r := router.New().Get(e)

	server := getServer(r)
	defer server.Close()
	client := server.Client()

	t.Run("request_valid_should_succeed", func(t *testing.T) {
		resp, err := client.Get(server.URL + "/?id=1")
		assert.NoError(t, err)

		assert.Equal(t, http.StatusAccepted, resp.StatusCode)

		var response sceneries.GetSimpleResponse
		err = json.NewDecoder(resp.Body).Decode(&response)
		assert.NoError(t, err)
		assert.True(t, response.Success)
	})
	t.Run("request_missing_argument_should_return_bad_request", func(t *testing.T) {
		resp, err := client.Get(server.URL + "/")
		assert.NoError(t, err)

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

		var response errors.ParseErrorStruct
		err = json.NewDecoder(resp.Body).Decode(&response)
		assert.NoError(t, err)
		assert.NotEmpty(t, response.Errors)
	})

}

func TestPostSimple(t *testing.T) {
	t.Setenv("ROUTER_LOGGING", "true")
	t.Setenv("ENVIRONMENT", "LOCAL")
	var e *endpoint.Endpoint[sceneries.PostSimple, sceneries.PostSimpleResponse]
	if !assert.NotPanics(t, func() {
		e = endpoint.New[sceneries.PostSimple, sceneries.PostSimpleResponse]("/simple/:id")
	}) {
		return
	}
	r := router.New().Post(e)

	server := getServer(r)
	defer server.Close()
	client := server.Client()
	t.Run("", func(t *testing.T) {
		resp, err := client.Post(server.URL+"/simple/1", "application/json", bytes.NewBufferString(`{"name":"Guionardo","age":48}`))
		assert.NoError(t, err)
		assert.Equal(t, http.StatusAccepted, resp.StatusCode)

		var response sceneries.PostSimpleResponse
		err = json.NewDecoder(resp.Body).Decode(&response)
		assert.NoError(t, err)
		assert.Equal(t, 48, response.Length)
	})
}
