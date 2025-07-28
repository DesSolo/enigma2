//go:generate mockery --case snake --with-expecter --name SecretsProvider

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
	slogchi "github.com/samber/slog-chi"

	"enigma/internal/pkg/adapters/template"
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
	template        template.Template
	externalURL     string
	router          *chi.Mux
}

// NewServer ...
func NewServer(secretsProvider SecretsProvider, template template.Template, externalURL string) *Server {
	return &Server{
		secretsProvider: secretsProvider,
		externalURL:     externalURL,
		template:        template,
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
	s.router.Use(middleware.RequestID)
	s.router.Use(slogchi.New(slog.Default()))
	s.router.Use(middleware.Recoverer)

	s.router.Get("/", s.indexHandler)
	s.router.Post("/post/", s.createSecretHandler)
	s.router.Get("/get/{token}", s.viewSecretHandler)

	s.router.Get("/health", healthHandler())
}
