package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"helm.sh/helm/v3/pkg/chart/loader"
	"helm.sh/helm/v3/pkg/chartutil"
)

var (
	pkgLocation string
)

func init() {
	cmdPackage := &cobra.Command{
		Use:   "package [chart]",
		Short: "packages an ECiDA module",
        Long: `Packages a module found in [chart]. Defaults to working directory. Follows the
Helm naming convention where the packaged name will be chart-version.gz`,
        Args: cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {

			var chart string

			if len(args) == 0 {
				chart, _ = os.Getwd()
			} else {
				chart = args[0]
			}

			return packageCmd(chart, pkgLocation)
		},
	}

	cmdPackage.Flags().StringVarP(&pkgLocation, "directory", "d", "./", "directory in which to save the chart")

	rootCmd.AddCommand(cmdPackage)
}

func packageCmd(chart string, destination string) error {
	helmChart, err := loader.Load(chart)

	if err != nil {
		return fmt.Errorf("failed to load chart from %s: %w\n", chart, err)
	}

	path, err := chartutil.Save(helmChart, destination)

	if err != nil {
		return fmt.Errorf("failed to export chart to %s: %w\n", destination, err)
	}

	fmt.Printf("wrote package to %s\n", path)

	return nil
}
