package benchmark

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/tingtt/prc_hub_bench/infrastructure/externalapi/backend"
)

func UsersIdStarPost(c *backend.Client, userId int64, wantedStatusCode int) (d time.Duration, err error) {
	start := time.Now()

	r, err := c.PostUsersIdStar(context.Background(), rand.Int63n(99)+1)

	d = time.Since(start)

	if err != nil {
		return
	}
	if r.StatusCode != wantedStatusCode {
		err = fmt.Errorf("failed to request (POST /users/:id/star): expected %d, found %d", wantedStatusCode, r.StatusCode)
		return
	}

	// TODO: chceck response body
	// b, err := io.ReadAll(r.Body)
	// if err != nil {
	// 	return
	// }

	return
}
