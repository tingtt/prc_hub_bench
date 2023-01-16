package benchmark

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"strconv"
	"time"

	"github.com/tingtt/prc_hub_bench/infrastructure/externalapi/backend"
)

func UsersSignInPost(c *backend.Client, ctx context.Context, b backend.LoginBody, wantedStatusCode int) (d time.Duration, token string, err error) {
	// Request
	start := time.Now()
	r, err := c.PostUsersSignIn(ctx, b)
	d = time.Since(start)

	// Check status code
	if err != nil {
		return
	}
	if r.StatusCode != wantedStatusCode {
		err = fmt.Errorf("failed to request (POST /users/sign_in): expected %d, found %d", wantedStatusCode, r.StatusCode)
		return
	}

	// Read body
	b2, err := io.ReadAll(r.Body)
	if err != nil {
		return
	}
	t := &backend.Token{}
	err = json.Unmarshal(b2, t)
	if err != nil {
		return
	}
	return d, t.Token, nil
}

func usersSignInPost(c *backend.Client, b backend.LoginBody, wantedStatusCode int) (token string, err error) {
	// Request
	r, err := c.PostUsersSignIn(context.Background(), b)
	// Check status code
	if err != nil {
		return
	}

	// Read body
	b2, err := io.ReadAll(r.Body)
	if err != nil {
		return
	}
	t := &backend.Token{}
	err = json.Unmarshal(b2, t)
	if err != nil {
		return
	}

	// log response
	err = writeFile("./.log/users_sign_in_POST_"+strconv.Itoa(r.StatusCode)+".json", b2)
	if err != nil {
		return
	}

	// Check status code
	if r.StatusCode != wantedStatusCode {
		err = fmt.Errorf("failed to request (POST /users/sign_in): expected %d, found %d", wantedStatusCode, r.StatusCode)
		return
	}

	return t.Token, nil
}
