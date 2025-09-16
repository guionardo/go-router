package tools_test

import (
	"testing"

	"github.com/guionardo/go-router/pkg/tools"
	pathtools "github.com/guionardo/go/pkg/path_tools"
	"github.com/stretchr/testify/assert"
)

func TestGetProjectRootFolder(t *testing.T) {
	t.Run("current_project_should_return_valid_folder", func(t *testing.T) {
		got, err := tools.GetProjectRootFolder()
		assert.NoError(t, err)
		assert.True(t, pathtools.DirExists(got))
	})
	t.Run("invalid_project_should_return_error", func(t *testing.T) {
		tmp := t.TempDir()
		got, err := tools.GetProjectRootFolder(tmp)
		assert.Error(t, err)
		assert.Empty(t, got)
	})
}

func TestGetModuleName(t *testing.T) {
	t.Run("get_current_module_should_return_gorouter", func(t *testing.T) {
		got, err := tools.GetModuleName()
		assert.NoError(t, err)
		assert.Contains(t, tools.ToolsImport, got)
	})
}
