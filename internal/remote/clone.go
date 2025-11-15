package remote

import (
	"encoding/json"
	"fmt"
	"net/http"
)

const (
	CloneEndpoint = "/functions/v1/clone"
)

var ()

func (c *client) Clone(req CloneRequest) (CloneResponse, error) {
	url := fmt.Sprintf("%s%s", c.baseURL, CloneEndpoint)

	// Serialize request to JSON
	jsonData, err := json.Marshal(req)
	if err != nil {
		return CloneResponse{}, fmt.Errorf("failed to marshal request: %w", err)
	}

	// Make POST request
	resp, err := c.post(url, jsonData)
	if err != nil {
		return CloneResponse{}, err
	}
	defer resp.Body.Close()

	// Parse response
	var cloneResponse CloneResponse
	if err := json.NewDecoder(resp.Body).Decode(&cloneResponse); err != nil {
		return CloneResponse{}, fmt.Errorf("failed to decode response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		if resp.StatusCode == http.StatusUnauthorized {
			return CloneResponse{}, ErrUnauthorized
		}
		return CloneResponse{}, fmt.Errorf("clone failed: %s", cloneResponse.Message)
	}

	return cloneResponse, nil
}
