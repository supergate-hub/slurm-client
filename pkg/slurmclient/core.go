package slurmclient

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

// CoreClient is the low-level HTTP client for Slurm REST API.
// It handles authentication, versioning, and common HTTP operations.
// Domain-specific clients (Jobs, Nodes, etc.) use this as their underlying transport.
type CoreClient struct {
	// HTTPClient is the underlying HTTP client.
	HTTPClient *http.Client

	// BaseURL is the base URL without version path, e.g. "http://host:6820/slurm"
	BaseURL string

	// APIVersion is the API version, e.g. "v0.0.40"
	APIVersion string

	// AuthToken is the JWT token for authentication.
	AuthToken string

	// UserAgent is the User-Agent header value.
	UserAgent string
}

// NewCoreClient creates a new CoreClient from the given configuration.
func NewCoreClient(cfg Config, domain Domain, version string) (*CoreClient, error) {
	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	normalizedVersion := normalizeVersion(version)
	baseURL := fmt.Sprintf("%s/%s", cfg.Endpoint, domain)

	httpClient := cfg.HTTPClient
	if httpClient == nil {
		httpClient = defaultHTTPClient(cfg.TLSConfig)
	}

	return &CoreClient{
		HTTPClient: httpClient,
		BaseURL:    baseURL,
		APIVersion: normalizedVersion,
		AuthToken:  cfg.AuthToken,
		UserAgent:  "slurm-go-sdk/1.0",
	}, nil
}

// URL builds the full URL for a given path.
func (c *CoreClient) URL(path string) string {
	return fmt.Sprintf("%s/%s%s", c.BaseURL, c.APIVersion, path)
}

// Do executes an HTTP request and decodes the response into v.
func (c *CoreClient) Do(ctx context.Context, method, path string, query url.Values, body any, v any) error {
	// Build URL
	reqURL := c.URL(path)
	if len(query) > 0 {
		reqURL += "?" + query.Encode()
	}

	// Prepare body
	var bodyReader io.Reader
	if body != nil {
		jsonBody, err := json.Marshal(body)
		if err != nil {
			return fmt.Errorf("failed to marshal request body: %w", err)
		}
		bodyReader = bytes.NewReader(jsonBody)
	}

	// Create request
	req, err := http.NewRequestWithContext(ctx, method, reqURL, bodyReader)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	req.Header.Set("Authorization", "Bearer "+c.AuthToken)
	req.Header.Set("Accept", "application/json")
	if c.UserAgent != "" {
		req.Header.Set("User-Agent", c.UserAgent)
	}
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	// Execute request
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	// Read response body
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}

	// Check for HTTP errors
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return &APIError{
			StatusCode: resp.StatusCode,
			Status:     resp.Status,
			Body:       string(respBody),
		}
	}

	// Decode response if v is provided
	if v != nil && len(respBody) > 0 {
		if err := json.Unmarshal(respBody, v); err != nil {
			return fmt.Errorf("failed to decode response: %w", err)
		}
	}

	return nil
}

// Get performs an HTTP GET request.
func (c *CoreClient) Get(ctx context.Context, path string, query url.Values, v any) error {
	return c.Do(ctx, http.MethodGet, path, query, nil, v)
}

// Post performs an HTTP POST request.
func (c *CoreClient) Post(ctx context.Context, path string, body any, v any) error {
	return c.Do(ctx, http.MethodPost, path, nil, body, v)
}

// Delete performs an HTTP DELETE request.
func (c *CoreClient) Delete(ctx context.Context, path string, query url.Values) error {
	return c.Do(ctx, http.MethodDelete, path, query, nil, nil)
}

// APIError represents an error response from the Slurm REST API.
type APIError struct {
	StatusCode int
	Status     string
	Body       string
}

func (e *APIError) Error() string {
	return fmt.Sprintf("slurm api error: %s (status: %d): %s", e.Status, e.StatusCode, e.Body)
}

