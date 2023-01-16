package benchmark

import (
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"github.com/tingtt/prc_hub_bench/infrastructure/externalapi/backend"

	"github.com/goombaio/namegenerator"
)

type Result struct {
	Score uint64 `json:"score" yaml:"score"`
	Error string `json:"error,omitempty" yaml:"error,omitempty"`
	Logs  []Log  `json:"logs,omitempty" yaml:"logs,omitempty"`
}

type Log struct {
	Target   string `json:"target" yaml:"target"`
	Duration string `json:"duration" yaml:"duration"`
}

func Run(c *backend.Client, d time.Duration, o struct{ Verbose bool }) (r Result) {
	d2, err := ResetPost(c)
	if err != nil {
		r.Error = err.Error()
		return
	}
	r.Logs = append(r.Logs, Log{
		"POST /reset (スコアに関係しない)",
		fmt.Sprintf("%d ms", d2.Abs().Milliseconds())},
	)

	LOGIN_USER := backend.LoginBody{Email: "throbbing-pond@prchub.com", Password: "throbbing-pond"}
	var TOKEN string
	d2, TOKEN, err = UsersSignInPost(c, LOGIN_USER, 200)
	if err != nil {
		r.Error = err.Error()
		return
	}
	r.Logs = append(r.Logs, Log{
		"POST /users/sign_in",
		fmt.Sprintf("%d ms", d2.Abs().Milliseconds())},
	)
	nameGenerator := namegenerator.NewNameGenerator(time.Now().UTC().UnixNano())

	var USER_ID, EVENT_ID string
	events, d2, err := EventsGet(c, backend.GetEventsParams{}, http.StatusOK)
	if err != nil {
		r.Error = err.Error()
		return
	}
	EVENT_ID = string(events[rand.Int63n(int64(len(events)-1))].Id)
	r.Logs = append(r.Logs, Log{
		"GET /events",
		fmt.Sprintf("%d ms", d2.Abs().Milliseconds())},
	)
	users, d2, err := UsersGet(c, TOKEN, http.StatusOK)
	if err != nil {
		r.Error = err.Error()
		return
	}
	USER_ID = string(users[rand.Int63n(int64(len(users)-1))].Id)
	r.Logs = append(r.Logs, Log{
		"GET /events",
		fmt.Sprintf("%d ms", d2.Abs().Milliseconds())},
	)

	multiRequests(&r, d, o,
		req{
			Name: "POST /users/sign_in",
			Req: func() (d time.Duration, err error) {
				d, TOKEN, err = UsersSignInPost(c, LOGIN_USER, 200)
				return d, err
			},
			Point: 2,
		},
		req{
			Name: "GET /events?embed=uesr&embed=documents",
			Req: func() (time.Duration, error) {
				events, d, err := EventsGet(c,
					backend.GetEventsParams{
						Embed: &[]string{"user", "documents"},
					},
					http.StatusOK,
				)
				if err == nil {
					EVENT_ID = string(events[rand.Int63n(int64(len(events)-1))].Id)
				}
				return d, err
			},
			Point: 5,
		},
		req{
			Name: "GET /events/:id?embed=uesr&embed=documents",
			Req: func() (time.Duration, error) {
				return EventsIdGet(c,
					EVENT_ID,
					backend.GetEventsIdParams{
						Embed: &[]string{"user", "documents"},
					},
					http.StatusOK,
				)
			},
			Point: 5,
		},
		req{
			Name: "POST /events",
			Req: func() (time.Duration, error) {
				name := nameGenerator.Generate()
				tmpBool := true
				return EventsPost(c,
					TOKEN,
					backend.CreateEventBody{
						Datetimes: &[]backend.CreateEventDatetime{{
							End:   time.Now().Add(time.Hour * 2).Format("2006-01-02T15:04:05Z07:00"),
							Start: time.Now().Format("2006-01-02T15:04:05Z07:00"),
						}, {
							End:   time.Now().Add(time.Hour * 26).Format("2006-01-02T15:04:05Z07:00"),
							Start: time.Now().Add(time.Hour * 24).Format("2006-01-02T15:04:05Z07:00"),
						}, {
							End:   time.Now().Add(time.Hour * 50).Format("2006-01-02T15:04:05Z07:00"),
							Start: time.Now().Add(time.Hour * 48).Format("2006-01-02T15:04:05Z07:00"),
						}},
						Description: &name,
						Location:    &name,
						Name:        name,
						Published:   &tmpBool,
					},
					http.StatusCreated,
				)
			},
			Point: 5,
		},
		req{
			Name: "GET /users",
			Req: func() (time.Duration, error) {
				users, d, err := UsersGet(c, TOKEN, http.StatusOK)
				if err == nil {
					USER_ID = string(users[rand.Int63n(int64(len(users)-1))].Id)
				}
				return d, err
			},
			Point: 5,
		},
		req{
			Name: "GET /users/:id",
			Req: func() (time.Duration, error) {
				return UsersIdGet(c, USER_ID, TOKEN, http.StatusOK)
			},
			Point: 3,
		},
		req{
			Name: "POST /users/:id/star",
			Req: func() (time.Duration, error) {
				return UsersIdStarPost(c, USER_ID, http.StatusOK)
			},
			Point: 3,
		},
	)

	return
}

type req struct {
	Name  string
	Req   func() (time.Duration, error)
	Point uint
}

func request(r *Result, d *time.Duration, req req, o struct{ Verbose bool }) {
	// Execute request and get duration
	d2, err := req.Req()
	if err != nil {
		r.Error = err.Error()
		return
	}

	// Sub from time left
	*d -= d2
	if *d <= 0 {
		return
	}

	// Add log
	r.Logs = append(r.Logs, Log{
		req.Name,
		fmt.Sprintf("%d ms", d2.Abs().Milliseconds()),
	})
	if o.Verbose {
		fmt.Printf("%v\n", Log{req.Name, fmt.Sprintf("%d ms", d2.Abs().Milliseconds())})
	}
	// Add score point
	r.Score += uint64(req.Point)
}

func loopRequest(r *Result, d *time.Duration, req req, o struct{ Verbose bool }) {
	for i := 0; true; i++ {
		request(r, d, req, o)
		if i%20 == 0 {
			req.Point += 1
		}
		time.Sleep(time.Second)

		if r.Error != "" {
			break
		}
	}
}

func multiRequests(r *Result, d time.Duration, o struct{ Verbose bool }, rr ...req) {
	for _, req := range rr {
		go loopRequest(r, &d, req, o)
	}
	for {
		if r.Error != "" || d <= 0 {
			// stop
			return
		}
		// continue
	}
}
