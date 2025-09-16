package code_test

import (
	"reflect"
	"testing"

	"github.com/guionardo/go-router/pkg/code"
	"github.com/guionardo/go-router/pkg/tools"
	"github.com/stretchr/testify/assert"
)

func TestFindStructsFromBaseCode(t *testing.T) {
	root, _ := tools.GetProjectRootFolder()
	findings, err := code.FindStructsFromBaseCode(root)
	assert.NoError(t, err)
	assert.NotEmpty(t, findings)
}

func Test_extractKeyValues(t *testing.T) {

	tests := []struct {
		name    string
		input   string
		want    map[string]string
		wantErr bool
	}{
		{"empty_should_return_error", "", nil, true},
		{"with_values_should_return_map", `gorouter:path="/home/test" enabled=true`, map[string]string{"path": "/home/test", "enabled": "true"}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := code.ExtractKeyValues(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("extractKeyValues() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("extractKeyValues() = %v, want %v", got, tt.want)
			}
		})
	}
}
