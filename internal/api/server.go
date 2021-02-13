package api

import (
	"enigma/config"
	"enigma/internal/storage"
	"fmt"
	"log"
	"net/http"
)

var (
	sConfig  *config.ServerConfig
	sStorage storage.SecretStorage
)

// Run ...
func Run(config *config.ServerConfig) error {
	sConfig = config
	sStorage = config.SecretStorage

	viewURLPattern = genRegexpGetView(sConfig.TokenBytes)
	templateIndex = loadTemplate("templates/index.html")
	templateGet = loadTemplate("templates/get.html")

	if err := awaitSecretStorage(sStorage); err != nil {
		return err
	}

	http.HandleFunc("/",
		methodMiddleware("GET", indexHandler),
	)
	http.HandleFunc("/post/",
		methodMiddleware("POST", createHandler),
	)
	http.HandleFunc("/get/",
		methodMiddleware("GET", viewHandler),
	)

	log.Printf(
		"service started port: %d response_address: %s token_bytes: %d\n",
		sConfig.ListenPort, sConfig.ResponseAddress, sConfig.TokenBytes,
	)

	return http.ListenAndServe(fmt.Sprintf(":%d", sConfig.ListenPort), nil)
}
