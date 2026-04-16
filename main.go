package main

import (
	"fmt"
	"os"

	"bizzmod-cli/cmd"
)

var version = "dev"

func main() {
	if err := cmd.NewRootCmd(version).Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
