package meta

type Annotations map[string]string

func CreateAnnotations() Annotations {
    return make(Annotations)
}
