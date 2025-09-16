package tools

import (
	"fmt"
	"time"
)

var layouts = NewLastSuccessIterator(time.DateTime, time.RFC3339, time.RFC3339Nano, time.RFC1123, time.RFC1123Z, time.ANSIC, time.DateOnly, time.TimeOnly)

func ParseTimeLayouts(s string) (t time.Time, err error) {
	for layout := range layouts.Iter() {
		if t, err = time.Parse(layout, s); err == nil {
			break
		}
	}
	if err != nil {
		err = fmt.Errorf("time \"%s\" couldn't be parsed", s)
	}
	return t, err
}
