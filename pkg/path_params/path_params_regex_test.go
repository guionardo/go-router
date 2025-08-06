package path_params

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_createRegexFromPattern(t *testing.T) {
	tests := []struct {
		name          string
		pattern       string
		expectedNames []string
		testURL       string
		shouldMatch   bool
	}{
		{
			name:          "no_pattern",
			pattern:       "/users",
			expectedNames: nil,
			testURL:       "/users",
			shouldMatch:   true,
		},
		{
			name:          "basic_pattern_gin",
			pattern:       "/users/{id}",
			expectedNames: []string{"id"},
			testURL:       "/users/123",
			shouldMatch:   true,
		},
		{
			name:          "multiple_parameters_gin",
			pattern:       "/users/{userId}/posts/{postId}",
			expectedNames: []string{"userId", "postId"},
			testURL:       "/users/123/posts/456",
			shouldMatch:   true,
		},
		{
			name:          "no match",
			pattern:       "/users/{id}",
			expectedNames: []string{"id"},
			testURL:       "/posts/123",
			shouldMatch:   false,
		},
		{
			name:          "basic pattern",
			pattern:       "/users/:id",
			expectedNames: []string{"id"},
			testURL:       "/users/123",
			shouldMatch:   true,
		},
		{
			name:          "multiple parameters",
			pattern:       "/users/:userId/posts/:postId",
			expectedNames: []string{"userId", "postId"},
			testURL:       "/users/123/posts/456",
			shouldMatch:   true,
		},
		{
			name:          "no match",
			pattern:       "/users/:id",
			expectedNames: []string{"id"},
			testURL:       "/posts/123",
			shouldMatch:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			regex, names, err := createRegexFromPattern(tt.pattern)
			if !assert.NoError(t, err) {
				return
			}
			if !assert.Equal(t, tt.expectedNames, names) {
				return
			}

			if regex != nil {

				matches := regex.MatchString(tt.testURL)
				assert.Equal(t, tt.shouldMatch, matches)
				if tt.shouldMatch {
					assert.Equal(t, tt.shouldMatch, matches, "expected match %t for URL %s, got %t", tt.shouldMatch, tt.testURL, matches)
				}
			}
		})
	}
}
