package app

import (
	"log/slog"
	"os"
)

func configureLogger(di *container) {
	handler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.Level(di.Config().Logging.Level),
	})

	slog.SetDefault(
		slog.New(handler),
	)
}
