package main

import (
	"os"

	"github.com/leo201313/Blockchain_with_Go/cli"
)

func main() {
	defer os.Exit(0)
	cmd := cli.CommandLine{}
	cmd.Run()
}
