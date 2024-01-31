userservice-go-sdk:
	go build ./cmd/userservice

protos:
	protoc --go_out=./cmd/userservice --go_opt=paths=source_relative --go-grpc_out=./cmd/userservice --go-grpc_opt=paths=source_relative --proto_path=../userservice-proto ../userservice-proto/user.proto