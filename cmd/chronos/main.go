package main

import (
	"log"
	"os"

	"github.com/flohansen/chronos/internal/cli"
)

func run() error {
	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "init":
			return cli.Init()
		}
	} else {
		return cli.Run()
	}

	return nil
}

func main() {
	if err := run(); err != nil {
		log.Fatalf("error: %s", err)
	}
}
