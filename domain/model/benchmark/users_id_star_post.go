package benchmark

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"strconv"
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

	// Chceck response body
	b, err := io.ReadAll(r.Body)
	if err != nil {
		return
	}
	// Unmarshal
	count := StarCount{}
	err = json.Unmarshal(b, &count)
	if err != nil {
		return
	}

	return
}

func usersIdStarPost(c *backend.Client, userId int64, wantedStatusCode int) (
	count StarCount,
	err error,
) {
	r, err := c.PostUsersIdStar(context.Background(), rand.Int63n(99)+1)
	if err != nil {
		return
	}
	if r.StatusCode != wantedStatusCode {
		err = fmt.Errorf("failed to request (POST /users/:id/star): expected %d, found %d", wantedStatusCode, r.StatusCode)
		return
	}

	// log response
	b, err := io.ReadAll(r.Body)
	if err != nil {
		return
	}
	err = writeFile("./.log/users_id_start_POST_"+strconv.Itoa(r.StatusCode)+".json", b)
	if err != nil {
		return
	}

	// Unmarshal
	err = json.Unmarshal(b, &count)
	if err != nil {
		return
	}

	return
}
