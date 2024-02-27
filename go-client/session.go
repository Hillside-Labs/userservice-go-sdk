package userup

import (
	"context"
	"crypto/rand"
	"encoding/json"
	"fmt"

	"github.com/hillside-labs/userservice-go-sdk/pkg/userapi"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// NewSessionID creates a new random ID for a session ID.
func NewSessionID() string {
	// We're not handling the error b/c it is not something we can
	// deal with. It also isn't clear what the failure modes are.
	p, _ := rand.Prime(rand.Reader, 64)
	return fmt.Sprintf("%x", p)
}

// SessionEventLoggerConfig represents the configuration for the SessionEventLogger.
type SessionEventLoggerConfig struct {
	Source      string       // Source represents the source of the events.
	SpecVersion string       // SpecVersion represents the version of the event specification.
	UserService *UserService // UserService represents the user service client.
}

// SessionEventLogger represents a logger for logging session events in the user service.
type SessionEventLogger struct {
	config SessionEventLoggerConfig // config represents the configuration for the EventLogger.
}

// NewSessionLoggerConfig creates a new SessionEventLoggerConfig with the specified sessionID, source and UserService.
// It returns the created EventLoggerConfig.
func NewSessionLoggerConfig(source string, userService *UserService) EventLoggerConfig {
	return EventLoggerConfig{
		Source:      source,
		SpecVersion: "1.0",
		UserService: userService,
	}
}

// NewLogger creates a new EventLogger with the specified configuration.
// It returns the created EventLogger.
func NewSessionLogger(config SessionEventLoggerConfig) SessionEventLogger {
	return SessionEventLogger{
		config: config,
	}
}

// LogEvent logs an event in the user service.
// It takes the user ID, data type, schema, subject, and data as input parameters.
// It returns the logged event and an error if any.
func (e SessionEventLogger) LogEvent(ctx context.Context, sessionKey string, dataType string, schema string, subject string, data interface{}) error {

	dataBytes, err := json.Marshal(data)
	if err != nil {
		return err
	}

	event := &userapi.SessionEvent{
		SessionKey:      sessionKey,
		Source:          e.config.Source,
		Type:            dataType,
		Data:            dataBytes,
		Specversion:     e.config.SpecVersion,
		Timestamp:       timestamppb.Now(),
		Datacontenttype: "application/json",
		Subject:         subject,
		Dataschema:      schema,
	}
	_, err = e.config.UserService.client.LogSessionEvent(ctx, &userapi.SessionEventRequest{
		Event: event,
	})
	return err
}
