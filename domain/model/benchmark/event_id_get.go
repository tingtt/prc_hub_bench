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

func EventsIdGet(c *backend.Client, ctx context.Context, id string, p backend.GetEventsIdParams, wantedStatusCode int) (d time.Duration, err error) {
	start := time.Now()

	r, err := c.GetEventsId(
		ctx,
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

	// Chceck response body
	b, err := io.ReadAll(r.Body)
	if err != nil {
		return
	}
	// Unmarshal
	event := Event{}
	err = json.Unmarshal(b, &event)
	if err != nil {
		return
	}
	err = event.validate()
	if err != nil {
		return
	}

	return
}

func eventsIdGet(c *backend.Client, id string, p backend.GetEventsIdParams, wantedStatusCode int) (event EventEmbed, err error) {
	r, err := c.GetEventsId(
		context.Background(),
		id,
		&p,
	)
	if err != nil {
		return
	}

	// log response
	b, err := io.ReadAll(r.Body)
	if err != nil {
		return
	}
	err = writeFile("./.log/events_id_GET_"+strconv.Itoa(r.StatusCode)+".json", b)
	if err != nil {
		return
	}

	// Check status code
	if r.StatusCode != wantedStatusCode {
		err = fmt.Errorf("failed to request (GET /events/:id): expected %d, found %d", wantedStatusCode, r.StatusCode)
		return
	}

	// Unmarshal
	err = json.Unmarshal(b, &event)
	if err != nil {
		return
	}
	err = event.validate()
	if err != nil {
		return
	}

	return
}
