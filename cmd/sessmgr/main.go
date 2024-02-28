package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	userup "github.com/hillside-labs/userservice-go-sdk/go-client"
	"github.com/urfave/cli/v2"
)

func main() {

	app := cli.App{
		Name:  "sessmgr",
		Usage: "Inspect and work with Anonymous Sessions.",
		Action: func(c *cli.Context) error {
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

			client.AddUserToSession(context.Background(), sessID, uint64(123))
			fmt.Println("Added session", sessID, "to user 123")

			return nil
		},
		Commands: []*cli.Command{
			{
				Name:  "add",
				Usage: "Add a new session.",
				Flags: []cli.Flag{
					&cli.BoolFlag{
						Name:    "generate",
						Aliases: []string{"g"},
					},
					&cli.StringFlag{
						Name:    "json",
						Aliases: []string{"j"},
					},
				},
				Action: func(c *cli.Context) error {
					client, err := userup.NewClient("localhost:9000")
					if err != nil {
						log.Fatal(err)
					}
					defer client.Close()

					sessID := c.Args().First()
					if c.Bool("generate") {
						sessID = userup.NewSessionID()
					}

					doc := make(map[string]interface{})
					json.Unmarshal([]byte(c.String("json")), &doc)

					err = client.AddSession(context.Background(), sessID, doc)
					if err != nil {
						fmt.Println("Error creating a session")
						log.Fatal(err)
					}

					return nil
				},
			},
			{
				Name:  "ls",
				Usage: "List our existing sessions.",
			},
			{
				Name:  "get",
				Usage: "Get a specific session's events.",
			},

			{
				Name:  "identify",
				Usage: "Identify a session as belonging to a specific user.",
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
