package main

import (
	"log"

	"github.com/docopt/docopt-go"
)

func run() error {
	usage := `ECiDA CLI tool

Usage:
  ecli list <chart>
  ecli package <chart> <destination>`

	opts, err := docopt.ParseDoc(usage)

	if err != nil {
		return err
	}

	if list, _ := opts.Bool("list"); list {
		chart, _ := opts.String("<chart>")

		return listCmd(chart)
	}

	if pkg, _ := opts.Bool("package"); pkg {
		chart, _ := opts.String("<chart>")
		destination, _ := opts.String("<destination>")

		return packageCmd(chart, destination)
	}

	return nil
}

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}
