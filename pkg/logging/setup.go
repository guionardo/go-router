package logging

import (
	"log/slog"
	"sync"

	"github.com/guionardo/go-router/pkg/config"
)

var (
	defaultLogger *Logger
	Get           = sync.OnceValue(func() *Logger {
		if defaultLogger == nil {

			defaultLogger = New(config.LogEnabled)
			if config.LogEnabled {
				defaultLogger.Info("Router logging enabled")
			}
		}
		return defaultLogger
	})
)

func Set(logger *slog.Logger) {
	defaultLogger = &Logger{
		infoFunc:  logger.Info,
		debugFunc: logger.Debug,
		warnFunc:  logger.Warn,
	}
	Get = func() *Logger {
		return defaultLogger
	}
}
