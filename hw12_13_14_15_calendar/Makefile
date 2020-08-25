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
	$(GOPATH)/bin/golangci-lint run ./...

.PHONY: build
