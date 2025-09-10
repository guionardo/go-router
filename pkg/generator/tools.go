package generator

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"path"
	"regexp"
	"slices"
	"strings"
)

var (
	matchFirstCap = regexp.MustCompile("(.)([A-Z][a-z]+)")
	matchAllCap   = regexp.MustCompile("([a-z0-9])([A-Z])")
)

func toSnakeCase(str string) string {
	snake := matchFirstCap.ReplaceAllString(str, "${1}_${2}")
	snake = matchAllCap.ReplaceAllString(snake, "${1}_${2}")
	return strings.ToLower(snake)
}

func getProjectRootFolder() (string, error) {
	currentPath, err := os.Getwd()
	if err != nil {
		return "", err
	}
	for len(currentPath) > 1 {
		if stat, err := os.Stat(path.Join(currentPath, "go.mod")); err == nil && !stat.IsDir() {
			return currentPath, nil
		}
		currentPath = path.Dir(currentPath)
	}
	return "", errors.New("couldn't find the project root")
}

func getModuleName(goModFile string) (string, error) {
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

func findSourceInFileForType(fileName string, structName string) error {
	f, err := os.Open(fileName)
	if err != nil {
		return err
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		words := strings.Split(strings.TrimSpace(scanner.Text()), " ")
		if pStructName := slices.Index(words, structName); pStructName >= 0 {
			if pStructType := slices.Index(words, "struct"); pStructType > pStructName {
				return nil
			}
		}
	}
	return fmt.Errorf("file %s doesn't contains struct %s", fileName, structName)
}

func findSourceFileForType(packageFolder string, structName string) (string, error) {
	files, err := os.ReadDir(packageFolder)
	if err != nil {
		return "", err
	}
	for _, file := range files {
		if matched, err := path.Match("*.go", file.Name()); !matched || err != nil {
			continue
		}
		sourceFileName := path.Join(packageFolder, file.Name())
		if err := findSourceInFileForType(sourceFileName, structName); err == nil {
			return sourceFileName, nil
		}
	}
	return "", fmt.Errorf("couldn't find struct %s in files of folder %s", structName, packageFolder)
}
