package logger

import (
	"context"
	"log/slog"

	"go.opentelemetry.io/otel/trace"

	"github.com/angelokurtis/go-home-automations/internal/errors"
)

// Handler wraps a slog.Handler and injects trace information into logs
type Handler struct {
	slog.Handler
}

// newHandler creates a new Handler instance
func newHandler(handler slog.Handler) *Handler {
	return &Handler{Handler: handler}
}

// Handle processes a log record and injects trace context if available
func (h *Handler) Handle(ctx context.Context, r slog.Record) error {
	// Check if trace context exists in the context
	if span := trace.SpanContextFromContext(ctx); span.IsValid() {
		// Extract trace and span IDs from the context
		traceID := span.TraceID().String()
		spanID := span.SpanID().String()

		// Add trace and span IDs as attributes to the log record
		r.AddAttrs(
			slog.String("trace-id", traceID),
			slog.String("span-id", spanID),
		)
	}

	// Wrap the original handler's Handle call with error handling
	return errors.WithStack(h.Handler.Handle(ctx, r))
}
