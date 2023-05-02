package main

import (
	"log"

	"github.com/flowshot-io/navigator-cli/internal/cli"
)

func main() {
	app, err := cli.New()
	if err != nil {
		log.Fatal(err)
	}

	if err := app.Execute(); err != nil {
		return
	}
}
