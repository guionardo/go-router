package reflections_test

import (
	"testing"

	reflections "github.com/guionardo/go-router/pkg/reflect"
	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	tt := reflections.New[*reflections.MockStruct]()
	assert.NoError(t, tt.Error)

	fileName := tt.FindContentOnFiles("type MockStruct")
	assert.NotEmpty(t, fileName)
}
