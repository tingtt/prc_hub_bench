package marshaler_yaml

import (
	"prc_hub_bench/domain/model/benchmark"
	"prc_hub_bench/presentation"

	"github.com/go-yaml/yaml"
)

func New() presentation.IMarshaler {
	return &m{}
}

type m struct{}

func (*m) Marshal(r benchmark.Result) ([]byte, error) {
	return yaml.Marshal(r)
}
