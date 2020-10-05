ifndef $(GOPATH)
	GOPATH=$(shell go env GOPATH)
	export GOPATH
endif

build:
	go build -o calendar cmd/app/*

test:
	go test -race ./...

lint:
	go get github.com/golangci/golangci-lint/cmd/golangci-lint@v1.30.0
		go mod tidy
	$(GOPATH)/bin/golangci-lint run ./...

proto:
	protoc --proto_path=./internal/server/grpc/protobufs --go_out=./internal/server/grpc/events_grpc --plugin=$(GOPATH)/bin/protoc-gen-go --go_opt=paths=source_relative events.proto

.PHONY: build
