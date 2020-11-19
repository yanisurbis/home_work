FROM golang as builder
WORKDIR /scheduler
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY  . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app cmd/scheduler/main.go

FROM alpine
RUN apk --no-cache add ca-certificates
WORKDIR /root/go
COPY --from=builder /scheduler/app .
COPY ./configs ./configs
CMD ["./app"]
