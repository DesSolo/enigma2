package api

import (
	"fmt"
	"html"
	"log"
	"net/http"
	"regexp"
	"strconv"
)

var (
	templateIndex []byte
	templateGet   []byte

	viewURLPattern *regexp.Regexp
)

func indexHandler(w http.ResponseWriter, r *http.Request) {
	if _, err := w.Write(templateIndex); err != nil {
		log.Println("error return response err:", err)
	}
}

func createHandler(w http.ResponseWriter, r *http.Request) {
	msgFormValue := r.FormValue("msg")
	dueFormValue := r.FormValue("due")
	if len(msgFormValue) == 0 || len(dueFormValue) == 0 || len(msgFormValue) >= 65535 {
		raiseError(w, http.StatusBadRequest)
		return
	}

	dues, err := strconv.Atoi(dueFormValue)
	if err != nil || !(dues >= 1 && dues <= 4) {
		raiseError(w, http.StatusBadRequest)
		return
	}

	secret, err := saveSecret(sStorage, msgFormValue, dues)
	if err != nil {
		log.Println("error save secret err:", err)
		raiseError(w, http.StatusInternalServerError)
		return
	}

	if _, err := fmt.Fprintf(w, "%s/get/%s", sConfig.ResponseAddress, secret); err != nil {
		log.Println("error return response err:", err)
		raiseError(w, http.StatusInternalServerError)
		return
	}
}

func viewHandler(w http.ResponseWriter, r *http.Request) {
	key := viewURLPattern.FindStringSubmatch(r.URL.Path)
	if len(key) == 0 {
		raiseError(w, http.StatusNotFound)
		return
	}

	secret, err := getSecret(sStorage, key[1])
	if err != nil {
		log.Println("error get secret err:", err)
		raiseError(w, http.StatusNotFound)
		return
	}
	if secret == "" {
		raiseError(w, http.StatusNotFound)
		return
	}

	if _, err := fmt.Fprintf(w, string(templateGet), html.EscapeString(secret)); err != nil {
		log.Println("error return response err:", err)
		raiseError(w, http.StatusInternalServerError)
	}
}
