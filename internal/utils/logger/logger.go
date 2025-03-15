package logger

import (
	"bridgebot/pkg/prettylog"
	"log/slog"
	"fmt"
)

var (
	Info  func(msg string, args ...any)
	Warn  func(msg string, args ...any)
	Error func(msg string, args ...any)
	Debug func(msg string, args ...any)

	Infof  func(format string, args ...any)
	Warnf  func(format string, args ...any)
	Errorf func(format string, args ...any)
	Debugf func(format string, args ...any)
)

func init() {
	handler := prettylog.NewHandler(&slog.HandlerOptions{
		Level:     slog.LevelDebug,
		AddSource: true,
	})
	logger := slog.New(handler)
	slog.SetDefault(logger)
	Info = logger.Info
	Warn = logger.Warn
	Error = logger.Error
	Debug = logger.Debug

	// Printf-like logging (formatted)
	Infof = func(format string, args ...any) {
		logger.Info(fmt.Sprintf(format, args...))
	}
	Warnf = func(format string, args ...any) {
		logger.Warn(fmt.Sprintf(format, args...))
	}
	Errorf = func(format string, args ...any) {
		logger.Error(fmt.Sprintf(format, args...))
	}
	Debugf = func(format string, args ...any) {
		logger.Debug(fmt.Sprintf(format, args...))
	}
}
