package slurmclient

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

// Request executes an HTTP request and returns the response.
// It handles authentication and content-type headers automatically.
func (c *ServiceClient) Request(ctx context.Context, method, path string, body any) (*http.Response, error) {
	return c.RequestWithQuery(ctx, method, path, nil, body)
}

// RequestWithQuery executes an HTTP request with query parameters.
// The path can be either a relative path (e.g., "/nodes") or a full URL from ServiceURL().
func (c *ServiceClient) RequestWithQuery(ctx context.Context, method, path string, query url.Values, body any) (*http.Response, error) {
	// Build the full URL
	// If path already contains the full URL (starts with http), use it directly
	var reqURL string
	if strings.HasPrefix(path, "http://") || strings.HasPrefix(path, "https://") {
		reqURL = path
	} else {
		reqURL = c.Endpoint + path
	}
	if len(query) > 0 {
		reqURL += "?" + query.Encode()
	}

	// Prepare the request body
	var bodyReader io.Reader
	if body != nil {
		jsonBody, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request body: %w", err)
		}
		bodyReader = bytes.NewReader(jsonBody)
	}

	// Create the request
	req, err := http.NewRequestWithContext(ctx, method, reqURL, bodyReader)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	req.Header.Set("Authorization", "Bearer "+c.AuthToken)
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	req.Header.Set("Accept", "application/json")

	// Execute the request
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}

	return resp, nil
}

// Get performs an HTTP GET request.
func (c *ServiceClient) Get(ctx context.Context, path string) (*http.Response, error) {
	return c.Request(ctx, http.MethodGet, path, nil)
}

// GetWithQuery performs an HTTP GET request with query parameters.
func (c *ServiceClient) GetWithQuery(ctx context.Context, path string, query url.Values) (*http.Response, error) {
	return c.RequestWithQuery(ctx, http.MethodGet, path, query, nil)
}

// Post performs an HTTP POST request.
func (c *ServiceClient) Post(ctx context.Context, path string, body any) (*http.Response, error) {
	return c.Request(ctx, http.MethodPost, path, body)
}

// Delete performs an HTTP DELETE request.
func (c *ServiceClient) Delete(ctx context.Context, path string) (*http.Response, error) {
	return c.Request(ctx, http.MethodDelete, path, nil)
}

// DeleteWithQuery performs an HTTP DELETE request with query parameters.
func (c *ServiceClient) DeleteWithQuery(ctx context.Context, path string, query url.Values) (*http.Response, error) {
	return c.RequestWithQuery(ctx, http.MethodDelete, path, query, nil)
}

// ResponseError represents an error response from the Slurm REST API.
type ResponseError struct {
	StatusCode int
	Status     string
	Body       string
}

func (e *ResponseError) Error() string {
	return fmt.Sprintf("slurm api error: %s (status: %d): %s", e.Status, e.StatusCode, e.Body)
}

// CheckResponse checks if the response indicates an error and returns a ResponseError if so.
func CheckResponse(resp *http.Response) error {
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		return nil
	}

	body, _ := io.ReadAll(resp.Body)
	resp.Body.Close()

	return &ResponseError{
		StatusCode: resp.StatusCode,
		Status:     resp.Status,
		Body:       string(body),
	}
}

// DecodeResponse reads and decodes the JSON response body into the given value.
// It also checks for HTTP errors.
func DecodeResponse[T any](resp *http.Response, v *T) error {
	defer resp.Body.Close()

	if err := CheckResponse(resp); err != nil {
		return err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}

	if err := json.Unmarshal(body, v); err != nil {
		return fmt.Errorf("failed to decode response: %w", err)
	}

	return nil
}

