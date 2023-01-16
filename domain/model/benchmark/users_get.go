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

	// Chceck response body
	b, err := io.ReadAll(r.Body)
	if err != nil {
		return
	}
	users := []User{}
	err = json.Unmarshal(b, &users)
	if err != nil {
		return
	}
	err = users[0].validate()
	if err != nil {
		return
	}

	return
}

func usersGet(c *backend.Client, bearer string, wantedStatusCode int) (users []User, err error) {
	r, err := c.GetUsers(
		context.Background(),
		func(ctx context.Context, req *http.Request) error {
			req.Header.Add("Authorization", "Bearer "+bearer)
			return nil
		},
	)
	if err != nil {
		return
	}
	if r.StatusCode != wantedStatusCode {
		err = fmt.Errorf("failed to request (GET /users): expected %d, found %d", wantedStatusCode, r.StatusCode)
		return
	}

	// log response
	b, err := io.ReadAll(r.Body)
	if err != nil {
		return
	}
	err = writeFile("./.log/users_GET_"+strconv.Itoa(r.StatusCode)+".json", b)
	if err != nil {
		return
	}

	// Unmarshal
	err = json.Unmarshal(b, &users)
	if err != nil {
		return
	}
	err = users[0].validate()
	if err != nil {
		return
	}

	return
}
