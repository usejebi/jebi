package remote

import (
	"encoding/json"
	"fmt"
	"net/http"
)

const (
	PushEndpoint = "/functions/v1/push"
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

	if resp.StatusCode != http.StatusCreated {
		if resp.StatusCode == http.StatusUnauthorized {
			return PushResponse{}, ErrUnauthorized
		}
		var errorResponse ErrorResponse
		if err := json.NewDecoder(resp.Body).Decode(&errorResponse); err != nil {
			return PushResponse{}, fmt.Errorf("failed to decode error response: %w", err)
		}

		if resp.StatusCode == http.StatusConflict && errorResponse.Code == "PROJECT_NAME_ALREADY_EXISTS" {
			return PushResponse{}, ErrProjectNameAlreadyExists
		}

		return PushResponse{}, fmt.Errorf("push failed: %s", errorResponse.Message)
	}

	// Parse response
	var pushResponse PushResponse
	if err := json.NewDecoder(resp.Body).Decode(&pushResponse); err != nil {
		return PushResponse{}, fmt.Errorf("failed to decode response: %w", err)
	}

	return pushResponse, nil
}
