package main

import (
	"os"

	"github.com/Levi9111/create-express-modular-go/cmd"
)

func main() {
	if err := cmd.RootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
