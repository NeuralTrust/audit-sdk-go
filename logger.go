package audit

import (
	"log/slog"
	"os"
)

func newLogger(level LogLevel) *slog.Logger {
	var slogLevel slog.Level
	switch level {
	case LogLevelDebug:
		slogLevel = slog.LevelDebug
	case LogLevelWarn:
		slogLevel = slog.LevelWarn
	case LogLevelError:
		slogLevel = slog.LevelError
	default:
		slogLevel = slog.LevelInfo
	}

	return slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slogLevel,
	}))
}

