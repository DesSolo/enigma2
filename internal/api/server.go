package api

import (
	"context"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
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
func (s *Server) Run(addr string) error {
	return http.ListenAndServe(addr, s.router) // nolint:gosec,wrapcheck
}
