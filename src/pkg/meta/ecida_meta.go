package meta

import (
	"fmt"
	"strings"

	"helm.sh/helm/v3/pkg/chart"
)

type EcidaMetadata map[string]string

func EcidaMetaFromChart(chart *chart.Chart) EcidaMetadata {

	meta := NewEcidaMetadata()

	for key, value := range chart.Metadata.Annotations {
		if strings.HasPrefix(key, "ecida.") {
			meta[strings.TrimPrefix(key, "ecida.")] = value
		}
	}

	return meta
}

func NewEcidaMetadata() EcidaMetadata {
	return make(map[string]string)
}

func (emeta *EcidaMetadata) ToChartAnnotations() ChartAnnotations {
	chartMeta := make(ChartAnnotations)

	for key, value := range *emeta {
		chartMeta[fmt.Sprintf("ecida.%s", key)] = value
	}

	return chartMeta
}
