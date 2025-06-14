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

