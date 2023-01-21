package benchmark

import (
	"encoding/json"
	"errors"
	"strconv"
	"time"
)

type Id string

func (id *Id) UnmarshalJSON(b []byte) error {
	var num int64
	if err1 := json.Unmarshal(b, &num); err1 == nil {
		*id = Id(strconv.FormatInt(num, 10))
		return nil
	}
	var str string
	if err1 := json.Unmarshal(b, &str); err1 == nil {
		*id = Id(str)
		return nil
	} else {
		return err1
	}
}

type Event struct {
	Id          Id              `json:"id"`
	Name        string          `json:"name"`
	Description *string         `json:"description,omitempty"`
	Location    *string         `json:"location,omitempty"`
	Datetimes   []EventDatetime `json:"datetimes"`
	Published   bool            `json:"published"`
	Completed   bool            `json:"completed"`
	UserId      Id              `json:"user_id"`
}

type EventDatetime struct {
	Start time.Time `json:"start"`
	End   time.Time `json:"end"`
}

type EventDocument struct {
	EventId Id     `json:"event_id"`
	Id      Id     `json:"id"`
	Name    string `json:"name"`
	Url     string `json:"url"`
}

type EventEmbed struct {
	Event
	User      User            `json:"user"`
	Documents []EventDocument `json:"documents"`
}

type User struct {
	Id                  Id      `json:"id"`
	Name                string  `json:"name"`
	PostEventAvailabled bool    `json:"post_event_availabled"`
	Manage              bool    `json:"manage"`
	Admin               bool    `json:"admin"`
	TwitterId           *string `json:"twitter_id,omitempty"`
	GithubUsername      *string `json:"github_username,omitempty"`
	StarCount           uint64  `json:"star_count"`
}

type StarCount struct {
	Count uint `json:"count"`
}

func (e Event) validate() error {
	if e.Id == "" {
		return errors.New("event.id require")
	}
	if e.Name == "" {
		return errors.New("event.name require")
	}
	for _, ed := range e.Datetimes {
		if err := ed.validate(); err != nil {
			return errors.New("event.[]" + err.Error())
		}
	}
	if e.UserId == "" {
		return errors.New("event.user_id require")
	}

	return nil
}

func (e EventDatetime) validate() error {
	if e.Start.IsZero() {
		return errors.New("datetimes.start require")
	}
	if e.End.IsZero() {
		return errors.New("datetimes.end require")
	}

	return nil
}

func (e EventDocument) validate() error {
	if e.EventId == "" {
		return errors.New("document.event_id require")
	}
	if e.Id == "" {
		return errors.New("document.id require")
	}
	if e.Name == "" {
		return errors.New("document.name require")
	}
	if e.Url == "" {
		return errors.New("document.url require")
	}

	return nil
}

func (e EventEmbed) validate() error {
	if err := e.Event.validate(); err != nil {
		return errors.New(err.Error())
	}
	if err := e.User.validate(); err != nil {
		return errors.New("event." + err.Error())
	}
	for _, ed := range e.Documents {
		if err := ed.validate(); err != nil {
			return errors.New("event.[]" + err.Error())
		}
	}

	return nil
}

func (u User) validate() error {
	if u.Id == "" {
		return errors.New("user.id require")
	}
	if u.Name == "" {
		return errors.New("user.id require")
	}
	// TODO: validate star_count

	return nil
}
