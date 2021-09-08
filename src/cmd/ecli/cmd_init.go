package main

import (
	"fmt"

	"ecida/pkg/meta"

	"github.com/spf13/cobra"
	"helm.sh/helm/v3/pkg/chart"
	"helm.sh/helm/v3/pkg/chartutil"
)

var (
	initLocation string
)

func init() {
	cmdInit := &cobra.Command{
		Use:   "init <name>",
		Short: "Initialise a new ECiDA module",
		Long:  "Initialise a new ECiDA module as a Helm chart in the current working directory",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			name := args[0]
			return initCmd(name, initLocation)
		},
	}

	cmdInit.Flags().StringVarP(&initLocation, "directory", "d", "", "directory in which to initialise the module")

	rootCmd.AddCommand(cmdInit)
}

func NewEcidaChart(name string) *chart.Metadata {
	metadata := meta.NewEcidaMetadata()

	metadata["ischart"] = "true"

	return &chart.Metadata{
		Name:        name,
		Description: fmt.Sprintf("The %s ECiDA helmchart chart", name),
		APIVersion:  "v2",
		Version:     "0.1.0",
		Type:        "application",
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
