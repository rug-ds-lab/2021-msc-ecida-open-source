package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/docopt/docopt-go"

	"helm.sh/helm/v3/pkg/chart/loader"
	"helm.sh/helm/v3/pkg/chartutil"
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
		fmt.Printf("listing %s\n", chart)

        helmChart, err := loader.Load(chart)

        if err != nil {
            return fmt.Errorf("failed to load chart: %w", err)
        }

        for annotation, value := range helmChart.Metadata.Annotations {
            if strings.HasPrefix(annotation, "ecida") {
                fmt.Printf("%s: %s\n", annotation, value)
            }
        }
	}

    if pkg, _ := opts.Bool("package"); pkg {
        chart, _ := opts.String("<chart>")
        destination, _ := opts.String("<destination>")

        helmChart, err := loader.Load(chart)

        if err != nil {
            return fmt.Errorf("failed to load chart: %w", err)
        }

        path, err := chartutil.Save(helmChart, destination)

        if err != nil {
            return fmt.Errorf("failed to export chart to %s: %w", destination, err)
        }
        
        fmt.Printf("wrote package to %s\n", path)
    }

	return nil
}

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}
