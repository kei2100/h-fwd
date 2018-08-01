package main

import (
	"log"
	"os"

	"github.com/kei2100/fwxy/cli"
)

func main() {
	if err := cli.RootCmd.Execute(); err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
}
