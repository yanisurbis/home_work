#### Результатом выполнения следующих домашних заданий является сервис «Календарь»:
- [Домашнее задание №12 «Заготовка сервиса Календарь»](./docs/12_README.md)
- [Домашнее задание №13 «Внешние API от Календаря»](./docs/13_README.md)
- [Домашнее задание №14 «Кроликизация Календаря»](./docs/14_README.md)
- [Домашнее задание №15 «Докеризация и интеграционное тестирование Календаря»](./docs/15_README.md)

**Домашнее задание не принимается, если не принято ДЗ, предшедствующее ему.**

### Helpers

```
cd internal/grpc
protoc protobufs/events.proto --go_out=plugins=grpc:.
```

### To start

- start server, main.go
- start client, client.go

### Sending request

- add `userid = 1` to headers
- add `from = Date.now()/1000` to body, form-url-encoded 