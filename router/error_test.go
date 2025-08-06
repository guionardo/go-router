package router

import (
	"errors"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestError(t *testing.T) {
	e := NewError(http.StatusBadRequest, errors.New("test"))
	assert.Equal(t, "status: 400 - Bad Request, error: test", e.Error())

	e = NewErrorF(http.StatusBadGateway, "database error: %d", 10)
	assert.Equal(t, "status: 502 - Bad Gateway, error: database error: 10", e.Error())

	e = NewError(http.StatusInternalServerError, nil)
	assert.Equal(t, "status: 500 - Internal Server Error, error: nil error", e.Error())
}
