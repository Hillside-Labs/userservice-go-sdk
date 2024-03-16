SRC=$(find . -name '*.go')

userservice-go-sdk: $(SRC)
	go build ./go-client

sessmgr: $(SRC)
	go build ./cmd/sessmgr/

sample-client:
	go build ./cmd/sample-client

protos:
	protoc --go_out=./pkg/userapi --go_opt=paths=source_relative --go-grpc_out=./pkg/userapi --go-grpc_opt=paths=source_relative --proto_path=../userservice-proto ../userservice-proto/user.proto
