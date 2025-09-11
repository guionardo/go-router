package tools_test

import (
	"testing"
	"time"

	"github.com/guionardo/go-router/pkg/tools"
	"github.com/stretchr/testify/assert"
)

func TestParseTime(t *testing.T) {
	t.Run("valid_datetime", func(t *testing.T) {
		d, err := tools.ParseTimeLayouts("2025-01-01 10:00:00")
		assert.NoError(t, err)
		assert.Equal(t, time.Date(2025, 1, 1, 10, 0, 0, 0, time.UTC), d)
	})
	t.Run("invalid_datetime", func(t *testing.T) {
		_, err := tools.ParseTimeLayouts("2025")
		assert.Error(t, err)
	})
}
