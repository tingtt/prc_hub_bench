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

func EventsIdDocumentsGet(c *backend.Client, ctx context.Context, bearer string, id string, p backend.GetEventsIdDocumentsParams, wantedStatusCode int) (documents []EventDocument, d time.Duration, err error) {
	start := time.Now()

	r, err := c.GetEventsIdDocuments(
		ctx,
		id,
		&p,
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
		err = fmt.Errorf("failed to request (GET /events/:id/documents): expected %d, found %d", wantedStatusCode, r.StatusCode)
		return
	}

	// Chceck response body
	b, err := io.ReadAll(r.Body)
	if err != nil {
		return
	}
	// Unmarshal
	err = json.Unmarshal(b, &documents)
	if err != nil {
		return
	}
	for _, ed := range documents {
		err = ed.validate()
		if err != nil {
			return
		}
	}

	return
}

func eventsIdDocumentsGet(c *backend.Client, bearer string, id string, p backend.GetEventsIdDocumentsParams, wantedStatusCode int) (documents []EventDocument, err error) {
	r, err := c.GetEventsIdDocuments(
		context.Background(),
		id,
		&p,
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
	err = writeFile("./.log/events_id_documents_GET_"+strconv.Itoa(r.StatusCode)+".json", b)
	if err != nil {
		return
	}

	// Check status code
	if r.StatusCode != wantedStatusCode {
		err = fmt.Errorf("failed to request (GET /events/:id/documents): expected %d, found %d", wantedStatusCode, r.StatusCode)
		return
	}

	// Unmarshal
	err = json.Unmarshal(b, &documents)
	if err != nil {
		return
	}
	for _, ed := range documents {
		err = ed.validate()
		if err != nil {
			return
		}
	}

	return
}
