package gowrite

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
)

type StorageService struct {
	Client *AppwriteClient
}

// Bucket представляет хранилище в Appwrite.
type Bucket struct {
	ID                    string   `json:"$id"`
	Name                  string   `json:"name"`
	Permissions           []string `json:"$permissions"`
	FileSecurity          bool     `json:"fileSecurity"`
	Enabled               bool     `json:"enabled"`
	MaximumFileSize       int64    `json:"maximumFileSize"`
	AllowedFileExtensions []string `json:"allowedFileExtensions"`
	Compression           string   `json:"compression"`
	Encryption            bool     `json:"encryption"`
	Antivirus             bool     `json:"antivirus"`
}

// File представляет файл в Appwrite.
type File struct {
	ID             string                 `json:"$id"`
	BucketID       string                 `json:"bucketId"`
	Name           string                 `json:"name"`
	Signature      string                 `json:"signature"`
	MimeType       string                 `json:"mimeType"`
	SizeOriginal   int64                  `json:"sizeOriginal"`
	Permissions    []string               `json:"$permissions"`
	ChunksTotal    int                    `json:"chunksTotal"`
	ChunksUploaded int                    `json:"chunksUploaded"`
	Data           map[string]interface{} `json:"-"`
}

func NewStorage(client *AppwriteClient) *StorageService {
	return &StorageService{client}
}

// Custom UnmarshalJSON для File
func (f *File) UnmarshalJSON(b []byte) error {
	type Alias File
	aux := &struct {
		*Alias
	}{
		Alias: (*Alias)(f),
	}

	if err := json.Unmarshal(b, &aux); err != nil {
		return err
	}

	var raw map[string]interface{}
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}

	// Удаляем известные поля
	delete(raw, "$id")
	delete(raw, "bucketId")
	delete(raw, "name")
	delete(raw, "signature")
	delete(raw, "mimeType")
	delete(raw, "sizeOriginal")
	delete(raw, "$permissions")
	delete(raw, "chunksTotal")
	delete(raw, "chunksUploaded")
	f.Data = raw

	return nil
}

// ListBuckets получает список всех бакетов.
func (s *StorageService) ListBuckets() ([]*Bucket, error) {
	respBody, err := s.Client.sendRequest("GET", "/storage/buckets", nil)
	if err != nil {
		return nil, err
	}

	var result struct {
		Buckets []*Bucket `json:"buckets"`
	}

	err = json.Unmarshal(respBody, &result)
	if err != nil {
		return nil, err
	}

	return result.Buckets, nil
}

// CreateBucket создает новый бакет.
func (s *StorageService) CreateBucket(bucketID, name string, permissions []string, fileSecurity, enabled bool, maximumFileSize int64, allowedFileExtensions []string, compression string, encryption, antivirus bool) (*Bucket, error) {
	payload := map[string]interface{}{
		"bucketId":              bucketID,
		"name":                  name,
		"permissions":           permissions,
		"fileSecurity":          fileSecurity,
		"enabled":               enabled,
		"maximumFileSize":       maximumFileSize,
		"allowedFileExtensions": allowedFileExtensions,
		"compression":           compression,
		"encryption":            encryption,
		"antivirus":             antivirus,
	}

	respBody, err := s.Client.sendRequest("POST", "/storage/buckets", payload)
	if err != nil {
		return nil, err
	}

	var bucket Bucket
	err = json.Unmarshal(respBody, &bucket)
	if err != nil {
		return nil, err
	}

	return &bucket, nil
}

// GetBucket получает бакет по его ID.
func (s *StorageService) GetBucket(bucketID string) (*Bucket, error) {
	path := fmt.Sprintf("/storage/buckets/%s", bucketID)
	respBody, err := s.Client.sendRequest("GET", path, nil)
	if err != nil {
		return nil, err
	}

	var bucket Bucket
	err = json.Unmarshal(respBody, &bucket)
	if err != nil {
		return nil, err
	}

	return &bucket, nil
}

// UpdateBucket обновляет бакет.
func (s *StorageService) UpdateBucket(bucketID, name string, permissions []string, fileSecurity, enabled bool, maximumFileSize int64, allowedFileExtensions []string, compression string, encryption, antivirus bool) (*Bucket, error) {
	payload := map[string]interface{}{
		"name":                  name,
		"permissions":           permissions,
		"fileSecurity":          fileSecurity,
		"enabled":               enabled,
		"maximumFileSize":       maximumFileSize,
		"allowedFileExtensions": allowedFileExtensions,
		"compression":           compression,
		"encryption":            encryption,
		"antivirus":             antivirus,
	}

	path := fmt.Sprintf("/storage/buckets/%s", bucketID)
	respBody, err := s.Client.sendRequest("PUT", path, payload)
	if err != nil {
		return nil, err
	}

	var bucket Bucket
	err = json.Unmarshal(respBody, &bucket)
	if err != nil {
		return nil, err
	}

	return &bucket, nil
}

// DeleteBucket удаляет бакет.
func (s *StorageService) DeleteBucket(bucketID string) error {
	path := fmt.Sprintf("/storage/buckets/%s", bucketID)
	_, err := s.Client.sendRequest("DELETE", path, nil)
	return err
}

// ListFiles получает список файлов в бакете.
func (s *StorageService) ListFiles(bucketID string) ([]*File, error) {
	path := fmt.Sprintf("/storage/buckets/%s/files", bucketID)
	respBody, err := s.Client.sendRequest("GET", path, nil)
	if err != nil {
		return nil, err
	}

	var result struct {
		Files []*File `json:"files"`
	}

	err = json.Unmarshal(respBody, &result)
	if err != nil {
		return nil, err
	}

	return result.Files, nil
}

// CreateFile загружает новый файл в бакет.
func (s *StorageService) CreateFile(bucketID, fileID, filePath string, permissions []string) (*File, error) {
	path := fmt.Sprintf("/storage/buckets/%s/files", bucketID)

	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// Добавляем fileId
	if err := writer.WriteField("fileId", fileID); err != nil {
		return nil, err
	}

	// Добавляем permissions
	for _, permission := range permissions {
		if err := writer.WriteField("permissions[]", permission); err != nil {
			return nil, err
		}
	}

	// Добавляем файл
	part, err := writer.CreateFormFile("file", filepath.Base(filePath))
	if err != nil {
		return nil, err
	}
	if _, err := io.Copy(part, file); err != nil {
		return nil, err
	}

	writer.Close()

	req, err := http.NewRequest("POST", s.Client.Endpoint+"/v1"+path, body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("X-Appwrite-Project", s.Client.ProjectID)
	req.Header.Set("X-Appwrite-Key", s.Client.APIKey)
	req.Header.Set("X-Appwrite-Response-Format", "1.6.0")

	resp, err := s.Client.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("HTTP %d: %s", resp.StatusCode, string(respBody))
	}

	var fileResp File
	err = json.Unmarshal(respBody, &fileResp)
	if err != nil {
		return nil, err
	}

	return &fileResp, nil
}

// GetFile получает файл по его ID.
func (s *StorageService) GetFile(bucketID, fileID string) (*File, error) {
	path := fmt.Sprintf("/storage/buckets/%s/files/%s", bucketID, fileID)
	respBody, err := s.Client.sendRequest("GET", path, nil)
	if err != nil {
		return nil, err
	}

	var file File
	err = json.Unmarshal(respBody, &file)
	if err != nil {
		return nil, err
	}

	return &file, nil
}

// UpdateFile обновляет файл.
func (s *StorageService) UpdateFile(bucketID, fileID, name string, permissions []string) (*File, error) {
	payload := map[string]interface{}{
		"name":        name,
		"permissions": permissions,
	}

	path := fmt.Sprintf("/storage/buckets/%s/files/%s", bucketID, fileID)
	respBody, err := s.Client.sendRequest("PUT", path, payload)
	if err != nil {
		return nil, err
	}

	var file File
	err = json.Unmarshal(respBody, &file)
	if err != nil {
		return nil, err
	}

	return &file, nil
}

// DeleteFile удаляет файл.
func (s *StorageService) DeleteFile(bucketID, fileID string) error {
	path := fmt.Sprintf("/storage/buckets/%s/files/%s", bucketID, fileID)
	_, err := s.Client.sendRequest("DELETE", path, nil)
	return err
}

// DownloadFile скачивает файл.
func (s *StorageService) DownloadFile(bucketID, fileID string) ([]byte, error) {
	path := fmt.Sprintf("/storage/buckets/%s/files/%s/download", bucketID, fileID)

	req, err := http.NewRequest("GET", s.Client.Endpoint+"/v1"+path, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("X-Appwrite-Project", s.Client.ProjectID)
	req.Header.Set("X-Appwrite-Key", s.Client.APIKey)
	req.Header.Set("X-Appwrite-Response-Format", "1.6.0")

	resp, err := s.Client.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		respBody, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("HTTP %d: %s", resp.StatusCode, string(respBody))
	}

	content, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return content, nil
}

// GetFilePreview получает превью файла.
func (s *StorageService) GetFilePreview(bucketID, fileID string, params map[string]string) ([]byte, error) {
	path := fmt.Sprintf("/storage/buckets/%s/files/%s/preview", bucketID, fileID)

	// Добавляем параметры
	if len(params) > 0 {
		q := "?"
		for key, value := range params {
			q += fmt.Sprintf("%s=%s&", key, value)
		}
		path += q[:len(q)-1]
	}

	req, err := http.NewRequest("GET", s.Client.Endpoint+"/v1"+path, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("X-Appwrite-Project", s.Client.ProjectID)
	req.Header.Set("X-Appwrite-Key", s.Client.APIKey)
	req.Header.Set("X-Appwrite-Response-Format", "1.6.0")

	resp, err := s.Client.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		respBody, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("HTTP %d: %s", resp.StatusCode, string(respBody))
	}

	content, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return content, nil
}

// ViewFile получает содержимое файла для просмотра.
func (s *StorageService) ViewFile(bucketID, fileID string) ([]byte, error) {
	path := fmt.Sprintf("/storage/buckets/%s/files/%s/view", bucketID, fileID)

	req, err := http.NewRequest("GET", s.Client.Endpoint+"/v1"+path, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("X-Appwrite-Project", s.Client.ProjectID)
	req.Header.Set("X-Appwrite-Key", s.Client.APIKey)
	req.Header.Set("X-Appwrite-Response-Format", "1.6.0")

	resp, err := s.Client.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		respBody, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("HTTP %d: %s", resp.StatusCode, string(respBody))
	}

	content, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return content, nil
}

// GetFileDownloadURL формирует URL для скачивания файла.
func (s *StorageService) GetFileDownloadURL(bucketID, fileID string) string {
	return fmt.Sprintf("%s/v1/storage/buckets/%s/files/%s/download?project=%s",
		s.Client.Endpoint, bucketID, fileID, s.Client.ProjectID)
}
