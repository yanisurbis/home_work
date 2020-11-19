ifndef $(GOPATH)
	GOPATH=$(shell go env GOPATH)
	export GOPATH
endif

build:
	go build -o ./bin/calendar ./cmd/calendar

test:
	go test -race ./internal/...

goimports:
	GO111MODULE=off go get golang.org/x/tools/cmd/goimports
	$(GOPATH)/bin/goimports -w ./internal

gci:
	GO111MODULE=off go get github.com/daixiang0/gci
	$(GOPATH)/bin/gci -w .

install-lint-deps:
	(which golangci-lint > /dev/null) || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v1.30.0

lint: install-lint-deps
	golangci-lint run ./...

#lint:
#	go get github.com/golangci/golangci-lint/cmd/golangci-lint@v1.30.0
#		go mod tidy
#	$(GOPATH)/bin/golangci-lint run ./...

resetdb:
	GO111MODULE=off go get -u github.com/pressly/goose/cmd/goose
	$(GOPATH)/bin/goose -dir migrations postgres "user=yanis password=yanis dbname=events sslmode=disable" down
	$(GOPATH)/bin/goose -dir migrations postgres "user=yanis password=yanis dbname=events sslmode=disable" up

proto:
	#protoc protobufs/events.proto --go_out=plugins=grpc:.
	protoc --proto_path=./internal/server/grpc/protobufs --go_out=./internal/server/grpc/events_grpc --plugin=$(GOPATH)/bin/protoc-gen-go --go_opt=paths=source_relative events.proto

build:
	docker-compose -f docker-compose.yml build

up:
	docker-compose -f docker-compose.yml up

down:
	docker-compose -f docker-compose.yml down

integration-tests:
	# there should be a better way
	docker rm calendar_test_db || true
	docker rm calendar_test_queue || true
	docker-compose -f docker-compose.yml -f docker-compose.test.yml -p tests build
	docker-compose -f docker-compose.yml -f docker-compose.test.yml -p tests up --abort-on-container-exit --exit-code-from sut

.PHONY: build test lint
