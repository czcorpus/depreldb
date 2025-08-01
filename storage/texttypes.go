package storage

type PreconfTextTypeMapping struct {
	data map[string]byte
}

func (mapping *PreconfTextTypeMapping) RawToReadable(val byte) string {
	for k, v := range mapping.data {
		if v == val {
			return k
		}
	}
	return ""
}

func (mapping *PreconfTextTypeMapping) ReadableToRaw(val string) byte {
	return mapping.data[val]
}

func NewPreconfTextTypeMapping() *PreconfTextTypeMapping {
	return &PreconfTextTypeMapping{
		data: make(map[string]byte),
	}
}
