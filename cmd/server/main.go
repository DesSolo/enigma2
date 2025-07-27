package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"os/signal"
	"syscall"

	"enigma/internal/app"
	"enigma/pkg/closer"
)

var version = "local"

const banner = `
 _______ __   _ _____  ______ _______ _______
 |______ | \  |   |   |  ____ |  |  | |_____|
 |______ |  \_| __|__ |_____| |  |  | |     |

 version: %s

`

func main() {
	application := app.New()

	fmt.Printf(banner, version) // nolint:forbidigo

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	go func() {
		<-ctx.Done()
		slog.DebugContext(ctx, "shutting down")

		cancel()

		if err := closer.Close(); err != nil {
			slog.ErrorContext(ctx, "closer.Close", "err", err)
		}
	}()

	if err := application.Run(ctx); err != nil {
		log.Fatalf("failed to run applications err: %s", err.Error())
	}
}
