package benchmark

import (
	"time"

	"github.com/tingtt/prc_hub_bench/domain/model/benchmark"
	"github.com/tingtt/prc_hub_bench/infrastructure/externalapi/backend"
)

type Result = benchmark.Result

func Run(c *backend.Client, d time.Duration) Result {
	return benchmark.Run(c, d)
}
