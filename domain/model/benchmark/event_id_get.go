package benchmark

import (
	"context"
	"fmt"
	"time"

	"github.com/tingtt/prc_hub_bench/infrastructure/externalapi/backend"
)

func EventsIdGet(c *backend.Client, id int64, p backend.GetEventsIdParams, wantedStatusCode int) (d time.Duration, err error) {
	start := time.Now()

	r, err := c.GetEventsId(
		context.Background(),
		id,
		&p,
	)

	d = time.Since(start)

	if err != nil {
		return
	}
	if r.StatusCode != wantedStatusCode {
		err = fmt.Errorf("failed to request (GET /events/:id): expected %d, found %d", wantedStatusCode, r.StatusCode)
		return
	}

	// TODO: chceck response body
	// b, err := io.ReadAll(r.Body)
	// if err != nil {
	// 	return
	// }

	return
}
