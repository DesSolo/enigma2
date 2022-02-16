package api

import (
	"net/http"
)

func raiseError(w http.ResponseWriter, method int) {
	http.Error(w, http.StatusText(method), method)
}
