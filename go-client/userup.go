package userup

import (
	"context"
	"encoding/json"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/structpb"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/google/uuid"
	"github.com/hillside-labs/userservice-go-sdk/pkg/userapi"
)

// NewClient creates a new instance of the UserService client.
// It establishes a gRPC connection to the specified URI and returns the client.
// The URI should be in the format "host:port".
// If the connection cannot be established, an error is returned.
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

type UserID struct {
	ID         uint64
	UUID       uuid.UUID
	ExternalID string
}

func UID(id uint64) UserID {
	return UserID{ID: id}
}

func UUID(uuid uuid.UUID) UserID {
	return UserID{UUID: uuid}
}

func ExtID(id string) UserID {
	return UserID{ExternalID: id}
}

type User struct {
	ID         UserID
	Username   string
	Attributes map[string]interface{}
	Traits     map[string]interface{}
}

type UserSearchParams struct {
	ID               UserID
	Username         string
	AttributeFilters []*userapi.AttributeFilter
	TraitFilters     []*userapi.TraitFilter
}

func (usp UserSearchParams) WithAttribute(name string, value interface{}, operator ...userapi.Operator) UserSearchParams {
	op := userapi.Operator_EQUALS
	if len(operator) > 0 {
		op = operator[0]
	}

	pbValue, _ := structpb.NewValue(value)
	filter := userapi.AttributeFilter{
		Name:     name,
		Value:    pbValue,
		Operator: op,
	}
	usp.AttributeFilters = append(usp.AttributeFilters, &filter)
	return usp
}

func (usp UserSearchParams) WithTrait(name string, value interface{}, operator ...userapi.Operator) UserSearchParams {
	op := userapi.Operator_EQUALS
	if len(operator) > 0 {
		op = operator[0]
	}

	pbValue, _ := structpb.NewValue(value)
	filter := userapi.TraitFilter{
		Name:     name,
		Value:    pbValue,
		Operator: op,
	}
	usp.TraitFilters = append(usp.TraitFilters, &filter)
	return usp
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
	UserID          UserID
}

type UserService struct {
	addr   string
	client userapi.UsersClient
	close  func() error
}

func (us *UserService) Close() {
	us.close()
}

// AddUser adds a new user to the user service.
// It takes a context and a pointer to a User struct as input.
// It returns a pointer to the created User and an error, if any.
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
		ExternalId: user.ID.ExternalID,
		Uuid:       user.ID.UUID.String(),
		Username:   user.Username,
		Attributes: attrStruct,
		Traits:     traitStruct,
	})
	if err != nil {
		return nil, err
	}
	return UserResponseToUser(userResp), nil
}

func rpcUserID(id UserID) *userapi.UserID {
	return &userapi.UserID{
		Id:         id.ID,
		Uuid:       id.UUID.String(),
		ExternalId: id.ExternalID,
	}
}

func clientUserID(id *userapi.UserID) UserID {
	return UserID{
		ID:         id.Id,
		UUID:       uuid.MustParse(id.Uuid),
		ExternalID: id.ExternalId,
	}
}

// AddAttribute adds an attribute to a user with the specified ID.
// It takes a context, user ID, attribute key, and attribute value as parameters.
// The attribute value can be of any type and will be converted to a structpb.Value.
// Returns an error if there was a problem adding the attribute.
func (us UserService) AddAttribute(ctx context.Context, id UserID, key string, value interface{}) error {
	attrVal, err := structpb.NewValue(value)
	if err != nil {
		return err
	}
	_, err = us.client.AddAttribute(ctx, &userapi.AttributeRequest{
		UserId: rpcUserID(id),
		Key:    key,
		Value:  attrVal,
	})
	return err
}

// AddTrait adds a trait to a user identified by their ID.
// It takes a context, user ID, trait key, and trait value as parameters.
// The trait value can be of any type and will be converted to a structpb.Value.
// Returns an error if there was a problem adding the trait.
func (us UserService) AddTrait(ctx context.Context, id UserID, key string, value interface{}) error {
	traitVal, err := structpb.NewValue(value)
	if err != nil {
		return err
	}
	_, err = us.client.AddTrait(ctx, &userapi.TraitRequest{
		UserId: rpcUserID(id),
		Key:    key,
		Value:  traitVal,
	})
	return err
}

// UpdateUser updates a user with the provided user data.
// It takes a context.Context and a pointer to a User struct as input.
// It returns a pointer to the updated User struct and an error if any.
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
		Id:         rpcUserID(user.ID),
		Username:   user.Username,
		Attributes: attrStruct,
		Traits:     traitStruct,
	})
	if err != nil {
		return nil, err
	}
	return UserResponseToUser(userResp), nil
}

// GetUser retrieves a user by their ID.
// It makes a request to the user service API to fetch the user details.
// Returns the user object if found, otherwise returns an error.
func (us UserService) GetUser(ctx context.Context, id UserID) (*User, error) {
	userResp, err := us.client.Get(ctx, &userapi.UserRequest{
		Id: rpcUserID(id),
	})
	if err != nil {
		return nil, err
	}
	return UserResponseToUser(userResp), nil
}

func UserSearchToUserQuery(usp *UserSearchParams) *userapi.UserQuery {
	query := userapi.UserQuery{
		UserId:           rpcUserID(usp.ID),
		Username:         usp.Username,
		AttributeFilters: usp.AttributeFilters,
		TraitFilters:     usp.TraitFilters,
	}
	return &query
}

func (us UserService) FindUser(ctx context.Context, usp *UserSearchParams) ([]*User, error) {
	userQuery := UserSearchToUserQuery(usp)
	usersResp, err := us.client.Find(ctx, userQuery)
	if err != nil {
		return nil, err
	}
	users := make([]*User, 0, len(usersResp.Users))
	for _, u := range usersResp.Users {
		users = append(users, UserResponseToUser(u))
	}
	return users, nil
}

// DeleteUser deletes a user by their ID.
func (us UserService) DeleteUser(ctx context.Context, id UserID) error {
	_, err := us.client.Delete(ctx, &userapi.UserRequest{
		Id: rpcUserID(id),
	})
	return err
}

func (us UserService) DeleteAttribute(ctx context.Context, userId UserID, key string) error {
	_, err := us.client.DeleteAttribute(ctx, &userapi.AttributeRequest{
		UserId: rpcUserID(userId),
		Key:    key,
	})
	return err
}

func (us UserService) DeleteTrait(ctx context.Context, userId UserID, key string) error {
	_, err := us.client.DeleteTrait(ctx, &userapi.TraitRequest{
		UserId: rpcUserID(userId),
		Key:    key,
	})
	return err
}

func (us UserService) SearchUserTraits(ctx context.Context, userId UserID, names []string, begin time.Time, end time.Time) ([]interface{}, error) {

	traitsResp, err := us.client.SearchUserTraits(ctx, &userapi.SearchUserTraitsRequest{
		UserId: rpcUserID(userId),
		Names:  names,
		Begin:  timestamppb.New(begin),
		End:    timestamppb.New(end),
	})
	if err != nil {
		return nil, err
	}

	traits := make([]interface{}, 0, len(traitsResp.Traits))

	for _, v := range traitsResp.Traits {
		traitMap := v.AsMap()
		traits = append(traits, traitMap)
	}

	return traits, nil
}

func (us UserService) GetUsersByTraits(ctx context.Context, names []string, begin time.Time, end time.Time) ([]*User, error) {

	usersResp, err := us.client.GetUsersByTraits(ctx, &userapi.SearchUserTraitsRequest{
		Names: names,
		Begin: timestamppb.New(begin),
		End:   timestamppb.New(end),
	})
	if err != nil {
		return nil, err
	}

	users := make([]*User, 0, len(usersResp.Users))

	for _, u := range usersResp.Users {
		users = append(users, UserResponseToUser(u))
	}

	return users, nil
}

func (us UserService) GetUsersByEvents(ctx context.Context, types []string, sources []string, schemas []string, begin time.Time, end time.Time) ([]*User, error) {

	usersResp, err := us.client.GetUsersByEvents(ctx, &userapi.SearchUserEventsRequest{
		Types:   types,
		Sources: sources,
		Schemas: schemas,
		Begin:   timestamppb.New(begin),
		End:     timestamppb.New(end),
	})
	if err != nil {
		return nil, err
	}

	users := make([]*User, 0, len(usersResp.Users))

	for _, u := range usersResp.Users {
		users = append(users, UserResponseToUser(u))
	}

	return users, nil
}

func (us UserService) SearchEvents(ctx context.Context, userId UserID, types []string, begin time.Time, end time.Time) ([]Event, error) {
	eventsResp, err := us.client.SearchEvents(ctx, &userapi.SearchEventsRequest{
		UserId: rpcUserID(userId),
		Names:  types,
		Begin:  timestamppb.New(begin),
		End:    timestamppb.New(end),
	})
	if err != nil {
		return nil, err
	}
	events := make([]Event, len(eventsResp.Events))
	for i, event := range eventsResp.Events {
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
			UserID:          clientUserID(event.UserId),
		}
	}
	return events, nil
}

func UserResponseToUser(userResp *userapi.UserResponse) *User {
	attributes := make(map[string]interface{})
	for k, v := range userResp.Attributes.Fields {
		attributes[k] = v.AsInterface()
	}
	traits := make(map[string]interface{})
	for k, v := range userResp.Traits.Fields {
		traits[k] = v.AsInterface()
	}
	return &User{
		ID:         clientUserID(userResp.Id),
		Username:   userResp.Username,
		Attributes: attributes,
		Traits:     traits,
	}
}

// QueryUsers queries the user service with the given query and returns a list of users and an error, if any.
// The query parameter specifies the criteria for filtering the users.
// The returned list of users contains the user ID, username, UUID, attributes, and traits.
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
			ID:         clientUserID(user.Id),
			Username:   user.Username,
			Attributes: attributes,
			Traits:     traits,
		}
	}
	return users, nil
}

// QueryAttributes queries the attributes of a user based on the provided query.
// It takes a context.Context and a *Query as input parameters.
// It returns a map[string]interface{} containing the attributes of the user and an error if any.
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

// QueryTraits queries the traits of users based on the provided query.
// It takes a context.Context and a *Query as input parameters.
// It returns a map[string]interface{} containing the traits of the users and an error if any.
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

// QueryEvents queries events based on the provided query parameters.
// It returns a slice of Event objects and an error if any.
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
			UserID:          clientUserID(event.UserId),
		}
	}
	return events, nil
}

func (us UserService) AddSession(ctx context.Context, sessionKey string, sessionData map[string]interface{}) error {
	obj, err := structpb.NewStruct(sessionData)
	if err != nil {
		return err
	}

	_, err = us.client.AddSession(ctx, &userapi.SessionRequest{
		Session: &userapi.Session{
			Key:    sessionKey,
			Object: obj,
		},
	})

	return err
}

func (us UserService) IdentifySession(ctx context.Context, sessionKey string, userID UserID) error {
	_, err := us.client.IdentifySession(ctx, &userapi.IdentifySessionRequest{
		SessionKey: []string{sessionKey},
		UserId:     rpcUserID(userID),
	})

	return err
}

type SessionQuery struct {
	UserID  UserID
	Keys    []string
	Begin   time.Time
	End     time.Time
	Limit   int32
	Offset  int32
	OrderBy string
}

func (us UserService) GetSessions(ctx context.Context, query *SessionQuery) ([]*userapi.Session, error) {
	sessionsResp, err := us.client.GetSessions(ctx, &userapi.GetSessionsRequest{
		UserId:      rpcUserID(query.UserID),
		SessionKeys: query.Keys,
		Begin:       timestamppb.New(query.Begin),
		End:         timestamppb.New(query.End),
		Limit:       query.Limit,
		Offset:      query.Offset,
		OrderBy:     query.OrderBy,
	})

	if err != nil {
		return nil, err
	}

	return sessionsResp.Sessions, nil
}

type SessionEventQuery struct {
	UserID      UserID
	SessionKeys []string
	Begin       time.Time
	End         time.Time
	Limit       int32
	Offset      int32
	OrderBy     string
}

func (us UserService) GetSessionEvents(ctx context.Context, query *SessionEventQuery) ([]*userapi.Event, error) {
	sessionEventsResp, err := us.client.GetSessionEvents(ctx, &userapi.GetSessionEventsRequest{
		SessionKeys: query.SessionKeys,
		UserId:      rpcUserID(query.UserID),
		Begin:       timestamppb.New(query.Begin),
		End:         timestamppb.New(query.End),
		Limit:       query.Limit,
		Offset:      query.Offset,
		OrderBy:     query.OrderBy,
	})

	if err != nil {
		return nil, err
	}

	return sessionEventsResp.Events, nil
}
