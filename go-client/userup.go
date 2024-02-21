package userup

import (
	"context"
	"encoding/json"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/structpb"

	"github.com/hillside-labs/userservice-go-sdk/pkg/userapi"
)

func NewClient(uri string) (*UserService, error) {
	conn, err := grpc.Dial(uri, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	client := userapi.NewUsersClient(conn)

	return &UserService{
		addr:   uri,
		client: client,
		close:  conn.Close,
	}, nil
}

type User struct {
	Id         uint64
	Username   string
	Uuid       string
	Attributes map[string]interface{}
	Traits     map[string]interface{}
}

type Event struct {
	Timestamp       time.Time
	ID              string
	Source          string
	SpecVersion     string
	Type            string
	DataContentType string
	DataSchema      string
	Subject         string
	Data            []byte
	UserID          uint64
}

type UserService struct {
	addr   string
	client userapi.UsersClient
	close  func() error
}

func (us *UserService) Close() {
	us.close()
}

func (us UserService) AddUser(ctx context.Context, user *User) (*User, error) {

	attrStruct, err := structpb.NewStruct(user.Attributes)
	if err != nil {
		return nil, err
	}
	traitStruct, err := structpb.NewStruct(user.Traits)
	if err != nil {
		return nil, err
	}
	userResp, err := us.client.Create(ctx, &userapi.NewUser{
		Username:   user.Username,
		Uuid:       user.Uuid,
		Attributes: attrStruct,
		Traits:     traitStruct,
	})
	if err != nil {
		return nil, err
	}
	attributes := make(map[string]interface{})
	for k, v := range userResp.Attributes.Fields {
		attributes[k] = v.AsInterface()
	}
	traits := make(map[string]interface{})
	for k, v := range userResp.Traits.Fields {
		traits[k] = v.AsInterface()
	}
	return &User{
		Id:         userResp.Id,
		Username:   userResp.Username,
		Uuid:       userResp.Uuid,
		Attributes: attributes,
		Traits:     traits,
	}, nil
}

func (us UserService) AddAttribute(ctx context.Context, id uint64, key string, value interface{}) error {
	attrVal, err := structpb.NewValue(value)
	if err != nil {
		return err
	}
	_, err = us.client.AddAttribute(ctx, &userapi.AttributeRequest{
		UserId: id,
		Key:    key,
		Value:  attrVal,
	})
	return err
}

func (us UserService) AddTrait(ctx context.Context, id uint64, key string, value interface{}) error {
	traitVal, err := structpb.NewValue(value)
	if err != nil {
		return err
	}
	_, err = us.client.AddTrait(ctx, &userapi.TraitRequest{
		UserId: id,
		Key:    key,
		Value:  traitVal,
	})
	return err
}

func (us UserService) UpdateUser(ctx context.Context, user *User) (*User, error) {
	attrStruct, err := structpb.NewStruct(user.Attributes)
	if err != nil {
		return nil, err
	}
	traitStruct, err := structpb.NewStruct(user.Traits)
	if err != nil {
		return nil, err
	}
	userResp, err := us.client.Update(ctx, &userapi.UserRequest{
		Id:         user.Id,
		Username:   user.Username,
		Uuid:       user.Uuid,
		Attributes: attrStruct,
		Traits:     traitStruct,
	})
	if err != nil {
		return nil, err
	}
	attributes := make(map[string]interface{})
	for k, v := range userResp.Attributes.Fields {
		attributes[k] = v.AsInterface()
	}
	traits := make(map[string]interface{})
	for k, v := range userResp.Traits.Fields {
		traits[k] = v.AsInterface()
	}
	return &User{
		Id:         userResp.Id,
		Username:   userResp.Username,
		Uuid:       userResp.Uuid,
		Attributes: attributes,
		Traits:     traits,
	}, nil
}

func (us UserService) GetUser(ctx context.Context, id uint64) (*User, error) {
	userResp, err := us.client.Get(ctx, &userapi.UserRequest{
		Id: id,
	})
	if err != nil {
		return nil, err
	}
	attributes := make(map[string]interface{})
	for k, v := range userResp.Attributes.Fields {
		attributes[k] = v.AsInterface()
	}
	traits := make(map[string]interface{})
	for k, v := range userResp.Traits.Fields {
		traits[k] = v.AsInterface()
	}
	return &User{
		Id:         userResp.Id,
		Username:   userResp.Username,
		Uuid:       userResp.Uuid,
		Attributes: attributes,
		Traits:     traits,
	}, nil
}

func (us UserService) DeleteUser(ctx context.Context, id uint64) error {
	_, err := us.client.Delete(ctx, &userapi.UserRequest{
		Id: id,
	})
	return err
}

func (us UserService) QueryUsers(ctx context.Context, query *Query) ([]*User, error) {
	queryJson, err := json.Marshal(query)
	if err != nil {
		return nil, err
	}
	userResp, err := us.client.QueryUsers(ctx, &userapi.QueryRequest{Query: queryJson})
	if err != nil {
		return nil, err
	}
	users := make([]*User, len(userResp.Users))
	for i, user := range userResp.Users {
		attributes := make(map[string]interface{})
		for k, v := range user.Attributes.Fields {
			attributes[k] = v.AsInterface()
		}
		traits := make(map[string]interface{})
		for k, v := range user.Traits.Fields {
			traits[k] = v.AsInterface()
		}
		users[i] = &User{
			Id:         user.Id,
			Username:   user.Username,
			Uuid:       user.Uuid,
			Attributes: attributes,
			Traits:     traits,
		}
	}
	return users, nil
}

func (us UserService) QueryAttributes(ctx context.Context, query *Query) (map[string]interface{}, error) {
	queryJson, err := json.Marshal(query)
	if err != nil {
		return nil, err
	}
	attrResp, err := us.client.QueryAttributes(ctx, &userapi.QueryRequest{Query: queryJson})
	if err != nil {
		return nil, err
	}
	return attrResp.Attributes.AsMap(), nil
}

func (us UserService) QueryTraits(ctx context.Context, query *Query) (map[string]interface{}, error) {
	queryJson, err := json.Marshal(query)
	if err != nil {
		return nil, err
	}
	traitResp, err := us.client.QueryTraits(ctx, &userapi.QueryRequest{Query: queryJson})
	if err != nil {
		return nil, err
	}
	return traitResp.Traits.AsMap(), nil
}

func (us UserService) QueryEvents(ctx context.Context, query *Query) ([]Event, error) {
	queryJson, err := json.Marshal(query)
	if err != nil {
		return nil, err
	}
	eventResp, err := us.client.QueryEvents(ctx, &userapi.QueryRequest{Query: queryJson})
	if err != nil {
		return nil, err
	}
	events := make([]Event, len(eventResp.Events))
	for i, event := range eventResp.Events {
		events[i] = Event{
			Timestamp:       event.Timestamp.AsTime(),
			ID:              event.Id,
			Source:          event.Source,
			SpecVersion:     event.Specversion,
			Type:            event.Type,
			DataContentType: event.Datacontenttype,
			DataSchema:      event.Dataschema,
			Subject:         event.Subject,
			Data:            event.Data,
			UserID:          event.UserId,
		}
	}
	return events, nil
}
