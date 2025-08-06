package inspect

import (
	"fmt"
	"strconv"
	"time"
)

func stringToBool(v string) bool {
	bv, _ := strconv.ParseBool(v)
	return bv
}

func stringToInt(v string) int {
	iv, _ := strconv.Atoi(v)
	return iv
}

func boolValue[T comparable](v bool, vTrue, vFalse T) T {
	if v {
		return vTrue
	}
	return vFalse
}

func strToTime(v string) (time.Time, error) {
	for _, templates := range []string{time.RFC1123, time.RFC3339, time.DateTime, time.DateOnly, time.TimeOnly} {
		if tv, err := time.Parse(templates, v); err == nil {
			return tv, nil
		}
	}
	return time.Time{}, fmt.Errorf("invalid time.Time value: '%s'", v)
}
