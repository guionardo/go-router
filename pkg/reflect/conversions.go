package reflections

import (
	"fmt"
	"time"
)

func BoolValue[T comparable](v bool, vTrue, vFalse T) T {
	if v {
		return vTrue
	}
	return vFalse
}

func StrToTime(v string) (time.Time, error) {
	for _, templates := range []string{time.RFC1123, time.RFC3339, time.DateTime, time.DateOnly, time.TimeOnly} {
		if tv, err := time.Parse(templates, v); err == nil {
			return tv, nil
		}
	}
	return time.Time{}, fmt.Errorf("invalid time.Time value: '%s'", v)
}
