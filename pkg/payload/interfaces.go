package payload

import "io"

type Payloader interface {
	Marshal(w io.Writer) error
	Unmarshal(r io.Reader) error
}
