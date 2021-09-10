package depsresolver

import (
	"ecida/pkg/meta"
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

func makeConnectString(chart *chart.Chart) string {
    // inspect the service block in the values to read the name and port
    svc, svcexists := chart.Values["service"]

    if !svcexists {
        return ""
    }

    // TODO: type checking and proper parsing
    service := svc.(map[string]interface{})

    name := service["name"].(string)
    port := service["port"].(string)
    return fmt.Sprintf("%s:%s", name, port)
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
            rootChart.Values[depName] = makeConnectString(childDep)
        }

        if err != nil {
            return nil, fmt.Errorf("couldn't load chart: %s: %w\n", newPath, err)
        }

        deps = append(deps, childDeps...)
    }
    
    return deps, nil
}
