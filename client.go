package appwrite

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type AppwriteClient struct {
	Endpoint   string
	ProjectID  string
	APIKey     string
	HTTPClient *http.Client
}

func NewClient(endpoint, projectID, apiKey string) *AppwriteClient {
	return &AppwriteClient{
		Endpoint:   endpoint,
		ProjectID:  projectID,
		APIKey:     apiKey,
		HTTPClient: &http.Client{},
	}
}

func (client *AppwriteClient) sendRequest(method, path string, body interface{}) ([]byte, error) {
	url := fmt.Sprintf("%s/v1%s", client.Endpoint, path)

	var reqBody []byte
	var err error

	if body != nil {
		reqBody, err = json.Marshal(body)
		if err != nil {
			return nil, err
		}
	}

	req, err := http.NewRequest(method, url, bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Appwrite-Project", client.ProjectID)
	req.Header.Set("X-Appwrite-Key", client.APIKey)
	req.Header.Set("X-Appwrite-Response-Format", "1.6.0")

	resp, err := client.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBody, err := ioutil.ReadAll(resp.Body)
	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("HTTP %d: %s", resp.StatusCode, string(respBody))
	}

	if err != nil {
		return nil, err
	}

	return respBody, nil
}

func (client *AppwriteClient) SendRequest(method, path string, body interface{}) ([]byte, error) {
	return client.sendRequest(method, path, body)
}
