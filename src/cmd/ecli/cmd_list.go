package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"helm.sh/helm/v3/pkg/chart/loader"
)

func init() {
	rootCmd.AddCommand(&cobra.Command{
		Use:   "list [chart]",
		Short: "List properties in chart",
		Long:  "List properties in chart. Defaults to current working directory.",
		RunE: func(cmd *cobra.Command, args []string) error {
			var (
				workDir string
			)

			if len(args) == 0 {
				workDir, _ = os.Getwd()
			} else {
				workDir = args[0]
			}

			return listCmd(workDir)
		},
	})
}

func listCmd(chart string) error {
	helmChart, err := loader.Load(chart)

	if err != nil {
		return fmt.Errorf("failed to load chart %s: %w\n", chart, err)
	}

	for annotation, value := range helmChart.Metadata.Annotations {
		if strings.HasPrefix(annotation, "ecida") {
			fmt.Printf("%s: %s\n", annotation, value)
		}
	}

	return nil
}
