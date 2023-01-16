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

func EventsGet(c *backend.Client, p backend.GetEventsParams, wantedStatusCode int) (events []EventEmbed, d time.Duration, err error) {
	start := time.Now()

	r, err := c.GetEvents(
		context.Background(),
		&p,
	)

	d = time.Since(start)

	if err != nil {
		return
	}
	if r.StatusCode != wantedStatusCode {
		err = fmt.Errorf("failed to request (GET /events): expected %d, found %d", wantedStatusCode, r.StatusCode)
		return
	}

	// Chceck response body
	b, err := io.ReadAll(r.Body)
	if err != nil {
		return
	}
	// Unmarshal
	err = json.Unmarshal(b, &events)
	if err != nil {
		return
	}
	err = events[0].Event.validate()
	if err != nil {
		return
	}
	if p.Embed != nil {
		for _, v := range *p.Embed {
			if v == "user" {
				err = events[0].validate()
				if err != nil {
					return
				}
			}
			if v == "documents" {
				for _, ed := range events[0].Documents {
					err = ed.validate()
					if err != nil {
						return
					}
				}
			}
		}
	}

	return
}

func eventsGet(c *backend.Client, p backend.GetEventsParams, wantedStatusCode int) (events []EventEmbed, err error) {
	r, err := c.GetEvents(
		context.Background(),
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
	err = writeFile("./.log/events_GET_"+strconv.Itoa(r.StatusCode)+".json", b)
	if err != nil {
		return
	}

	// Check status code
	if r.StatusCode != wantedStatusCode {
		err = fmt.Errorf("failed to request (GET /events): expected %d, found %d", wantedStatusCode, r.StatusCode)
		return
	}

	// Unmarshal
	err = json.Unmarshal(b, &events)
	if err != nil {
		return
	}
	err = events[0].validate()
	if err != nil {
		return
	}

	return
}
