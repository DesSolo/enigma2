package api

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

const (
	defaultReadTimeout = 5 * time.Second
)

// SecretsProvider ...
type SecretsProvider interface {
	SaveSecret(ctx context.Context, message string, dues int) (string, error)
	GetSecret(ctx context.Context, message string) (string, error)
}

// Server ...
type Server struct {
	secretsProvider SecretsProvider
	router          *chi.Mux
}

// NewServer ...
func NewServer(secretsProvider SecretsProvider) *Server {
	return &Server{
		secretsProvider: secretsProvider,
		router:          chi.NewRouter(),
	}
}

// LoadHandlers ...
func (s *Server) LoadHandlers(indexTemplate, viewSecretTemplate []byte, externalURL string) {
	s.router.Use(middleware.Recoverer)
	s.router.Use(middleware.RequestID)
	s.router.Use(middleware.Logger)

	s.router.Get("/", indexHandler(indexTemplate))
	s.router.Post("/post/", createSecretHandler(s.secretsProvider, externalURL))
	s.router.Get("/get/{token}", viewSecretHandler(s.secretsProvider, viewSecretTemplate))

	s.router.Get("/health", healthHandler())
}

// Run ...
func (s *Server) Run(ctx context.Context, addr string) error {
	server := http.Server{
		Addr:        addr,
		Handler:     s.router,
		ReadTimeout: defaultReadTimeout,
	}

	go func() {
		<-ctx.Done()
		if err := server.Shutdown(ctx); err != nil {
			slog.Error("failed to shutdown server", "err", err)
		}
	}()

	if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf(" server.ListenAndServe: %w", err)
	}

	return nil
}
