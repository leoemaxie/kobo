package telemetry

import (
	"log/slog"
	"os"
)

// InitLogger configures the standard library structured logger.
func InitLogger() {
	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if err, ok := a.Value.Any().(error); ok {
				return slog.String(a.Key, err.Error())
			}
			return a
		},
	})
	logger := slog.New(handler)
	slog.SetDefault(logger)
}
