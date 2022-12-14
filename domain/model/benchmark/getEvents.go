package benchmark

import (
	"context"
	"fmt"
	"time"

	"github.com/tingtt/prc_hub_bench/infrastructure/externalapi/backend"
)

func getEvents(c *backend.Client, p backend.GetEventsParams, wantedStatusCode int) (d time.Duration, err error) {
	start := time.Now()

	r, err := c.GetEvents(
		context.Background(),
		&p,
	)

	d = time.Since(start)

	if err != nil {
		return
	}
	if r.StatusCode != wantedStatusCode {
		err = fmt.Errorf("failed to request (GET /events): expected %d, found %d", wantedStatusCode, r.StatusCode)
		return
	}

	// TODO: chceck response body
	// b, err := io.ReadAll(r.Body)
	// if err != nil {
	// 	return
	// }

	return
}
