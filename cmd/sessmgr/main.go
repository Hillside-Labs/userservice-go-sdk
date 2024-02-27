package main

import (
	"context"
	"fmt"
	"log"

	userup "github.com/hillside-labs/userservice-go-sdk/go-client"
	"github.com/hillside-labs/userservice-go-sdk/pkg/userapi"
)

func main() {
	client, err := userup.NewClient("localhost:9000")
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	fmt.Println("Created a connection")

	sessID := userup.NewSessionID()
	err = client.AddSession(context.Background(), sessID, map[string]interface{}{"hello": "world"})
	if err != nil {
		fmt.Println("Error creating a session")
		log.Fatal(err)
	}

	fmt.Println("Created a new session with session id ", sessID)

	sl := userup.NewSessionLogger(userup.SessionEventLoggerConfig(userup.NewSessionLoggerConfig("io.userup.test.sessmgr", client)))

	fmt.Println("New session logger created")

	err = sl.LogEvent(context.Background(), sessID, "application/json", "test_schema", "view_cart", map[string]interface{}{"url": "/greatjeans/2/blue"})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Sent a session event")

	client.IdentifySession(context.Background(), sessID, uint64(123))
	fmt.Println("Added session", sessID, "to user 123")

	sessions, err := client.GetSessions(context.Background(), &userup.SessionQuery{
		UserId: &userapi.UserID{Id: uint64(123)},
	},)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Got sessions for user 123")
	for _, session := range sessions {
		fmt.Printf("%+v\n", session)
	}
}
