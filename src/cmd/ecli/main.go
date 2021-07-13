package main

import (
	"fmt"
	"log"

	"github.com/docopt/docopt-go"
)

func run() error {
	usage := `ECiDA CLI tool

Usage:
  ecli list <chart>`

	opts, err := docopt.ParseDoc(usage)

	if err != nil {
		return err
	}

	if list, _ := opts.Bool("list"); list {
        chart, _ := opts.String("<chart>")
		fmt.Printf("listing %s\n", chart)
	}

	return nil
}

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}
