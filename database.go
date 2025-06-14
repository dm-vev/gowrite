package gowrite

import (
	"encoding/json"
	"fmt"
	"github.com/dm-vev/gowrite/query"
	"net/url"
	"sync"

	jsoniter "github.com/json-iterator/go"
)

var _json = jsoniter.ConfigFastest

type DatabaseService struct {
	Client *AppwriteClient
}

// Database represents an Appwrite database.
type Database struct {
	ID      string `json:"$id"`
	Name    string `json:"name"`
	Enabled bool   `json:"enabled"`
}

// Collection represents an Appwrite collection.
type Collection struct {
	ID               string        `json:"$id"`
	Name             string        `json:"name"`
	Permissions      []string      `json:"$permissions"`
	DocumentSecurity bool          `json:"documentSecurity"`
	Enabled          bool          `json:"enabled"`
	Attributes       []interface{} `json:"attributes"`
	Indexes          []interface{} `json:"indexes"`
}

// Document represents an Appwrite document.
type Document struct {
	ID          string                 `json:"$id"`
	Collection  string                 `json:"$collectionId"`
	Database    string                 `json:"$databaseId"`
	Permissions []string               `json:"$permissions"`
	Data        map[string]interface{} `json:"-"`
}

// Attribute represents a collection attribute.
type Attribute struct {
	Key      string `json:"key"`
	Type     string `json:"type"`
	Required bool   `json:"required"`
	Array    bool   `json:"array"`
}

// Permission constants
const (
	ReadAny    = "read(\"any\")"
	WriteAny   = "write(\"any\")"
	ReadUsers  = "read(\"users\")"
	WriteUsers = "write(\"users\")"
)

func NewDatabases(client *AppwriteClient) *DatabaseService {
	return &DatabaseService{client}
}

// ListDatabases retrieves a list of databases.
func (db *DatabaseService) ListDatabases() ([]*Database, error) {
	respBody, err := db.Client.sendRequest("GET", "/databases", nil)
	if err != nil {
		return nil, err
	}

	var result struct {
		Databases []*Database `json:"databases"`
	}

	err = json.Unmarshal(respBody, &result)
	if err != nil {
		return nil, err
	}

	return result.Databases, nil
}

// CreateDatabase creates a new database.
func (db *DatabaseService) CreateDatabase(databaseID, name string, enabled bool) (*Database, error) {
	payload := map[string]interface{}{
		"databaseId": databaseID,
		"name":       name,
		"enabled":    enabled,
	}

	respBody, err := db.Client.sendRequest("POST", "/databases", payload)
	if err != nil {
		return nil, err
	}

	var database Database
	err = json.Unmarshal(respBody, &database)
	if err != nil {
		return nil, err
	}

	return &database, nil
}

// GetDatabase retrieves a database by its ID.
func (db *DatabaseService) GetDatabase(databaseID string) (*Database, error) {
	path := fmt.Sprintf("/databases/%s", databaseID)
	respBody, err := db.Client.sendRequest("GET", path, nil)
	if err != nil {
		return nil, err
	}

	var database Database
	err = json.Unmarshal(respBody, &database)
	if err != nil {
		return nil, err
	}

	return &database, nil
}

// UpdateDatabase updates a database.
func (db *DatabaseService) UpdateDatabase(databaseID, name string, enabled bool) (*Database, error) {
	payload := map[string]interface{}{
		"name":    name,
		"enabled": enabled,
	}

	path := fmt.Sprintf("/databases/%s", databaseID)
	respBody, err := db.Client.sendRequest("PUT", path, payload)
	if err != nil {
		return nil, err
	}

	var database Database
	err = json.Unmarshal(respBody, &database)
	if err != nil {
		return nil, err
	}

	return &database, nil
}

// DeleteDatabase deletes a database.
func (db *DatabaseService) DeleteDatabase(databaseID string) error {
	path := fmt.Sprintf("/databases/%s", databaseID)
	_, err := db.Client.sendRequest("DELETE", path, nil)
	return err
}

// ListCollections retrieves a list of collections in a database.
func (db *DatabaseService) ListCollections(databaseID string) ([]*Collection, error) {
	path := fmt.Sprintf("/databases/%s/collections", databaseID)
	respBody, err := db.Client.sendRequest("GET", path, nil)
	if err != nil {
		return nil, err
	}

	var result struct {
		Collections []*Collection `json:"collections"`
	}

	err = json.Unmarshal(respBody, &result)
	if err != nil {
		return nil, err
	}

	return result.Collections, nil
}

// CreateCollection creates a new collection.
func (db *DatabaseService) CreateCollection(databaseID, collectionID, name string, permissions []string, documentSecurity, enabled bool) (*Collection, error) {
	payload := map[string]interface{}{
		"collectionId":     collectionID,
		"name":             name,
		"permissions":      permissions,
		"documentSecurity": documentSecurity,
		"enabled":          enabled,
	}

	path := fmt.Sprintf("/databases/%s/collections", databaseID)
	respBody, err := db.Client.sendRequest("POST", path, payload)
	if err != nil {
		return nil, err
	}

	var collection Collection
	err = json.Unmarshal(respBody, &collection)
	if err != nil {
		return nil, err
	}

	return &collection, nil
}

// GetCollection retrieves a collection by its ID.
func (db *DatabaseService) GetCollection(databaseID, collectionID string) (*Collection, error) {
	path := fmt.Sprintf("/databases/%s/collections/%s", databaseID, collectionID)
	respBody, err := db.Client.sendRequest("GET", path, nil)
	if err != nil {
		return nil, err
	}

	var collection Collection
	err = json.Unmarshal(respBody, &collection)
	if err != nil {
		return nil, err
	}

	return &collection, nil
}

// UpdateCollection updates a collection.
func (db *DatabaseService) UpdateCollection(databaseID, collectionID, name string, permissions []string, documentSecurity, enabled bool) (*Collection, error) {
	payload := map[string]interface{}{
		"name":             name,
		"permissions":      permissions,
		"documentSecurity": documentSecurity,
		"enabled":          enabled,
	}

	path := fmt.Sprintf("/databases/%s/collections/%s", databaseID, collectionID)
	respBody, err := db.Client.sendRequest("PUT", path, payload)
	if err != nil {
		return nil, err
	}

	var collection Collection
	err = json.Unmarshal(respBody, &collection)
	if err != nil {
		return nil, err
	}

	return &collection, nil
}

// DeleteCollection deletes a collection.
func (db *DatabaseService) DeleteCollection(databaseID, collectionID string) error {
	path := fmt.Sprintf("/databases/%s/collections/%s", databaseID, collectionID)
	_, err := db.Client.sendRequest("DELETE", path, nil)
	return err
}

// CreateDocument creates a new document.
func (db *DatabaseService) CreateDocument(databaseID, collectionID, documentID string, data map[string]interface{}, permissions []string) (*Document, error) {
	payload := map[string]interface{}{
		"documentId":  documentID,
		"data":        data,
		"permissions": permissions,
	}

	path := fmt.Sprintf("/databases/%s/collections/%s/documents", databaseID, collectionID)
	respBody, err := db.Client.sendRequest("POST", path, payload)
	if err != nil {
		return nil, err
	}

	var document Document
	err = json.Unmarshal(respBody, &document)
	if err != nil {
		return nil, err
	}

	return &document, nil
}

func (d *Document) UnmarshalJSON(b []byte) error {
	// Временная структура для известных полей
	type Alias Document
	aux := &struct {
		*Alias
	}{
		Alias: (*Alias)(d),
	}

	// Декодируем известные поля
	if err := json.Unmarshal(b, &aux); err != nil {
		return err
	}

	// Декодируем весь JSON в карту
	var raw map[string]interface{}
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}

	// Удаляем системные поля из карты и сохраняем в Data
	delete(raw, "$id")
	delete(raw, "$createdAt")
	delete(raw, "$updatedAt")
	delete(raw, "$permissions")
	delete(raw, "$databaseId")
	delete(raw, "$collectionId")
	d.Data = raw

	return nil
}

// GetDocument retrieves a document by its ID.
func (db *DatabaseService) GetDocument(databaseID, collectionID, documentID string) (*Document, error) {
	path := fmt.Sprintf("/databases/%s/collections/%s/documents/%s", databaseID, collectionID, documentID)
	respBody, err := db.Client.sendRequest("GET", path, nil)
	if err != nil {
		return nil, err
	}

	var document Document
	err = document.UnmarshalJSON(respBody)
	if err != nil {
		return nil, err
	}

	return &document, nil
}

// UpdateDocument updates a document.
func (db *DatabaseService) UpdateDocument(databaseID, collectionID, documentID string, data map[string]interface{}, permissions []string) (*Document, error) {
	payload := map[string]interface{}{
		"data":        data,
		"permissions": permissions,
	}

	path := fmt.Sprintf("/databases/%s/collections/%s/documents/%s", databaseID, collectionID, documentID)
	respBody, err := db.Client.sendRequest("PATCH", path, payload)
	if err != nil {
		return nil, err
	}

	var document Document
	err = json.Unmarshal(respBody, &document)
	if err != nil {
		return nil, err
	}

	return &document, nil
}

// DeleteDocument deletes a document.
func (db *DatabaseService) DeleteDocument(databaseID, collectionID, documentID string) error {
	path := fmt.Sprintf("/databases/%s/collections/%s/documents/%s", databaseID, collectionID, documentID)
	_, err := db.Client.sendRequest("DELETE", path, nil)
	return err
}

// ListDocuments получает список всех документов в коллекции, обрабатывая пагинацию для получения
// всех документов, превышающих лимит в 5000 за один запрос.
func (db *DatabaseService) ListDocuments(databaseID, collectionID string, queries []string) ([]*Document, error) {
	const (
		maxLimit    = 800
		concurrency = 5
	)

	// Предварительно фильтруем запросы, убирая limit и offset
	baseQueries := make([]string, 0, len(queries))
	for _, q := range queries {
		var pq query.QueryOptions
		if err := _json.UnmarshalFromString(q, &pq); err == nil &&
			pq.Method != "limit" && pq.Method != "offset" {
			baseQueries = append(baseQueries, q)
		}
	}

	type pageResult struct {
		docs    []*Document
		err     error
		hasMore bool
	}

	fetchPage := func(off int) pageResult {
		q := url.Values{}
		for _, qs := range baseQueries {
			q.Add("queries[]", qs)
		}
		q.Add("queries[]", query.Limit(maxLimit))
		q.Add("queries[]", query.Offset(off))

		path := fmt.Sprintf("/databases/%s/collections/%s/documents?%s", databaseID, collectionID, q.Encode())

		respBody, err := db.Client.sendRequest("GET", path, nil)
		if err != nil {
			return pageResult{nil, err, false}
		}

		var result struct {
			Documents []*Document `json:"documents"`
		}

		if err = _json.Unmarshal(respBody, &result); err != nil {
			return pageResult{nil, err, false}
		}

		hasMore := len(result.Documents) == maxLimit
		return pageResult{result.Documents, nil, hasMore}
	}

	var (
		allDocs []*Document
		off     int
		mu      sync.Mutex
		wg      sync.WaitGroup
		errOnce sync.Once
		retErr  error
	)

	worker := func() {
		defer wg.Done()
		for {
			mu.Lock()
			myOff := off
			off += maxLimit
			mu.Unlock()

			res := fetchPage(myOff)
			if res.err != nil {
				errOnce.Do(func() { retErr = res.err })
				return
			}

			if len(res.docs) == 0 {
				return
			}

			mu.Lock()
			allDocs = append(allDocs, res.docs...)
			mu.Unlock()

			if !res.hasMore {
				return
			}
		}
	}

	for i := 0; i < concurrency; i++ {
		wg.Add(1)
		go worker()
	}

	wg.Wait()
	if retErr != nil {
		return nil, retErr
	}
	return allDocs, nil
}

func (db *DatabaseService) CountDocuments(databaseID, collectionID string, queries []string) (int, error) {
	const maxLimit = 800

	// Предварительно фильтруем запросы, убирая limit и offset
	baseQueries := make([]string, 0, len(queries))
	for _, q := range queries {
		var pq query.QueryOptions
		if err := _json.UnmarshalFromString(q, &pq); err == nil &&
			pq.Method != "limit" && pq.Method != "offset" {
			baseQueries = append(baseQueries, q)
		}
	}

	offset := 0
	totalCount := 0

	for {
		q := url.Values{}
		for _, queryStr := range baseQueries {
			q.Add("queries[]", queryStr)
		}
		q.Add("queries[]", query.Limit(maxLimit))
		q.Add("queries[]", query.Offset(offset))

		path := fmt.Sprintf("/databases/%s/collections/%s/documents?%s", databaseID, collectionID, q.Encode())

		respBody, err := db.Client.sendRequest("GET", path, nil)
		if err != nil {
			return 0, err
		}

		var result struct {
			Documents []*Document `json:"documents"`
		}

		if err = _json.Unmarshal(respBody, &result); err != nil {
			return 0, err
		}

		count := len(result.Documents)
		totalCount += count

		if count < maxLimit {
			break
		}

		offset += maxLimit
	}

	return totalCount, nil
}

// AttributeType defines allowed attribute types when creating attributes.
type AttributeType string

const (
	AttributeBoolean      AttributeType = "boolean"
	AttributeDatetime     AttributeType = "datetime"
	AttributeEmail        AttributeType = "email"
	AttributeEnum         AttributeType = "enum"
	AttributeFloat        AttributeType = "float"
	AttributeInteger      AttributeType = "integer"
	AttributeIP           AttributeType = "ip"
	AttributeRelationship AttributeType = "relationship"
	AttributeString       AttributeType = "string"
	AttributeURL          AttributeType = "url"
)

// CreateAttribute creates a new attribute for a collection.
func (db *DatabaseService) CreateAttribute(databaseID, collectionID, key string, attrType AttributeType, required bool, defaultValue interface{}, array bool, meta map[string]interface{}) (*Attribute, error) {
	payload := map[string]interface{}{
		"key":      key,
		"required": required,
		"default":  defaultValue,
		"array":    array,
	}
	for k, v := range meta {
		payload[k] = v
	}

	path := fmt.Sprintf("/databases/%s/collections/%s/attributes/%s", databaseID, collectionID, attrType)
	respBody, err := db.Client.sendRequest("POST", path, payload)
	if err != nil {
		return nil, err
	}

	var attr Attribute
	if err = json.Unmarshal(respBody, &attr); err != nil {
		return nil, err
	}

	return &attr, nil
}

// GetAttribute retrieves an attribute from a collection.
func (db *DatabaseService) GetAttribute(databaseID, collectionID, key string) (*Attribute, error) {
	path := fmt.Sprintf("/databases/%s/collections/%s/attributes/%s", databaseID, collectionID, key)
	respBody, err := db.Client.sendRequest("GET", path, nil)
	if err != nil {
		return nil, err
	}

	var attr Attribute
	if err = json.Unmarshal(respBody, &attr); err != nil {
		return nil, err
	}

	return &attr, nil
}

// DeleteAttribute deletes an attribute from a collection.
func (db *DatabaseService) DeleteAttribute(databaseID, collectionID, key string) error {
	path := fmt.Sprintf("/databases/%s/collections/%s/attributes/%s", databaseID, collectionID, key)
	_, err := db.Client.sendRequest("DELETE", path, nil)
	return err
}

// ListAttributes retrieves all attributes from a collection.
func (db *DatabaseService) ListAttributes(databaseID, collectionID string, queries []string) ([]*Attribute, error) {
	q := url.Values{}
	for _, qs := range queries {
		q.Add("queries[]", qs)
	}
	path := fmt.Sprintf("/databases/%s/collections/%s/attributes", databaseID, collectionID)
	if encoded := q.Encode(); encoded != "" {
		path += "?" + encoded
	}
	respBody, err := db.Client.sendRequest("GET", path, nil)
	if err != nil {
		return nil, err
	}

	var result struct {
		Attributes []*Attribute `json:"attributes"`
	}
	if err = json.Unmarshal(respBody, &result); err != nil {
		return nil, err
	}

	return result.Attributes, nil
}

// UpdateAttribute updates an attribute of a given type.
func (db *DatabaseService) UpdateAttribute(databaseID, collectionID, key string, attrType AttributeType, updates map[string]interface{}) (*Attribute, error) {
	var path string
	if attrType == AttributeRelationship {
		path = fmt.Sprintf("/databases/%s/collections/%s/attributes/%s/relationship", databaseID, collectionID, key)
	} else {
		path = fmt.Sprintf("/databases/%s/collections/%s/attributes/%s/%s", databaseID, collectionID, attrType, key)
	}

	respBody, err := db.Client.sendRequest("PATCH", path, updates)
	if err != nil {
		return nil, err
	}

	var attr Attribute
	if err = json.Unmarshal(respBody, &attr); err != nil {
		return nil, err
	}

	return &attr, nil
}
