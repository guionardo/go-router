package path_params

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPathParams_Inject(t *testing.T) {
	type PL struct {
		ID   string `json:"id"`
		Name string `json:"name"`
		Age  int    `json:"age"`
	}
	pp := &PathParams{
		params: map[string]string{
			"id":   "123",
			"name": "john",
		},
		paramsNames: []string{"id", "name"},
		regex:       nil,
	}
	pl := &PL{
		Age: 20,
	}
	err := pp.Inject(&pl)
	assert.NoError(t, err)
	assert.Equal(t, "123", pl.ID)
	assert.Equal(t, "john", pl.Name)
	assert.Equal(t, 20, pl.Age)
}
