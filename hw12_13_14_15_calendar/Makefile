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
	protoc ./internal/grpc/protobuf/events.proto --go_out=plugins=grpc:.

.PHONY: build
