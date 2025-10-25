package remote

import (
	"bytes"
	"fmt"
	"net/http"

	"github.com/jawahars16/jebi/internal/keystore"
)

type client struct {
	baseURL    string
	httpClient http.Client
	keystore   keystore.KeyStore
}

func NewAPIClient(baseURL string) *client {
	return &client{
		baseURL:    baseURL,
		httpClient: *http.DefaultClient,
		keystore:   keystore.NewDefault("."), // Use current directory for keystore
	}
}

func (c *client) post(url string, jsonData []byte) (*http.Response, error) {
	httpReq, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set content type
	httpReq.Header.Set("Content-Type", "application/json")

	// Load auth token from keystore and set Authorization header
	var accessToken string
	if err := c.keystore.Get("access_token", &accessToken); err == nil && accessToken != "" {
		httpReq.Header.Set("Authorization", fmt.Sprintf("Bearer %s", accessToken))
	}

	// Make the HTTP request
	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}

	// Check for HTTP errors
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		// Close the body here since we're not returning it
		resp.Body.Close()
		return nil, fmt.Errorf("API request failed with status %d", resp.StatusCode)
	}

	// Return the response with body open - caller is responsible for closing it
	return resp, nil
}
