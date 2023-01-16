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

func EventsIdDocumentsIdGet(c *backend.Client, ctx context.Context, bearer string, eventId string, documentId string, wantedStatusCode int) (d time.Duration, err error) {
	start := time.Now()

	r, err := c.GetEventsIdDocumentsDocumentId(
		ctx,
		eventId,
		documentId,
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
		err = fmt.Errorf("failed to request (GET /events/:id/documents/:document_id): expected %d, found %d", wantedStatusCode, r.StatusCode)
		return
	}

	// Chceck response body
	b, err := io.ReadAll(r.Body)
	if err != nil {
		return
	}
	// Unmarshal
	document := EventDocument{}
	err = json.Unmarshal(b, &document)
	if err != nil {
		return
	}
	err = document.validate()
	if err != nil {
		return
	}

	return
}

func eventsIdDocumentsIdGet(c *backend.Client, bearer string, eventId string, documentId string, wantedStatusCode int) (document EventDocument, err error) {
	r, err := c.GetEventsIdDocumentsDocumentId(
		context.Background(),
		eventId,
		documentId,
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
	err = writeFile("./.log/events_id_GET_"+strconv.Itoa(r.StatusCode)+".json", b)
	if err != nil {
		return
	}

	// Check status code
	if r.StatusCode != wantedStatusCode {
		err = fmt.Errorf("failed to request (GET /events/:id/documents/:document_id): expected %d, found %d", wantedStatusCode, r.StatusCode)
		return
	}

	// Unmarshal
	err = json.Unmarshal(b, &document)
	if err != nil {
		return
	}
	err = document.validate()
	if err != nil {
		return
	}

	return
}
