package presentation

import "prc_hub_bench/domain/model/benchmark"

type IMarshaler interface {
	Marshal(r benchmark.Result) ([]byte, error)
}
