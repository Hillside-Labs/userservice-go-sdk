package main

import (
	"context"
	"fmt"
	"log"

	"github.com/hillside-labs/userservice-go-sdk/pkg/userapi"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/structpb"
)

type UserService struct {
	addr   string
	client userapi.UsersClient
	ctx    context.Context
	close  func() error
}

func NewUserService(addr string) (*UserService, error) {
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	return &UserService{
		addr:   addr,
		client: userapi.NewUsersClient(conn),
		ctx:    context.Background(),
		close:  conn.Close,
	}, nil
}

func (us *UserService) Close() {
	us.close()
}

func main() {
	us, err := NewUserService("localhost:9001")
	if err != nil {
		log.Fatal(err)
	}
	defer us.Close()

	attrs, _ := structpb.NewStruct(map[string]interface{}{
		"avid user": true,
	})
	
	user, err := us.client.Create(us.ctx, &userapi.NewUser{
		Username:   "devangana",
		Attributes: attrs,
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(user)

	user, err = us.client.Get(us.ctx, &userapi.UserRequest{Id: 3})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(user)
}
