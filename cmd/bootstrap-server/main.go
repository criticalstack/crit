package main

import (
	"log"

	"github.com/criticalstack/crit/cmd/bootstrap-server/app"
)

func main() {
	if err := app.NewRootCmd().Execute(); err != nil {
		log.Fatal(err)
	}
}
