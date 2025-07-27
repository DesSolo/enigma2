package api

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/flosch/pongo2/v6"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

const (
	defaultReadTimeout = 5 * time.Second
)

// SecretsProvider ...
type SecretsProvider interface {
	SaveSecret(ctx context.Context, message string, dues int) (string, error)
	CheckExistsSecret(ctx context.Context, token string) error
	GetSecret(ctx context.Context, token string) (string, error)
}

// Server ...
type Server struct {
	secretsProvider SecretsProvider
	externalURL     string

	templateSet *pongo2.TemplateSet
	router      *chi.Mux
}

// NewServer ...
func NewServer(secretsProvider SecretsProvider, templateSet *pongo2.TemplateSet, externalURL string) *Server {
	return &Server{
		secretsProvider: secretsProvider,
		externalURL:     externalURL,
		templateSet:     templateSet,
		router:          chi.NewRouter(),
	}
}

// Run ...
func (s *Server) Run(ctx context.Context, addr string) error {
	s.initHandlers()

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

func (s *Server) initHandlers() {
	s.router.Use(middleware.Recoverer)
	s.router.Use(middleware.RequestID)
	s.router.Use(middleware.Logger)

	s.router.Get("/", s.indexHandler)
	s.router.Post("/post/", s.createSecretHandler)
	s.router.Get("/get/{token}", s.viewSecretHandler)

	s.router.Get("/health", healthHandler())
}
