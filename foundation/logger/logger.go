package logger

import (
	"context"
	"io"
	"log"
	"log/slog"
	"strings"
)

type Logger struct {
	*slog.Logger
}

func New(w io.Writer, level string) *Logger {
	var l slog.Level
	switch strings.ToLower(level) {
	case "debug":
		l = slog.LevelDebug
	case "info":
		l = slog.LevelInfo
	case "warn":
		l = slog.LevelWarn
	case "error":
		l = slog.LevelError
	default:
		l = slog.LevelInfo
	}

	opts := &slog.HandlerOptions{
		Level: l,
	}

	handler := slog.NewJSONHandler(w, opts)
	return &Logger{
		Logger: slog.New(handler),
	}
}

func (l *Logger) NewStdLogger() *log.Logger {
	return slog.NewLogLogger(l.Handler(), slog.LevelInfo)
}

func (l *Logger) Infow(msg string, keysAndValues ...interface{}) {
	l.Info(msg, keysAndValues...)
}

func (l *Logger) Errorw(msg string, err error, keysAndValues ...interface{}) {
	if err != nil {
		keysAndValues = append(keysAndValues, "error", err)
	}
	l.Error(msg, keysAndValues...)
}

func (l *Logger) Debugw(msg string, keysAndValues ...interface{}) {
	l.Debug(msg, keysAndValues...)
}

func (l *Logger) With(args ...interface{}) *Logger {
	return &Logger{
		Logger: l.Logger.With(args...),
	}
}

func (l *Logger) WithContext(ctx context.Context) *Logger {
	// In slog, context is passed to the log method (InfoContext, etc.),
	// but we can also extract trace IDs here if using a specific handler.
	// For now, just return l.
	return l
}
