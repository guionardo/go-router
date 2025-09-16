package config

import (
	"os"
	"slices"
	"strconv"
	"strings"
)

const (
	ROUTER_LOGGING = "ROUTER_LOGGING"
	ENVIRONMENT    = "ENVIRONMENT"
)

var (
	DevelopmentMode bool
	LogEnabled      bool
)

func init() {
	LogEnabled, _ = strconv.ParseBool(os.Getenv(ROUTER_LOGGING))
	DevelopmentMode = slices.Contains(
		[]string{"LOCAL", "DEV", "DEVELOPMENT"},
		strings.ToUpper(os.Getenv(ENVIRONMENT)))
}
