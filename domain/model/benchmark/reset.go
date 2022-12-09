package benchmark

import (
	"context"
	"fmt"
	"prc_hub_bench/infrastructure/externalapi/backend"
	"time"
)

func reset(c *backend.Client) (d time.Duration, err error) {
	// Request
	start := time.Now()
	r, err := c.PostReset(context.Background())
	d = time.Since(start)

	// Check status code
	if err == nil {
		return
	}
	if r.StatusCode != 200 {
		err = fmt.Errorf("failed to request (POST /reset): expected 200, found %d", r.StatusCode)
		return
	}
	return
}
