package main

import (
	"log"

	"github.com/criticalstack/crit/cmd/cinder/app"
)

func main() {
	if err := app.NewCommand().Execute(); err != nil {
		log.Fatal(err)
	}
}
