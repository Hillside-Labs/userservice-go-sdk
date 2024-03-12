package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/fatih/color"
	userup "github.com/hillside-labs/userservice-go-sdk/go-client"
	"github.com/rodaine/table"
	"github.com/urfave/cli/v2"
)

func main() {

	app := cli.App{
		Name:  "sessmgr",
		Usage: "Inspect and work with Anonymous Sessions.",
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

					fmt.Println("session created: ", sessID)
					out, _ := json.MarshalIndent(doc, "", "  ")
					fmt.Println(string(out))

					return nil
				},
			},
			{
				Name:  "ls",
				Usage: "List our existing sessions.",
				Action: func(c *cli.Context) error {
					client, err := userup.NewClient("localhost:9000")
					if err != nil {
						log.Fatal(err)
					}
					defer client.Close()

					sessions, err := client.GetSessions(context.Background(), &userup.SessionQuery{
						Limit: 100,
					})
					if err != nil {
						log.Fatal(err)
					}

					if len(sessions) == 0 {
						return nil
					}

					headerFmt := color.New(color.FgGreen, color.Underline).SprintfFunc()
					columnFmt := color.New(color.FgHiBlue, color.Bold).SprintfFunc()

					tbl := table.New("Session Key", "User ID", "Object")
					tbl.WithHeaderFormatter(headerFmt).WithFirstColumnFormatter(columnFmt)

					for _, sess := range sessions {

						objJSON, _ := json.MarshalIndent(sess.Object.AsMap(), "", "  ")
						tbl.AddRow(sess.Key, sess.UserId, string(objJSON))
					}

					tbl.Print()

					return nil
				},
			},
			{
				Name:  "get",
				Usage: "Get a specific session's events.",
				Action: func(c *cli.Context) error {
					client, err := userup.NewClient("localhost:9000")
					if err != nil {
						log.Fatal(err)
					}
					defer client.Close()

					sessID := c.Args().First()
					if sessID == "" {
						return fmt.Errorf("Please provide a Session ID argument.")
					}

					events, err := client.GetSessionEvents(context.Background(), &userup.SessionEventQuery{
						SessionKeys: []string{sessID},
					})

					if err != nil {
						log.Fatal(err)
					}

					headerFmt := color.New(color.FgGreen, color.Underline).SprintfFunc()
					columnFmt := color.New(color.FgHiBlue, color.Bold).SprintfFunc()

					tbl := table.New("Subject", "Type", "User ID", "Object")
					tbl.WithHeaderFormatter(headerFmt).WithFirstColumnFormatter(columnFmt)

					for _, event := range events {
						tbl.AddRow(event.Subject, event.Type, event.UserId, string(event.Data))
					}

					tbl.Print()

					return nil
				},
			},
			{
				Name:  "evt",
				Usage: "Create an event within a session",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name: "type",
					},
					&cli.StringFlag{
						Name: "subject",
					},
					&cli.StringFlag{
						Name: "data",
					},
				},
				Action: func(c *cli.Context) error {
					client, err := userup.NewClient("localhost:9000")
					if err != nil {
						log.Fatal(err)
					}
					defer client.Close()
					sessID := c.Args().First()
					if sessID == "" {
						return fmt.Errorf("A Session Key is required to create a session event.")
					}

					if c.String("type") == "" {
						return fmt.Errorf("A Type is required. The type is a reverse dns name that describes the type of event.")
					}

					if c.String("subject") == "" {
						return fmt.Errorf("A Subject is required. The subject is what names the session event.")
					}

					if c.String("data") == "" {
						return fmt.Errorf("Data is required for the event.")
					}

					loggerConfig := userup.EventLoggerConfig{
						Source:      "https://userup.io/demo/sessmgr/elarson",
						SpecVersion: "1.0",
						UserService: client,
					}
					logger := userup.NewLogger(loggerConfig)

					fmt.Println("data: ", c.String("data"))
					var data map[string]interface{}
					err = json.Unmarshal([]byte(c.String("data")), &data)
					if err != nil {
						return err
					}

					logger.LogSessionEvent(
						context.Background(),
						sessID,
						c.String("type"),
						"userup.demo.schema",
						c.String("subject"),
						data,
					)

					return nil
				},
			},
			{
				Name:  "identify",
				Usage: "Identify a session as belonging to a specific user.",
				Action: func(c *cli.Context) error {
					client, err := userup.NewClient("localhost:9000")
					if err != nil {
						log.Fatal(err)
					}
					defer client.Close()

					sessID := c.Args().First()
					userIDstr := c.Args().Get(1)

					if sessID == "" || userIDstr == "" {
						return fmt.Errorf("Missing args: SESSION_ID USER_ID")
					}

					userID, err := strconv.ParseUint(userIDstr, 10, 64)
					if err != nil {
						return err
					}

					return client.IdentifySession(context.Background(), sessID, userID)
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
