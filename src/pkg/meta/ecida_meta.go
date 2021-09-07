package meta

import (
	"fmt"
)

type EcidaMetadata map[string]string

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
