package main

import (
	"log"

	"github.com/mathi-ma51zaw/humctl-wrapper-demo/internal/commands"
)

func main() {
	if err := commands.Execute(); err != nil {
		log.Fatal(err)
	}
} 