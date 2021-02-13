package main

import (
	"enigma/config"
	"enigma/internal/api"
	"log"
)

func main() {
	config := config.NewSeverConfig()

	if err := api.Run(config); err != nil {
		log.Fatal(err)
	}
}
