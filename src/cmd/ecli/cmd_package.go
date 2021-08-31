package main

import (
	"fmt"

	"helm.sh/helm/v3/pkg/chart/loader"
	"helm.sh/helm/v3/pkg/chartutil"
)

func packageCmd(chart string, destination string) error {
	helmChart, err := loader.Load(chart)

	if err != nil {
		return fmt.Errorf("failed to load chart: %w", err)
	}

	path, err := chartutil.Save(helmChart, destination)

	if err != nil {
		return fmt.Errorf("failed to export chart to %s: %w", destination, err)
	}

	fmt.Printf("wrote package to %s\n", path)

	return nil
}
