ifndef $(GOPATH)
	GOPATH=$(shell go env GOPATH)
	export GOPATH
endif

build:
	go build -o calendar cmd/app/*

test:
	go test -race ./...

lint:
	GO111MODULE=off go get github.com/golangci/golangci-lint/cmd/golangci-lint@v1.30.0
	$(GOPATH)/bin/golangci-lint run ./...

resetdb:
	GO111MODULE=off go get -u github.com/pressly/goose/cmd/goose
	$(GOPATH)/bin/goose -dir migrations postgres "user=yanis password=yanis dbname=events sslmode=disable" down
	$(GOPATH)/bin/goose -dir migrations postgres "user=yanis password=yanis dbname=events sslmode=disable" up

proto:
	protoc ./internal/grpc/protobuf/events.proto --go_out=plugins=grpc:.

.PHONY: build
