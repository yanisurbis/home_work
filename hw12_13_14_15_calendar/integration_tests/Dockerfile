FROM golang as builder
WORKDIR /tests
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY  . .
#RUN cd integration_tests/
CMD go test -v ./integration_tests

#FROM alpine
#RUN apk --no-cache add ca-certificates
#WORKDIR /root/go
#COPY --from=builder /calendar/app .
#COPY ./configs ./configs
#CMD ["./app"]

# docker build -t calendar-api  -f cmd/calendar/Dockerfile .