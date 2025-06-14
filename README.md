# gowrite

`gowrite` предоставляет упрощенный клиент для [Appwrite](https://appwrite.io). Он реализует работу с базой данных и хранилищем, а также содержит утилиты для построения запросов.

## Установка

```
go get github.com/example/gowrite
```

## Использование

```go
package main

import (
    "github.com/example/gowrite"
)

func main() {
    client := appwrite.NewClient("https://localhost/v1", "<PROJECT_ID>", "<API_KEY>")
    databases := appwrite.NewDatabases(client)
    storage := appwrite.NewStorage(client)
    _ = databases
    _ = storage
}
```

Пакет содержит под‑пакеты:

- `id` – генерация уникальных идентификаторов;
- `query` – конструктор запросов для поиска документов.

Полную документацию по API Appwrite можно найти на официальном сайте проекта.
