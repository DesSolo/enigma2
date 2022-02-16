package api

import (
	"enigma/internal/api/service"
	"fmt"
	"html"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

func indexHandler(template []byte) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		if _, err := rw.Write(template); err != nil {
			log.Printf("fault write response err: %s", err.Error())
		}
	}
}

func createSecretHandler(s *service.SecretService, externalURL string) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		msgFormValue := r.FormValue("msg")
		dueFormValue := r.FormValue("due")
		if len(msgFormValue) == 0 || len(dueFormValue) == 0 || len(msgFormValue) >= 65535 {
			raiseError(rw, http.StatusBadRequest)
			return
		}

		dues, err := strconv.Atoi(dueFormValue)
		if err != nil || !(dues >= 1 && dues <= 4) {
			raiseError(rw, http.StatusBadRequest)
			return
		}

		token, err := s.Save(msgFormValue, dues)
		if err != nil {
			log.Printf("fault save secret err: %s", err.Error())
			raiseError(rw, http.StatusInternalServerError)
			return
		}

		if _, err := fmt.Fprintf(rw, "%s/get/%s", externalURL, token); err != nil {
			log.Printf("fault write response err: %s", err.Error())
			raiseError(rw, http.StatusInternalServerError)
			return
		}
	}
}

func viewSecretHandler(s *service.SecretService, template []byte) http.HandlerFunc {
	tpl := string(template)
	// todo bench this

	return func(rw http.ResponseWriter, r *http.Request) {
		token := chi.URLParam(r, "token")
		if token == "" {
			raiseError(rw, http.StatusNotFound)
			return
		}

		secret, err := s.Get(token)
		if err != nil {
			log.Printf("fault get secret err: %s", err.Error())
			raiseError(rw, http.StatusNotFound)
			return
		}

		if secret == "" {
			raiseError(rw, http.StatusNotFound)
			return
		}

		if _, err := fmt.Fprintf(rw, tpl, html.EscapeString(secret)); err != nil {
			log.Printf("fault return response err: %s", err.Error())
			raiseError(rw, http.StatusInternalServerError)
		}
	}
}
