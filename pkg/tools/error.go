package tools

import "fmt"

type ParseError struct {
	fieldName string
	err       error
}

func (pe *ParseError) Error() string {
	return fmt.Sprintf("%s - %s", pe.fieldName, pe.err.Error())
}

func NewParseError(fieldName string, err error) error {
	if err == nil {
		return nil
	}
	return &ParseError{
		fieldName: fieldName,
		err:       err,
	}
}

func GroupError(parseType string, funcs ...func() error) error {
	for _, f := range funcs {
		if err := f(); err != nil {
			return fmt.Errorf("%s - %w", parseType, err)
		}
	}
	return nil
}
