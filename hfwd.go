package main

import (
	"log"
	"os"

	"github.com/kei2100/h-fwd/cli"
)

func main() {
	if err := cli.RootCmd.Execute(); err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
}
