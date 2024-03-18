package main

import (
	"context"
	"fmt"
	"log"
	"os"

	userup "github.com/hillside-labs/userservice-go-sdk/go-client"
)

func main() {
	client, err := userup.NewClient("localhost:9000")
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	logger := userup.NewLogger(userup.NewLoggerConfig("https://userup.io/sample-client", client))

	user := &userup.User{
		Username: "jdoe2",
		Attributes: map[string]interface{}{
			"user_type": "admin",
			"email":     "jdoe@localhost.com",
			"ranking":   5,
		},
	}

	ctx := context.Background()
	userRet, err := client.AddUser(ctx, user)
	if err != nil {
		log.Fatal(err)
	}

	logger.LogEvent(context.Background(),
		userup.Event{
			Type:    "io.userup.user.created",
			Subject: "user",
			Data:    user,
		},
	)

	fmt.Println("User ID:", userRet.ID)
	user, err = client.GetUser(ctx, userRet.ID)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("+%v\n", user)

	client.AddAttribute(ctx, userRet.ID, "alias", "dumbledore")
	query := userup.Query{
		Joins: []userup.Join{
			{
				Table: "attributes",
				On:    "users.id = attributes.user_id",
				Filter: map[string]userup.Condition{
					"attribute": {
						"alias": "dumbledore",
					},
				},
			},
		},
	}
	users, err := client.QueryUsers(ctx, &query)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("query results")
	for _, u := range users {
		fmt.Printf("+%v\n", u)
	}

	os.Exit(0)
}
