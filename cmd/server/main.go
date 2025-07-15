package main

import (
	"context"
	"fmt"
	"log"

	"enigma/internal/app"
)

var version = "local"

const banner = `
 _______ __   _ _____  ______ _______ _______
 |______ | \  |   |   |  ____ |  |  | |_____|
 |______ |  \_| __|__ |_____| |  |  | |     |

 version: %s

`

func main() {
	application := app.New()

	fmt.Printf(banner, version)

	// TODO: grace
	if err := application.Run(context.Background()); err != nil {
		log.Fatalf("failed to run applications err: %s", err.Error())
	}
}
