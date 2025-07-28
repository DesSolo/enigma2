package api

import (
	"context"
	"log/slog"
	"net/http"

	"enigma/internal/pkg/adapters/template"
)

func (s *Server) shouldRenderTemplate(ctx context.Context, rw http.ResponseWriter, name string, data template.Data) {
	if err := s.template.RenderFile(name, rw, data); err != nil {
		slog.ErrorContext(ctx, "fault render template", "err", err)
		raiseError(rw, http.StatusInternalServerError)
	}
}

func raiseError(w http.ResponseWriter, method int) {
	http.Error(w, http.StatusText(method), method)
}
