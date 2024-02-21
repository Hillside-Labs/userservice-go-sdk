// Package userup provides functionality for logging events in the user service.
package userup

import (
	"context"
	"encoding/json"

	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/hillside-labs/userservice-go-sdk/pkg/userapi"
)

// EventLoggerConfig represents the configuration for the EventLogger.
type EventLoggerConfig struct {
	Source      string       // Source represents the source of the events.
	SpecVersion string       // SpecVersion represents the version of the event specification.
	UserService UserService  // UserService represents the user service client.
}

// EventLogger represents a logger for logging events in the user service.
type EventLogger struct {
	config EventLoggerConfig  // config represents the configuration for the EventLogger.
}

// NewLoggerConfig creates a new EventLoggerConfig with the specified source and UserService.
// It returns the created EventLoggerConfig.
func NewLoggerConfig(source string, userService UserService) EventLoggerConfig {
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
func (e EventLogger) LogEvent(ctx context.Context, userId uint64, dataType string, schema string, subject string, data interface{}) (*Event, error) {

	dataBytes, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	event := &userapi.Event{
		UserId:          userId,
		Source:          e.config.Source,
		Type:            dataType,
		Data:            dataBytes,
		Specversion:     e.config.SpecVersion,
		Timestamp:       timestamppb.Now(),
		Id:              uuid.New().String(),
		Datacontenttype: "application/json",
		Subject:         subject,
		Dataschema:      schema,
	}
	eventResp, err := e.config.UserService.client.LogEvent(ctx, &userapi.EventRequest{
		Event: event,
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
		UserID:          eventResp.Event.UserId,
	}, nil
}
