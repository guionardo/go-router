package tools_test

import (
	"testing"

	"github.com/guionardo/go-router/pkg/tools"
	"github.com/stretchr/testify/assert"
)

func TestLastSuccessIterator_Iter(t *testing.T) {
	iter := tools.NewLastSuccessIterator(1, 2, 3, 4, 5, 6)
	for item := range iter.Iter() {
		if item > 5 {
			break
		}
	}

	for item := range iter.Iter() {
		assert.Equal(t, 6, item)
		break
	}
}
