package main

import (
	"enigma/internal/api"
	"enigma/internal/config"
	"log"
)

func main() {
	cfg := config.NewSeverConfig()

	if err := api.Run(cfg); err != nil {
		log.Fatal(err)
	}
}
