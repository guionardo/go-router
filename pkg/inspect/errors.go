package inspect

type (
	ParseErrorStruct struct {
		Errors []string `json:"parsing_errors,omitempty"`
	}
	Unwrapper interface {
		Unwrap() []error
	}
)

func NewParseError(err error) *ParseErrorStruct {
	if err == nil {
		return nil
	}
	var errors []string
	if uw, ok := err.(Unwrapper); ok {
		errs := uw.Unwrap()
		if len(errs) > 0 {
			errors = make([]string, len(errs))
			for i, e := range errs {
				errors[i] = e.Error()
			}
		}
	} else {
		errors = []string{err.Error()}
	}
	return &ParseErrorStruct{errors}
}