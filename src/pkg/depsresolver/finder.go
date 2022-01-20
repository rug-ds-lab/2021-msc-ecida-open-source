package depsresolver

import (
	"ecida/pkg/meta"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"helm.sh/helm/v3/pkg/chart"
	"helm.sh/helm/v3/pkg/chart/loader"
)

func isPathDirectory(path string) bool {
	fileinfo, err := os.Stat(path)

	if err != nil {
		return false
	}

	return fileinfo.IsDir()
}

func concatPaths(base string, suffix string) string {
	if isPathDirectory(base) {
		return filepath.Join(base, suffix)
	} else {
		return filepath.Join(filepath.Base(base), suffix)
	}
}

func readChartFromPath(path string) (*chart.Chart, error) {
	return loader.Load(path)
}

func makeConnectString(chart *chart.Chart) (string, error) {
	// inspect the service block in the values to read the name and port
	svc, svcexists := chart.Values["service"]

	if !svcexists {
		return "", errors.New("service block is not defined in values.yaml")
	}

	service, ismap := svc.(map[string]interface{})

	if !ismap {
		return "", errors.New("service is not a block in values.yaml")
	}

	name, namevalid := service["name"].(string)
	port, portvalid := service["port"].(string)

	if !namevalid || !portvalid {
		return "", errors.New("service.name and service.port must be strings in values.yaml")
	}

	return fmt.Sprintf("%s:%s", name, port), nil
}

func FindDependencies(root string) ([]*chart.Chart, error) {

	rootChart, err := readChartFromPath(root)

	if err != nil {
		return nil, fmt.Errorf("couldn't load chart %s: %w\n", root, err)
	}

	metadata := meta.EcidaMetaFromChart(rootChart)

	deps := []*chart.Chart{rootChart}

	for depName, pkg := range metadata {
		newPath := concatPaths(root, pkg)

		childDeps, err := FindDependencies(newPath)

		for _, childDep := range childDeps {
			connectStr, err := makeConnectString(childDep)
			if err != nil {
				return nil, fmt.Errorf("failed to resolve dependency for %s: %w\n", newPath, err)
			}
			rootChart.Values[depName] = connectStr
		}

		if err != nil {
			return nil, fmt.Errorf("couldn't load chart: %s: %w\n", newPath, err)
		}

		deps = append(deps, childDeps...)
	}

	return deps, nil
}
