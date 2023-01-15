package benchmark

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/tingtt/prc_hub_bench/infrastructure/externalapi/backend"
)

func UsersGet(c *backend.Client, bearer string, wantedStatusCode int) (d time.Duration, err error) {
	start := time.Now()

	r, err := c.GetUsers(
		context.Background(),
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
		err = fmt.Errorf("failed to request (GET /users): expected %d, found %d", wantedStatusCode, r.StatusCode)
		return
	}

	// TODO: chceck response body
	// b, err := io.ReadAll(r.Body)
	// if err != nil {
	// 	return
	// }

	return
}
