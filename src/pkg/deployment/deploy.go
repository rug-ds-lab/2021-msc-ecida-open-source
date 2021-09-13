package deployment

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"helm.sh/helm/v3/pkg/chart"
	"helm.sh/helm/v3/pkg/chartutil"
)

const DeployTempDir = "_ecida_deploy"

func Deploy(pipeline string, charts []*chart.Chart) error {

	// we deploy everything in reverse so that the tail of the pipeline is
	// deployed first, and the head last.
	for i, j := 0, len(charts)-1; i < j; i, j = i+1, j-1 {
		charts[i], charts[j] = charts[j], charts[i]
	}

	// For now we deploy with the helm CLI which we assume to be installed. The
	// deployment process is therefore: write the chart to a temporary
	// location, and then deploy the chart from that location.
	path, err := ensureTempLocation()

	if err != nil {
		return err
	}

	for _, chart := range charts {
		err = deploySingle(pipeline, chart, path)

		if err != nil {
			return err
		}
	}

	return nil
}

func deploySingle(pipeline string, chart *chart.Chart, path string) error {

	// update the name of the package so that it's scoped to this pipeline
	chartName := fmt.Sprintf("%s-%s", pipeline, chart.Metadata.Name)

	chart.Metadata.Name = chartName

	// create a temporary archive
	path, err := chartutil.Save(chart, path)

	if err != nil {
		return err
	}

	err = runHelmDeploy(chartName, path)

	if err != nil {
		return fmt.Errorf("failed to apply helm configuration: %w", err)
	}

	// remove the temporary archive
	err = os.Remove(path)

	if err != nil {
		return err
	}

	return nil
}

func runHelmDeploy(releaseName string, path string) error {
	return exec.Command("helm", "upgrade", "--install", releaseName, path).Run()
}

func ensureTempLocation() (string, error) {
	dir := os.TempDir()
	tempDir := filepath.Join(dir, DeployTempDir)

	err := os.MkdirAll(tempDir, os.ModePerm)

	if err != nil {
		return "", fmt.Errorf("failed to create tempdir: %w", err)
	}

	return tempDir, nil
}
