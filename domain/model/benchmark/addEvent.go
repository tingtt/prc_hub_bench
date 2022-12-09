package benchmark

import (
	"context"
	"fmt"
	"net/http"
	"prc_hub_bench/infrastructure/externalapi/backend"
	"time"
)

func addEvent(c *backend.Client, bearer string, b backend.PostEventsJSONRequestBody, wantedStatusCode int) (d time.Duration, err error) {
	start := time.Now()

	r, err := c.PostEvents(context.Background(), b,
		func(ctx context.Context, req *http.Request) error {
			req.Header.Add("Authorization", "Bearer "+bearer)
			return nil
		},
	)

	d = time.Since(start)

	if err != nil {
		return
	}
	if r.StatusCode != wantedStatusCode {
		err = fmt.Errorf("failed to request (POST /events): expected %d, found %d", wantedStatusCode, r.StatusCode)
		return
	}

	// TODO: chceck response body
	// b2, err := io.ReadAll(r.Body)
	// if err != nil {
	// 	return
	// }

	return
}
