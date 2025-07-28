package api

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"

	"enigma/internal/pkg/adapters/template"
)

func (s *Server) shouldRenderTemplate(ctx context.Context, rw http.ResponseWriter, tmpl string, data template.Data) {
	if err := s.renderTemplate(rw, tmpl, data); err != nil {
		slog.ErrorContext(ctx, "fault render template", "err", err)
		raiseError(rw, http.StatusInternalServerError)
	}
}

func (s *Server) renderTemplate(rw http.ResponseWriter, name string, data template.Data) error {
	if err := s.template.RenderFile(name, rw, data); err != nil {
		return fmt.Errorf("could not render template: %w", err)
	}

	return nil
}

func raiseError(w http.ResponseWriter, method int) {
	http.Error(w, http.StatusText(method), method)
}
