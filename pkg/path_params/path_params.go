package path_params

import (
	"fmt"
	"net/url"
	"regexp"

	"github.com/lestrrat-go/urlenc"
)

type PathParams struct {
	pattern     string
	params      map[string]string
	paramsNames []string
	regex       *regexp.Regexp
	matchFunc   func(url *url.URL) error
	injectFunc  func(payload any) error
}

// NewPathParams creates a new PathParams instance
func NewPathParams(pattern string) (*PathParams, error) {
	regex, paramNames, err := createRegexFromPattern(pattern)
	if err != nil {
		return nil, err
	}

	pp := &PathParams{
		pattern:     pattern,
		params:      make(map[string]string, len(paramNames)),
		paramsNames: paramNames,
		regex:       regex,
		injectFunc:  nil,
	}
	if regex == nil {
		pp.injectFunc = func(payload any) error {
			return nil
		}
		pp.matchFunc = func(url *url.URL) error {
			if url.Path == pp.pattern {
				return nil
			}
			return fmt.Errorf("url [%s] does not match the pattern [%s]", url.Path, pp.pattern)
		}
	} else {
		pp.injectFunc = func(payload any) error {
			d, err := urlenc.Marshal(pp.params)
			if err != nil {
				return err
			}
			return urlenc.Unmarshal(d, payload)
		}
		pp.matchFunc = pp.match
	}

	params := make(map[string]string, len(paramNames))
	for _, name := range paramNames {
		if _, ok := params[name]; ok {
			return nil, fmt.Errorf("duplicated path param name: %s - pattern: %s", name, pattern)
		}
		params[name] = ""
	}
	return pp, nil
}

func (p *PathParams) match(url *url.URL) error {
	if p.regex == nil {
		if url.Path == p.pattern {
			return nil
		}
		return fmt.Errorf("url [%s] does not match the pattern [%s]", url.Path, p.pattern)
	}
	matches := p.regex.FindStringSubmatch(url.Path)
	if matches == nil {
		return fmt.Errorf("url [%s] does not match the pattern [%s]", url, p.regex.String())
	}

	// matches[0] is the full match, matches[1:] are the capture groups
	for i, name := range p.paramsNames {
		if i+1 < len(matches) {
			p.params[name] = matches[i+1]
		}
	}
	return nil
}

func (p *PathParams) Match(url *url.URL) error {
	return p.matchFunc(url)
}

func (p *PathParams) Get(key string) string {
	return p.params[key]
}

func (p *PathParams) Set(key, value string) {
	if p.params == nil {
		p.params = make(map[string]string)
	}
	p.params[key] = value
}

// GetAll returns all parameters
func (p *PathParams) GetAll() map[string]string {
	return p.params
}

func (p *PathParams) Inject(payload any) error {
	return p.injectFunc(payload)
}

func (p *PathParams) MatchAndInject(url *url.URL, payload any) (err error) {
	if err = p.Match(url); err == nil {
		err = p.Inject(payload)
	}
	return err
}

func (p *PathParams) ParamNames() []string {
	return p.paramsNames
}
