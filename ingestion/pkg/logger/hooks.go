package logger

import (
	"log/slog"
	"os"
)

// TODO: Parse a config file (debug run or real run?)
// Setup initializes global logger with context awareness
func Setup() {
	jsonHandler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
		AddSource: true,
	})

	logger := slog.New(&ContextHandler{Handler: jsonHandler})
	slog.SetDefault(logger)
}

