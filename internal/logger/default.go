package logger

import (
	"log/slog"
	"os"
	"time"

	"github.com/lmittmann/tint"
)

// init initializes the logger with predefined settings and returns the logger instance.
func init() {
	// Create a new handler with tint colorized output
	handler := tint.NewHandler(os.Stderr, &tint.Options{
		AddSource:  true,
		Level:      slog.LevelInfo,
		TimeFormat: time.Kitchen,
	})

	// Wrap the tint handler with otel handler if otel integration is enabled (optional)
	handler = newHandler(handler)

	// Create a new slog logger instance with the formatted handler
	l := slog.New(handler)

	// Set the newly created logger as the default logger for the application
	slog.SetDefault(l)
}
