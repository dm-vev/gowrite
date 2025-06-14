# gowrite

`gowrite` предоставляет упрощенный клиент для [Appwrite](https://appwrite.io). Он реализует работу с базой данных и хранилищем, а также содержит утилиты для построения запросов.

## Установка

```
go get github.com/dm-vev/gowrite
```

## Использование

```go
package main

import (
    "github.com/dm-vev/gowrite"
)

func main() {
    client := gowrite.NewClient("https://localhost/v1", "<PROJECT_ID>", "<API_KEY>")
    databases := gowrite.NewDatabases(client)
    storage := gowrite.NewStorage(client)
    _ = databases
    _ = storage
}
```


## CI/CD

Для запуска интеграционных тестов в GitHub Actions добавьте следующие секреты репозитория:

- `APPWRITE_ENDPOINT` — URL сервера Appwrite
- `APPWRITE_PROJECT_ID` — идентификатор проекта
- `APPWRITE_API_KEY` — API‑ключ с правами на работу с базой данных

Секреты можно добавить в разделе **Settings → Secrets and variables → Actions**. После этого workflow `.github/workflows/test.yml` автоматически выполнит `go test -v ./...`.

