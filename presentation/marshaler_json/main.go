package marshaler_json

import (
	"encoding/json"
	"prc_hub_bench/domain/model/benchmark"
	"prc_hub_bench/presentation"
)

func New(indent string) presentation.IMarshaler {
	return &m{indent: indent}
}

type m struct {
	indent string
}

func (m *m) Marshal(r benchmark.Result) ([]byte, error) {
	return json.MarshalIndent(r, "", m.indent)
}
