package benchmark

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"time"

	"github.com/tingtt/prc_hub_bench/infrastructure/externalapi/backend"
)

func UsersSignInPost(c *backend.Client, b backend.LoginBody, wantedStatusCode int) (d time.Duration, token string, err error) {
	// Request
	start := time.Now()
	r, err := c.PostUsersSignIn(context.Background(), b)
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
