package logging

import "log/slog"

type (
	Logger struct {
		infoFunc  logFunc
		debugFunc logFunc
		warnFunc  logFunc
	}
	logFunc func(message string, args ...any)
)

func New(enabled bool) *Logger {
	bf := func(f logFunc) logFunc {
		if enabled {
			return f
		}
		return func(string, ...any) {}
	}
	return &Logger{
		infoFunc:  bf(slog.Info),
		debugFunc: bf(slog.Debug),
		warnFunc:  bf(slog.Warn),
	}
}

func (l *Logger) Info(message string, args ...any) {
	l.infoFunc(message, args...)
}

func (l *Logger) Debug(message string, args ...any) {
	l.debugFunc(message, args...)
}

func (l *Logger) Warn(message string, args ...any) {
	l.warnFunc(message, args...)
}
