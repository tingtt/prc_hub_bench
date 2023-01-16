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

	var eventId string = string(events[rand.Int63n(int64(len(events)-1))].Id)
	event, err := eventsIdGet(c,
		eventId,
		backend.GetEventsIdParams{
			Embed: &[]string{"user", "documents"},
		},
		http.StatusOK,
	)
	if err != nil {
		fmt.Printf("Failed: GET /events/%s\n\terr: %v\n", eventId, err)
		return err
	}
	fmt.Printf("Success: GET /events/%s\n\tevent: %+v\n", eventId, event)

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

	userId := string(users[rand.Int63n(int64(len(users)-1))].Id)
	user, err := usersIdGet(c, userId, TOKEN, http.StatusOK)
	if err != nil {
		fmt.Printf("Failed: GET /users/%s\n\terr: %v\n", userId, err)
		return err
	}
	fmt.Printf("Success: GET /users/%s\n\tuser: %+v\n", userId, user)

	count, err := usersIdStarPost(c, userId, http.StatusOK)
	if err != nil {
		fmt.Printf("Failed: POST /users/%s/star\n\terr: %v\n", userId, err)
		return err
	}
	fmt.Printf("Success: POST /users/%s/star\n\tcount: %+v\n", userId, count)

	return nil
}
