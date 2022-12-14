package marshaler_json

import (
	"encoding/json"

	"github.com/tingtt/prc_hub_bench/domain/model/benchmark"
	"github.com/tingtt/prc_hub_bench/presentation"
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
