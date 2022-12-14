package presentation

import "github.com/tingtt/prc_hub_bench/domain/model/benchmark"

type IMarshaler interface {
	Marshal(r benchmark.Result) ([]byte, error)
}
