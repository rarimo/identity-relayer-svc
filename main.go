package main

import (
	"os"

	"gitlab.com/rarify-protocol/relayer-svc/internal/cli"
)

func main() {
	cli.Run(os.Args)
}
