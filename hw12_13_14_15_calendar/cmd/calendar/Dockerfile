FROM golang as builder
WORKDIR /calendar
COPY go.mod ./
RUN go mod download
COPY  ./ .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app cmd/calendar/main.go

FROM alpine
RUN apk --no-cache add ca-certificates
WORKDIR /root/go
COPY --from=builder /calendar/app .
CMD ["./app"]
EXPOSE 5432