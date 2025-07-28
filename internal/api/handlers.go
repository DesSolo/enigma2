package api

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"enigma/internal/pkg/adapters/template"
	"enigma/internal/pkg/storage"
)

func healthHandler() http.HandlerFunc {
	return func(rw http.ResponseWriter, _ *http.Request) {
		rw.Header().Set("Content-Type", "application/json")
		rw.WriteHeader(http.StatusOK)
		rw.Write([]byte(`{"healthy": true}`)) // nolint:errcheck
	}
}

func (s *Server) indexHandler(rw http.ResponseWriter, r *http.Request) {
	s.shouldRenderTemplate(r.Context(), rw, "index.html", nil)
}

func (s *Server) createSecretHandler(rw http.ResponseWriter, r *http.Request) {
	msgFormValue := r.FormValue("msg")
	dueFormValue := r.FormValue("due")
	if len(msgFormValue) == 0 || len(dueFormValue) == 0 || len(msgFormValue) >= 65535 {
		slog.Warn("not valid",
			"len_msgFormValue", len(msgFormValue),
			"dueFormValue", dueFormValue,
		)
		raiseError(rw, http.StatusBadRequest)
		return
	}

	dues, err := strconv.Atoi(dueFormValue)
	if err != nil || (dues < 1 || dues > 4) {
		slog.Warn("not valid dues",
			"dues", dues,
			"err", err,
		)
		raiseError(rw, http.StatusBadRequest)
		return
	}

	ctx := r.Context()

	token, err := s.secretsProvider.SaveSecret(ctx, msgFormValue, dues)
	if err != nil {
		slog.ErrorContext(ctx, "fault save secret", "err", err)
		raiseError(rw, http.StatusInternalServerError)
		return
	}

	if _, err := fmt.Fprintf(rw, "%s/get/%s", s.externalURL, token); err != nil {
		slog.ErrorContext(ctx, "fault write response", "err", err)
		raiseError(rw, http.StatusInternalServerError)
		return
	}
}

func (s *Server) viewSecretHandler(rw http.ResponseWriter, r *http.Request) {
	token := chi.URLParam(r, "token")
	if token == "" {
		raiseError(rw, http.StatusNotFound)
		return
	}

	conform := r.URL.Query().Get("conform")
	if conform == "" {
		s.viewConformPage(token).ServeHTTP(rw, r)
		return
	}

	s.viewSecretPage(token).ServeHTTP(rw, r)
}

func (s *Server) viewConformPage(token string) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		if err := s.secretsProvider.CheckExistsSecret(ctx, token); err != nil {
			if errors.Is(err, storage.ErrNotFound) {
				rw.WriteHeader(http.StatusNotFound)
				s.shouldRenderTemplate(ctx, rw, "expired.html", nil)
				return
			}

			slog.ErrorContext(ctx, "fault check exist secret", "err", err)
			raiseError(rw, http.StatusInternalServerError)
			return
		}

		s.shouldRenderTemplate(ctx, rw, "conform.html", nil)
	}
}

func (s *Server) viewSecretPage(token string) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		secret, err := s.secretsProvider.GetSecret(ctx, token)
		if err != nil {
			if errors.Is(err, storage.ErrNotFound) {
				rw.WriteHeader(http.StatusNotFound)
				s.shouldRenderTemplate(ctx, rw, "expired.html", nil)
				return
			}

			slog.ErrorContext(ctx, "fault get secret", "err", err)
			raiseError(rw, http.StatusInternalServerError)
			return
		}

		s.shouldRenderTemplate(ctx, rw, "secret.html", template.Data{
			"secret": secret,
		})
	}
}
