package main

import (
	"os"

	"github.com/lil-yellow-flower/humctl-wrapper-demo/internal/commands"
)

func main() {
	if err := commands.Execute(); err != nil {
		os.Exit(1)
	}
} 