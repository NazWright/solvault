package main

import (
	"os"

	"github.com/NazWright/solvault/cmd/solvault/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
