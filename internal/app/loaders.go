package app

import (
	"fmt"
	"log/slog"
	"os"
)

func configureLogger(di *container) error {
	options := di.Config().Logging

	handlerOptions := &slog.HandlerOptions{
		Level: slog.Level(options.Level),
	}

	var handler slog.Handler

	switch options.Format {
	case "text":
		handler = slog.NewTextHandler(os.Stdout, handlerOptions)
	case "json":
		handler = slog.NewJSONHandler(os.Stdout, handlerOptions)
	default:
		return fmt.Errorf("unsupported log format: %s", options.Format)
	}

	slog.SetDefault(
		slog.New(handler),
	)

	return nil
}
