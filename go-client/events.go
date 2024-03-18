// Package userup provides functionality for logging events in the user service.
package userup

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/hillside-labs/userservice-go-sdk/pkg/userapi"
)

// EventLoggerConfig represents the configuration for the EventLogger.
type EventLoggerConfig struct {
	Source      string       // Source represents the source of the events.
	SpecVersion string       // SpecVersion represents the version of the event specification.
	UserService *UserService // UserService represents the user service client.
}

// EventLogger represents a logger for logging events in the user service.
type EventLogger struct {
	config EventLoggerConfig // config represents the configuration for the EventLogger.
}

// NewLoggerConfig creates a new EventLoggerConfig with the specified source and UserService.
// It returns the created EventLoggerConfig.
func NewLoggerConfig(source string, userService *UserService) EventLoggerConfig {
	return EventLoggerConfig{
		Source:      source,
		SpecVersion: "1.0",
		UserService: userService,
	}
}

// NewLogger creates a new EventLogger with the specified configuration.
// It returns the created EventLogger.
func NewLogger(config EventLoggerConfig) EventLogger {
	return EventLogger{
		config: config,
	}
}

// LogEvent logs an event in the user service.
// It takes the user ID, data type, schema, subject, and data as input parameters.
// It returns the logged event and an error if any.
func (e EventLogger) LogEvent(ctx context.Context, event Event) (*Event, error) {

	// check for required fields

	if event.Source == "" {
		event.Source = e.config.Source
	}

	if event.Type == "" {
		return nil, fmt.Errorf("`Type` is required for the event")
	}

	if event.ID == "" {
		event.ID = uuid.New().String()
	}

	if event.DataContentType == "" {
		event.DataContentType = "application/json"
	}

	if event.SpecVersion == "" {
		event.SpecVersion = e.config.SpecVersion
	}

	if event.Timestamp.IsZero() {
		event.Timestamp = time.Now()
	}

	jsonData, err := json.Marshal(event.Data)
	if err != nil {
		return nil, err
	}

	apiEvent := &userapi.Event{
		UserId:          rpcUserID(event.UserID),
		Source:          event.Source,
		Type:            event.Type,
		Data:            jsonData,
		Specversion:     event.SpecVersion,
		Timestamp:       timestamppb.New(event.Timestamp),
		Id:              event.ID,
		Datacontenttype: event.DataContentType,
		Subject:         event.Subject,
		Dataschema:      event.DataSchema,
	}
	eventResp, err := e.config.UserService.client.LogEvent(ctx, &userapi.EventRequest{
		Event: apiEvent,
	})
	if err != nil {
		return nil, err
	}
	return &Event{
		Timestamp:       eventResp.Event.Timestamp.AsTime(),
		ID:              eventResp.Event.Id,
		Source:          eventResp.Event.Source,
		SpecVersion:     eventResp.Event.Specversion,
		Type:            eventResp.Event.Type,
		DataContentType: eventResp.Event.Datacontenttype,
		DataSchema:      eventResp.Event.Dataschema,
		Subject:         eventResp.Event.Subject,
		Data:            eventResp.Event.Data,
		SessionKey:      eventResp.Event.SessionKey,
		UserID:          clientUserID(eventResp.Event.UserId),
	}, nil
}
