FROM golang as builder
WORKDIR /tests
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY  . .
CMD go test -v ./integration_tests