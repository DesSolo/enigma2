package app

import (
	"context"
	"fmt"
	"log/slog"
)

// App ...
type App struct{}

// New ...
func New() *App {
	return &App{}
}

// Run ...
func (a *App) Run(ctx context.Context) error {
	di := newContainer(ctx)

	configureLogger(di)

	bindAddr := di.Config().Server.Bind

	slog.Info("server running", "addr", bindAddr)

	if err := di.APIServer().Run(ctx, bindAddr); err != nil {
		return fmt.Errorf("failed to start API server: %w", err)
	}

	return nil
}
