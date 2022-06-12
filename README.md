
В сервисе:
- [X]  Подключение и подписка на канал в nats-streaming
- [X]  Полученные данные писать в Postgres
- [X]  Так же полученные данные сохранить in memory в сервисе (Кеш)
- [X]  В случае падения сервиса восстанавливать Кеш из Postgres
- [X]  Поднять http сервер и выдавать данные по id из кеша
- [X]  Сделать простейший интерфейс отображения полученных данных, для
их запроса по id

Чтобы поднять сервис, сначала нужно запустить контейнеры:
```bash
docker-compose up -d
```

Я решил не использовать сервис миграции для такой просто базы, поэтому 
ее нужно создать вручную:

```SQL
CREATE TABLE orders
(
    id BIGSERIAL PRIMARY KEY,
    order_num VARCHAR(255) NOT NULL UNIQUE,
    order_data jsonb NOT NULL
);
```

После чего следует запустить producer/main.go и  consumer/main.go
```bash
go run cmd/producer/main.go
go run cmd/consumer/main.go
```
