package userup

import (
	"context"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/structpb"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/hillside-labs/userservice-go-sdk/pkg/userapi"
)

type Integration struct {
	ID         int
	Name       string
	Schedule   string
	ExecPath   string
	ConfigPath string
	Enabled    bool
	Settings   map[string]interface{}
}

func (i *Integration) toProto() *userapi.Integration {
	settings, _ := structpb.NewStruct(i.Settings)
	return &userapi.Integration{
		ID:         int32(i.ID),
		Name:       i.Name,
		Schedule:   i.Schedule,
		ExecPath:   i.ExecPath,
		ConfigPath: i.ConfigPath,
		Enabled:    i.Enabled,
		Settings:   settings,
	}
}

func integrationFromProto(i *userapi.Integration) *Integration {
	return &Integration{
		ID:         int(i.ID),
		Name:       i.Name,
		Schedule:   i.Schedule,
		ExecPath:   i.ExecPath,
		ConfigPath: i.ConfigPath,
		Enabled:    i.Enabled,
		Settings:   i.Settings.AsMap(),
	}
}

type Job struct {
	IntegrationName string
	Started         time.Time
	Ended           time.Time
	Status          userapi.JobStatus
	Error           string
	ID              int
}

func (j *Job) toProto() *userapi.Job {
	return &userapi.Job{
		IntegrationName: j.IntegrationName,
		Started:         timestamppb.New(j.Started),
		Ended:           timestamppb.New(j.Ended),
		Status:          j.Status,
		Error:           j.Error,
		ID:              int32(j.ID),
	}
}

func jobFromProto(j *userapi.Job) *Job {
	return &Job{
		IntegrationName: j.IntegrationName,
		Started:         j.Started.AsTime(),
		Ended:           j.Ended.AsTime(),
		Status:          j.Status,
		Error:           j.Error,
		ID:              int(j.ID),
	}
}

type IntegrationsService struct {
	addr   string
	client userapi.IntegrationsClient
	close  func() error
}

func (is *IntegrationsService) Close() {
	is.close()
}

// NewIntegrationsClient creates a new instance of the UserService Integrations client.
// It establishes a gRPC connection to the specified URI and returns the client.
// The URI should be in the format "host:port".
// If the connection cannot be established, an error is returned.
func NewIntegrationsClient(uri string) (*IntegrationsService, error) {
	conn, err := grpc.Dial(uri, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	client := userapi.NewIntegrationsClient(conn)

	return &IntegrationsService{
		addr:   uri,
		client: client,
		close:  conn.Close,
	}, nil
}

// AddIntegration adds a new integration to the user service.
func (is *IntegrationsService) AddIntegration(ctx context.Context, integration *Integration) (*Integration, error) {
	req := &userapi.IntegrationAddRequest{
		Integration: integration.toProto(),
	}
	resp, err := is.client.AddIntegration(ctx, req)
	if err != nil {
		return nil, err
	}
	return integrationFromProto(resp.Integration), nil
}

func (is *IntegrationsService) GetIntegration(ctx context.Context, name string) (*Integration, error) {
	req := &userapi.IntegrationGetRequest{
		Name: name,
	}
	resp, err := is.client.GetIntegration(ctx, req)
	if err != nil {
		return nil, err
	}
	return integrationFromProto(resp.Integration), nil
}

func (is *IntegrationsService) UpdateIntegration(ctx context.Context, integration *Integration) (*Integration, error) {
	req := &userapi.IntegrationUpdateRequest{
		Integration: integration.toProto(),
	}
	resp, err := is.client.UpdateIntegration(ctx, req)
	if err != nil {
		return nil, err
	}
	return integrationFromProto(resp.Integration), nil
}

func (is *IntegrationsService) RemoveIntegration(ctx context.Context, name string) error {
	req := &userapi.IntegrationRemoveRequest{
		Name: name,
	}
	_, err := is.client.RemoveIntegration(ctx, req)
	if err != nil {
		return err
	}
	return nil
}

func (is *IntegrationsService) ListIntegrations(ctx context.Context) ([]*Integration, error) {
	req := &userapi.IntegrationListRequest{}
	resp, err := is.client.ListIntegrations(ctx, req)
	if err != nil {
		return nil, err
	}
	integrations := make([]*Integration, 0)
	for _, integration := range resp.Integrations {
		apiIntegration := integrationFromProto(integration)
		integrations = append(integrations, apiIntegration)
	}
	return integrations, nil
}

func (is *IntegrationsService) JobUpdate(ctx context.Context, job *Job) error {
	req := &userapi.JobUpdateRequest{
		Job: job.toProto(),
	}
	_, err := is.client.JobUpdate(ctx, req)
	if err != nil {
		return err
	}
	return nil
}

func (is *IntegrationsService) GetJobHistory(ctx context.Context, integrationName string) ([]*Job, error) {
	req := &userapi.JobGetHistoryRequest{
		IntegrationName: integrationName,
	}
	resp, err := is.client.GetJobHistory(ctx, req)
	if err != nil {
		return nil, err
	}
	jobs := make([]*Job, 0)
	for _, job := range resp.JobHistory {
		apiJob := jobFromProto(job)
		jobs = append(jobs, apiJob)
	}
	return jobs, nil
}
