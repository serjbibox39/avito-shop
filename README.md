# Avito Merch Store

Сервис для управления магазином мерча в Авито.

## Запуск проекта

1. Установите зависимости:
   ```bash
   go mod tidy

   migrate -path migrations -database "postgres://user:password@localhost:5432/avito_merch?sslmode=disable" up

   go run cmd/server/main.go