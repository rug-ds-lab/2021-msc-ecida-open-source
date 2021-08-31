package main

import (
	"fmt"
	"strings"

	"helm.sh/helm/v3/pkg/chart/loader"
)

func listCmd(chart string) error {
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

	return nil
}
