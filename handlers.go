package main

import (
	"fmt"
	"html"
	"log"
	"net/http"
	"strconv"
)

var (
	templateIndex = LoadTemplate("templates/index.html")
	templateGet = LoadTemplate("templates/get.html")

	viewURLPattern = GenRegexpGetView()
)

func indexHandler(w http.ResponseWriter, r *http.Request)  {
	if _, err := w.Write(templateIndex); err != nil {
		log.Println("error return response err:", err)
	}
}

func createHandler(w http.ResponseWriter, r *http.Request)  {
	msgFormValue := r.FormValue("msg")
	dueFormValue := r.FormValue("due")
	if len(msgFormValue) == 0 || len(dueFormValue) == 0 || len(msgFormValue) >= 65535 {
		http.Error(w, "400 - bad request", http.StatusBadRequest)
		return
	}
	dues, err := strconv.Atoi(dueFormValue)
	if err != nil || !(dues >= 1 && dues <= 4) {
		http.Error(w, "400 - bad request", http.StatusBadRequest)
		return
	}
	secret, err := SaveSecret(Config.SecretStorage, msgFormValue, dues)
	if err != nil {
		log.Println("error save secret err:", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	if _, err := fmt.Fprintf(w, "%s/get/%s", Config.ResponseAddress, secret); err != nil {
		log.Println("error return response err:", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
}

func viewHandler(w http.ResponseWriter, r *http.Request) {
	key := viewURLPattern.FindStringSubmatch(r.URL.Path)
	if len(key) == 0 {
		log.Println("url get miss match")
		http.Error(w, "404 - not found", http.StatusNotFound)
		return
	}
	secret, err := GetSecret(Config.SecretStorage, key[1])
	if err != nil {
		log.Println("error get secret err:", err)
		http.Error(w, "404 - not found", http.StatusNotFound)
		return
	}
	if secret == "" {
		http.Error(w, "404 - not found", http.StatusNotFound)
		return
	}
	if _, err := fmt.Fprintf(w, string(templateGet), html.EscapeString(secret)); err != nil {
		log.Println("error return response err:", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}
}
