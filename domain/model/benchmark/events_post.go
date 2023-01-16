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

func EventsPost(c *backend.Client, bearer string, b backend.PostEventsJSONRequestBody, wantedStatusCode int) (d time.Duration, err error) {
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

	// Chceck response body
	b2, err := io.ReadAll(r.Body)
	if err != nil {
		return
	}
	// Unmarshal
	event := Event{}
	err = json.Unmarshal(b2, &event)
	if err != nil {
		return
	}
	err = event.validate()
	if err != nil {
		return
	}

	return
}

func eventsPost(c *backend.Client, bearer string, b backend.PostEventsJSONRequestBody, wantedStatusCode int) (event Event, err error) {
	r, err := c.PostEvents(context.Background(), b,
		func(ctx context.Context, req *http.Request) error {
			req.Header.Add("Authorization", "Bearer "+bearer)
			return nil
		},
	)
	if err != nil {
		return
	}

	// log response
	b2, err := io.ReadAll(r.Body)
	if err != nil {
		return
	}
	err = writeFile("./.log/events_POST_"+strconv.Itoa(r.StatusCode)+".json", b2)
	if err != nil {
		return
	}

	// Check status code
	if r.StatusCode != wantedStatusCode {
		err = fmt.Errorf("failed to request (POST /events): expected %d, found %d", wantedStatusCode, r.StatusCode)
		return
	}

	// Unmarshal
	err = json.Unmarshal(b2, &event)
	if err != nil {
		return
	}
	err = event.validate()
	if err != nil {
		return
	}

	return
}
