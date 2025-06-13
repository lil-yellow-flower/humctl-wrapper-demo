package main

import (
	"log"
	"os"

	"github.com/lil-yellow-flower/humctl-wrapper-demo/internal/commands"
)

func main() {
	if err := commands.Execute(); err != nil {
		log.Printf("Error: %v\n", err)
		os.Exit(1)
	}
} 