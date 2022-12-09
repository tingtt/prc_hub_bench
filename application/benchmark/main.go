package benchmark

import (
	"prc_hub_bench/domain/model/benchmark"
	"prc_hub_bench/infrastructure/externalapi/backend"
	"time"
)

type Result = benchmark.Result

func Run(c *backend.Client, d time.Duration) Result {
	return benchmark.Run(c, d)
}
