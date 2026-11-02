package logger

import "log/slog"

func New(name string) *slog.Logger {
	return slog.Default()
}
