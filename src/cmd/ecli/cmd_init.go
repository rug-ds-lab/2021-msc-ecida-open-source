package main

import (
	"fmt"

    "ecida/pkg/meta"

	"helm.sh/helm/v3/pkg/chart"
	"helm.sh/helm/v3/pkg/chartutil"
)

func NewEcidaChart(name string) *chart.Metadata {
    metadata := meta.NewEcidaMetadata()

    metadata["ischart"] = "true"

    return &chart.Metadata{
        Name: name,
        Description: fmt.Sprintf("The %s ECiDA helmchart chart", name),
        APIVersion: "v2",
        Version: "0.1.0",
        Type: "application",
        Annotations: metadata.ToChartAnnotations(),
    }
}

func initCmd(name string, dirname string) error {
    fmt.Printf("initting %s and %s\n", name, dirname)

    path, err := chartutil.Create(name, dirname)

    if err != nil {
        return fmt.Errorf("failed to create new chart: %w", err)
    }

    chartfile := fmt.Sprintf("%s/Chart.yaml", path)

    cf := NewEcidaChart(name)

    return chartutil.SaveChartfile(chartfile, cf)
}
