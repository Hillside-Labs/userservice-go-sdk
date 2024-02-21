package userup

import (
	"context"

	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/hillside-labs/userservice-go-sdk/pkg/userapi"
)

type EventLoggerConfig struct {
	Source      string
	SpecVersion string
}

type EventLogger struct {
	config EventLoggerConfig
}

func (e EventLogger) LogEvent(ctx context.Context, us UserService, userId uint64, dataType string, schema string, subject string, dataContentType string, data []byte) (*Event, error) {
	event := &userapi.Event{
		UserId:          userId,
		Source:          e.config.Source,
		Type:            dataType,
		Data:            data,
		Specversion:     e.config.SpecVersion,
		Timestamp:       timestamppb.Now(),
		Id:              uuid.New().String(),
		Datacontenttype: dataContentType,
		Subject:         subject,
		Dataschema:      schema,
	}
	eventResp, err := us.client.LogEvent(ctx, &userapi.EventRequest{
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
