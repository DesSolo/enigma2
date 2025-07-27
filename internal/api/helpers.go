package api

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/flosch/pongo2/v6"
)

func (s *Server) shouldRenderTemplate(ctx context.Context, rw http.ResponseWriter, tmpl string, pongoContext pongo2.Context) {
	if err := s.renderTemplate(rw, tmpl, pongoContext); err != nil {
		slog.ErrorContext(ctx, "fault render template", "err", err)
		raiseError(rw, http.StatusInternalServerError)
	}
}

func (s *Server) renderTemplate(rw http.ResponseWriter, name string, pongoContext pongo2.Context) error {
	tmpl, err := s.templateSet.FromFile(name)
	if err != nil {
		return fmt.Errorf("could not load template: %w", err)
	}

	if err := tmpl.ExecuteWriter(pongoContext, rw); err != nil {
		return fmt.Errorf("could not render template: %w", err)
	}

	return nil
}

func raiseError(w http.ResponseWriter, method int) {
	http.Error(w, http.StatusText(method), method)
}
