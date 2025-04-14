package logger

import (
	"bridgebot/pkg/prettylog"
	"fmt"
	"log/slog"
	"os"
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

	Fatal  func(msg string, args ...any)
	Fatalf func(format string, args ...any)
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

	// Fatal logging
	Fatal = func(msg string, args ...any) {
		logger.Error(fmt.Sprintf(msg, args...))
		os.Exit(1)
	}
	Fatalf = func(format string, args ...any) {
		logger.Error(fmt.Sprintf(format, args...))
		os.Exit(1)
	}

}
