package meta

type ChartAnnotations map[string]string

func CreateAnnotations() ChartAnnotations {
    return make(ChartAnnotations)
}
