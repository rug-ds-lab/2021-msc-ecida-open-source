package main

import (
	"log"

	"github.com/docopt/docopt-go"
)

func run() error {
	usage := `ECiDA CLI tool

Usage:
  ecli list <chart>
  ecli package <chart> <destination>
  ecli publish <file>
  ecli init [--name <name>] [--dir <dirname>]

Options:
  -n <name>, --name <name>       The name of the package
  -d <dirname>, --dir <dirname>  Location of the module
`

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

    if publish, _ := opts.Bool("publish"); publish {
        file, _ := opts.String("<file>")
        
        return publishCmd(DirOrWorkdir(file))
    }

    if init, _ := opts.Bool("init"); init {
        name, _ := opts.String("--name")
        dirname, _ := opts.String("--dir")

        return initCmd(name, DirOrWorkdir(dirname))
    }

	return nil
}

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}
