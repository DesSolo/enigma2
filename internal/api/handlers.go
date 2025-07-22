package api

import (
	"errors"
	"fmt"
	"html"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"enigma/internal/pkg/storage"
)

func indexHandler(template []byte) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		if _, err := rw.Write(template); err != nil {
			slog.ErrorContext(r.Context(), "fault write response", "err", err)
		}
	}
}

func createSecretHandler(p SecretsProvider, externalURL string) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
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

		token, err := p.SaveSecret(ctx, msgFormValue, dues)
		if err != nil {
			slog.ErrorContext(ctx, "fault save secret", "err", err)
			raiseError(rw, http.StatusInternalServerError)
			return
		}

		if _, err := fmt.Fprintf(rw, "%s/get/%s", externalURL, token); err != nil {
			slog.ErrorContext(ctx, "fault write response", "err", err)
			raiseError(rw, http.StatusInternalServerError)
			return
		}
	}
}

func viewSecretHandler(p SecretsProvider, template []byte) http.HandlerFunc {
	tpl := string(template)

	return func(rw http.ResponseWriter, r *http.Request) {
		token := chi.URLParam(r, "token")
		if token == "" {
			raiseError(rw, http.StatusNotFound)
			return
		}

		ctx := r.Context()

		secret, err := p.GetSecret(ctx, token)
		if err != nil {
			if errors.Is(err, storage.ErrNotFound) {
				raiseError(rw, http.StatusNotFound)
				return
			}

			slog.ErrorContext(ctx, "fault get secret", "err", err)
			raiseError(rw, http.StatusInternalServerError)
			return
		}

		if _, err := fmt.Fprintf(rw, tpl, html.EscapeString(secret)); err != nil {
			slog.ErrorContext(ctx, "fault return response", "err", err)
			raiseError(rw, http.StatusInternalServerError)
		}
	}
}

func healthHandler() http.HandlerFunc {
	return func(rw http.ResponseWriter, _ *http.Request) {
		rw.Header().Set("Content-Type", "application/json")
		rw.WriteHeader(http.StatusOK)
		rw.Write([]byte(`{"healthy": true}`)) // nolint:errcheck
	}
}
