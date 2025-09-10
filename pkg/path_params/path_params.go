package path_params

import (
	"fmt"
	"net/url"
	"regexp"
	"strings"
)

func GetPathParamsNames(pattern string) ([]string, error) {
	_, paramsNames, err := createRegexFromPattern(pattern)
	return paramsNames, err
}

func InvalidatePathToHttp(pattern string) (string, error) {
	_, paramsNames, err := createRegexFromPattern(pattern)
	if err != nil {
		return "", err
	}
	ip := pattern
	for _, pn := range paramsNames {
		ip = strings.ReplaceAll(ip, fmt.Sprintf(":%s", pn), fmt.Sprintf("{%s}", pn))
	}
	return ip, nil
}
func createRegexFromPattern(pattern string) (*regexp.Regexp, []string, error) {
	url, err := url.Parse(pattern)
	if err != nil {
		return nil, nil, err
	}
	if strings.Contains(url.Path, ":") {
		return createRegexFromPatternGin(pattern)
	}

	if strings.Contains(url.Path, "{") {
		return createRegexFromPatternHttp(pattern)
	}

	return nil, []string{}, nil
}

// createRegexFromPatternGin converts a URL pattern with :param placeholders (Gin framework style) into a regex
// and returns the compiled regex along with the parameter names in order
// Pattern example: "/users/:id/posts/:postId"
// This will match URLs like: "/users/123/posts/456"
func createRegexFromPatternGin(pattern string) (*regexp.Regexp, []string, error) {
	var paramNames []string

	// Escape special regex characters except for our placeholders
	escaped := regexp.QuoteMeta(pattern)

	// Find all parameter placeholders in the original pattern
	// Gin uses :param syntax, so we look for :followed by word characters
	paramRegex := regexp.MustCompile(`:([a-zA-Z_][a-zA-Z0-9_]*)`)
	matches := paramRegex.FindAllStringSubmatch(pattern, -1)

	for _, match := range matches {
		paramNames = append(paramNames, match[1])
	}

	// Replace escaped placeholders with capture groups
	// We need to replace the escaped version of :param with ([^/]+)
	for _, name := range paramNames {
		escapedPlaceholder := regexp.QuoteMeta(":" + name)
		escaped = strings.Replace(escaped, escapedPlaceholder, "([^/]+)", 1)
	}

	// Compile the final regex
	regex, err := regexp.Compile("^" + escaped + "$")
	if err != nil {
		return nil, nil, err
	}

	return regex, paramNames, nil
}

// CreateRegexFromPatternHttp converts a URL pattern with {param} placeholders into a regex
// and returns the compiled regex along with the parameter names in order
func createRegexFromPatternHttp(pattern string) (*regexp.Regexp, []string, error) {
	var paramNames []string

	// Escape special regex characters except for our placeholders
	escaped := regexp.QuoteMeta(pattern)

	// Find all parameter placeholders in the original pattern
	paramRegex := regexp.MustCompile(`\{([^}]+)\}`)
	matches := paramRegex.FindAllStringSubmatch(pattern, -1)

	for _, match := range matches {
		paramNames = append(paramNames, match[1])
	}

	// Replace escaped placeholders with capture groups
	// We need to replace the escaped version of {param} with ([^/]+)
	for _, name := range paramNames {
		escapedPlaceholder := regexp.QuoteMeta("{" + name + "}")
		escaped = strings.Replace(escaped, escapedPlaceholder, "([^/]+)", 1)
	}

	// Compile the final regex
	regex, err := regexp.Compile("^" + escaped + "$")
	if err != nil {
		return nil, nil, err
	}

	return regex, paramNames, nil
}
