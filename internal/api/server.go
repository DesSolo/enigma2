package api

import (
	"enigma/internal/api/service"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

// Server ...
type Server struct {
	secretService *service.SecretService
	router        *chi.Mux
}

// NewServer ...
func NewServer(s *service.SecretService) *Server {
	return &Server{
		secretService: s,
		router:        chi.NewRouter(),
	}
}

// LoadHandlers ...
func (s *Server) LoadHandlers(indexTemplate, viewSecretTemplate []byte, externalURL string) {
	s.router.Use(middleware.Recoverer)
	s.router.Use(middleware.RequestID)
	s.router.Use(middleware.Logger)

	s.router.Get("/", indexHandler(indexTemplate))
	s.router.Post("/post/", createSecretHandler(s.secretService, externalURL))
	s.router.Get("/get/{token}", viewSecretHandler(s.secretService, viewSecretTemplate))

	s.router.Get("/health", healthHandler())
}

// Run ...
func (s *Server) Run(addr string) error {
	return http.ListenAndServe(addr, s.router)
}
