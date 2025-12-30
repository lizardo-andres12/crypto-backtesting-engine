package logger

import (
	"context"
	"log/slog"
)

type contextKey string

const RunIDKey contextKey = "run_id"

// ContextHandler wraps slog.Handler and adds the run_id extraction from context
type ContextHandler struct {
	slog.Handler
}

// Handle extracts the run_id from the injected context and logs it if present
func (ch *ContextHandler) Handle(ctx context.Context, record slog.Record) error {
	if v := ctx.Value(RunIDKey); v != nil {
		record.AddAttrs(slog.String("run_id", v.(string)))
	}
	return ch.Handler.Handle(ctx, record)
}

