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
## Примеры всех функций

### Клиент
```go
client := appwrite.NewClient("https://localhost/v1", "PROJECT_ID", "API_KEY")
_ , err := client.SendRequest("GET", "/health", nil)
```

### Работа с базой данных
```go
db := appwrite.NewDatabases(client)
_, _ = db.ListDatabases()
_, _ = db.CreateDatabase("db1", "Main", true)
_, _ = db.GetDatabase("db1")
_, _ = db.UpdateDatabase("db1", "Main", true)
_ = db.DeleteDatabase("db1")
_, _ = db.ListCollections("db1")
_, _ = db.CreateCollection("db1", "col1", "Tasks", nil, false, true)
_, _ = db.GetCollection("db1", "col1")
_, _ = db.UpdateCollection("db1", "col1", "Tasks", nil, false, true)
_ = db.DeleteCollection("db1", "col1")
_, _ = db.CreateDocument("db1", "col1", "doc1", map[string]interface{}{"a":1}, nil)
_, _ = db.GetDocument("db1", "col1", "doc1")
_, _ = db.UpdateDocument("db1", "col1", "doc1", map[string]interface{}{"a":2}, nil)
_ = db.DeleteDocument("db1", "col1", "doc1")
_, _ = db.ListDocuments("db1", "col1", nil)
_, _ = db.CountDocuments("db1", "col1", nil)
```

### Работа с хранилищем
```go
s := appwrite.NewStorage(client)
_, _ = s.ListBuckets()
_, _ = s.CreateBucket("b1", "Files", nil, false, true, 0, nil, "none", false, false)
_, _ = s.GetBucket("b1")
_, _ = s.UpdateBucket("b1", "Files", nil, false, true, 0, nil, "none", false, false)
_ = s.DeleteBucket("b1")
_, _ = s.ListFiles("b1")
// Для CreateFile необходим путь к файлу
_, _ = s.CreateFile("b1", "file1", "/path/to/file", nil)
_, _ = s.GetFile("b1", "file1")
_, _ = s.UpdateFile("b1", "file1", "new", nil)
_ = s.DeleteFile("b1", "file1")
_, _ = s.DownloadFile("b1", "file1")
_, _ = s.GetFilePreview("b1", "file1", nil)
_, _ = s.ViewFile("b1", "file1")
url := s.GetFileDownloadURL("b1", "file1")
_ = url
```

### Генерация идентификаторов
```go
id.Custom("fixed")
id.Unique()
```

### Конструктор запросов
```go
query.Equal("name", "Bob")
query.NotEqual("status", "done")
query.LessThan("age", 18)
query.LessThanEqual("age", 18)
query.GreaterThan("age", 30)
query.GreaterThanEqual("age", 30)
query.Search("title", "hello")
query.IsNull("deletedAt")
query.IsNotNull("createdAt")
query.Between("age", 10, 20)
query.StartsWith("name", "A")
query.EndsWith("name", "n")
query.Contains("tags", "go")
query.Select([]string{"name", "age"})
query.OrderAsc("name")
query.OrderDesc("age")
query.CursorBefore("doc1")
query.CursorAfter("doc2")
query.Limit(10)
query.Offset(20)
query.Or([]string{query.Equal("a", 1), query.Equal("b", 2)})
query.And([]string{query.Equal("a", 1), query.Equal("b", 2)})
```
