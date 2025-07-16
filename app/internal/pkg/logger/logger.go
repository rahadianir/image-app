package logger

import (
	"context"
	"log/slog"
	"os"
	"time"
)

type CustomHandler struct {
	handler slog.Handler
}

func NewLogger() *slog.Logger {
	return slog.New(&CustomHandler{
		handler: slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
				if a.Key == slog.MessageKey {
					a.Key = "message"
				}

				if a.Key == slog.TimeKey {
					a.Value = slog.StringValue(a.Value.Time().Format(time.RFC3339))
				}

				return a
			},
		}),
	})
}

func (h *CustomHandler) Handle(ctx context.Context, r slog.Record) error {
	return h.handler.Handle(ctx, r)
}

func (h *CustomHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return h.handler.Enabled(ctx, level)
}

func (h *CustomHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &CustomHandler{
		handler: h.handler.WithAttrs(attrs),
	}
}

func (h *CustomHandler) WithGroup(name string) slog.Handler {
	return &CustomHandler{
		handler: h.handler.WithGroup(name),
	}
}
