package main

import (
	"os"

	"gitlab.com/rarimo/relayer-svc/internal/cli"
)

func main() {
	cli.Run(os.Args)
}
