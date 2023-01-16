package benchmark

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/tingtt/prc_hub_bench/infrastructure/externalapi/backend"
)

func UsersIdGet(c *backend.Client, ctx context.Context, id string, bearer string, wantedStatusCode int) (d time.Duration, err error) {
	start := time.Now()

	r, err := c.GetUsersId(
		ctx,
		id,
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

	// Chceck response body
	b, err := io.ReadAll(r.Body)
	if err != nil {
		return
	}
	user := User{}
	err = json.Unmarshal(b, &user)
	if err != nil {
		return
	}
	err = user.validate()
	if err != nil {
		return
	}

	return
}

func usersIdGet(c *backend.Client, id string, bearer string, wantedStatusCode int) (user User, err error) {
	r, err := c.GetUsersId(
		context.Background(),
		id,
		func(ctx context.Context, req *http.Request) error {
			req.Header.Add("Authorization", "Bearer "+bearer)
			return nil
		},
	)
	if err != nil {
		return
	}

	// log response
	b, err := io.ReadAll(r.Body)
	if err != nil {
		return
	}
	err = writeFile("./.log/users_id_GET_"+strconv.Itoa(r.StatusCode)+".json", b)
	if err != nil {
		return
	}

	// Check status code
	if r.StatusCode != wantedStatusCode {
		err = fmt.Errorf("failed to request (GET /users/:id): expected %d, found %d", wantedStatusCode, r.StatusCode)
		return
	}

	// Unmarshal
	err = json.Unmarshal(b, &user)
	if err != nil {
		return
	}
	err = user.validate()
	if err != nil {
		return
	}

	return
}
