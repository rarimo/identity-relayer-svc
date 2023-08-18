package main

import (
	"os"

	"gitlab.com/rarimo/relayer-svc/internal/cli"
)

func main() {
	if !cli.Run(os.Args) {
		os.Exit(1)
	}
}
