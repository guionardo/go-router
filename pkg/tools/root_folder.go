package tools

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strings"

	pathtools "github.com/guionardo/go/pkg/path_tools"
)

var fsRoot string

func GetProjectRootFolder(cwds ...string) (root string, err error) {
	var cwd string
	if len(cwds) == 0 {
		cwd, err = os.Getwd()
	} else {
		cwd = cwds[0]
	}
	if err == nil {
		cwd, err = filepath.Abs(cwd)
	}
	if err != nil {
		return "", err
	}
	current := cwd
	for len(current) > len(fsRoot) {
		if pathtools.FileExists(path.Join(current, "go.mod")) {
			return current, nil
		}
		current = path.Dir(current)
	}
	return "", fmt.Errorf("could not find go.mod from current path: %s", cwd)
}

func GetModuleName(goModFileArgs ...string) (string, error) {
	var goModFile string
	if len(goModFileArgs) > 0 {
		goModFile = goModFileArgs[0]
	} else {
		root, err := GetProjectRootFolder()
		if err != nil {
			return "", err
		}
		goModFile = path.Join(root, "go.mod")
	}
	content, err := os.ReadFile(goModFile)
	if err != nil {
		return "", err
	}
	for row := range strings.SplitSeq(string(content), "\n") {
		if module, found := strings.CutPrefix(row, "module "); found {
			return strings.TrimSpace(module), nil
		}
	}
	return "", fmt.Errorf("couldn't find module param in %s", goModFile)
}

func init() {
	if runtime.GOOS == "windows" {
		fsRoot = "C:\\"
	} else {
		fsRoot = "/"
	}
}
