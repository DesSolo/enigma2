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
	di := newContainer()

	configureLogger(di)

	slog.Info("server running", "addr", di.Config().Server.Bind)

	if err := di.APIServer().Run(di.Config().Server.Bind); err != nil {
		return fmt.Errorf("failed to start API server: %w", err)
	}

	return nil
}
