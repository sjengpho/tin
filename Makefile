.DEFAULT_GOAL := build
PROTO_OUTPUT := ./proto
GO_BIN := $(or $(GOBIN),$(HOME)/go/bin)

build:
	go build -o ./bin/tin-server cmd/grpc/main.go
	go build -o ./bin/tin cmd/cli/main.go

install:
	go build -o $(GO_BIN)/tin-server cmd/grpc/main.go
	go build -o $(GO_BIN)/tin cmd/cli/main.go

test:
	go test ./...

protoc:
	protoc -I proto ./proto/grpc.proto --go_out=plugins=grpc:$(PROTO_OUTPUT)
