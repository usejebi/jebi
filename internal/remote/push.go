package remote

import (
	"encoding/json"
	"fmt"
	"net/http"
)

const (
	PushEndpoint = "/functions/v1/sync"
)

var (
	ErrProjectNameAlreadyExists = fmt.Errorf("project name already exists on remote")
	ErrUnauthorized             = fmt.Errorf("unauthorized access to remote server")
)

func (c *client) Push(req PushRequest) (PushResponse, error) {
	url := fmt.Sprintf("%s%s", c.baseURL, PushEndpoint)

	// Serialize request to JSON
	jsonData, err := json.Marshal(req)
	if err != nil {
		return PushResponse{}, fmt.Errorf("failed to marshal request: %w", err)
	}

	// Make POST request
	resp, err := c.post(url, jsonData)
	if err != nil {
		return PushResponse{}, err
	}
	defer resp.Body.Close()

	// Parse response
	var pushResponse PushResponse
	if err := json.NewDecoder(resp.Body).Decode(&pushResponse); err != nil {
		return PushResponse{}, fmt.Errorf("failed to decode response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		if resp.StatusCode == http.StatusUnauthorized {
			return PushResponse{}, ErrUnauthorized
		}
		if resp.StatusCode == http.StatusConflict && pushResponse.Code == "PROJECT_NAME_ALREADY_EXISTS" {
			return PushResponse{}, ErrProjectNameAlreadyExists
		}
		return PushResponse{}, fmt.Errorf("push failed: %s", pushResponse.Message)
	}

	return pushResponse, nil
}
