package benchmark

import (
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"github.com/goombaio/namegenerator"
	"github.com/tingtt/prc_hub_bench/infrastructure/externalapi/backend"
)

func TestEndpoints(c *backend.Client) error {
	_, err := ResetPost(c)
	if err != nil {
		return err
	}

	mkdirIfNotExist("./.log/")

	LOGIN_USER := backend.LoginBody{Email: "throbbing-pond@prchub.com", Password: "throbbing-pond"}
	var TOKEN string
	TOKEN, err = usersSignInPost(c, LOGIN_USER, 200)
	if err != nil {
		fmt.Printf("Failed: POST /users/sign_in\n\terr: %v\n", err)
		return err
	}
	fmt.Println("Success: POST /users/sign_in")
	nameGenerator := namegenerator.NewNameGenerator(time.Now().UTC().UnixNano())

	events, err := eventsGet(c,
		backend.GetEventsParams{
			Embed: &[]string{"user", "documents"},
		},
		http.StatusOK,
	)
	if err != nil {
		fmt.Printf("Failed: GET /events\n\terr: %v\n", err)
		return err
	}
	fmt.Printf("Success: GET /events\n\tevents[0]: %+v\n", events[0])

	event, err := eventsIdGet(c,
		rand.Int63n(99)+1,
		backend.GetEventsIdParams{
			Embed: &[]string{"user", "documents"},
		},
		http.StatusOK,
	)
	if err != nil {
		fmt.Printf("Failed: GET /events/{id}\n\terr: %v\n", err)
		return err
	}
	fmt.Printf("Success: GET /events/{id}\n\tevent: %+v\n", event)

	name := nameGenerator.Generate()
	tmpBool := true
	event2, err := eventsPost(c,
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
	if err != nil {
		fmt.Printf("Failed: POST /events\n\terr: %v\n", err)
		return err
	}
	fmt.Printf("Success: POST /events\n\tevent: %+v\n", event2)

	users, err := usersGet(c, TOKEN, http.StatusOK)
	if err != nil {
		fmt.Printf("Failed: GET /users\n\terr: %v\n", err)
		return err
	}
	fmt.Printf("Success: GET /users\n\tusers[0]: %+v\n", users[0])

	count, err := usersIdStarPost(c, rand.Int63n(99)+1, http.StatusOK)
	if err != nil {
		fmt.Printf("Failed: POST /users/{id}/star\n\terr: %v\n", err)
		return err
	}
	fmt.Printf("Success: POST /users/{id}/star\n\tcount: %+v\n", count)

	return nil
}
