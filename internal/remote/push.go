package remote

import (
	"encoding/json"
	"fmt"
)

const (
	PushEndpoint = "/functions/v1/push"
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

	return pushResponse, nil
}
