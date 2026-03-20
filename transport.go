package slurmclient

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"errors"
	"log/slog"
	"math"
	"net"
	"net/http"
	"strings"
	"time"
)

// Transport handles HTTP communication with slurmrestd.
// It is the single point for auth, retry, logging, and URL construction.
type Transport struct {
	baseURL    string
	domain     string // "slurm" or "slurmdb"
	version    Version
	authToken  string
	httpClient *http.Client
	logger     *slog.Logger

	// Retry configuration
	maxRetries int
	baseDelay  time.Duration
}

// TransportConfig configures a Transport.
type TransportConfig struct {
	BaseURL    string
	Domain     string
	Version    Version
	AuthToken  string
	HTTPClient *http.Client
	Logger     *slog.Logger
	MaxRetries int
	BaseDelay  time.Duration
}

func newTransport(cfg TransportConfig) *Transport {
	t := &Transport{
		baseURL:    strings.TrimRight(cfg.BaseURL, "/"),
		domain:     cfg.Domain,
		version:    cfg.Version,
		authToken:  cfg.AuthToken,
		httpClient: cfg.HTTPClient,
		logger:     cfg.Logger,
		maxRetries: cfg.MaxRetries,
		baseDelay:  cfg.BaseDelay,
	}
	if t.httpClient == nil {
		t.httpClient = &http.Client{Timeout: 30 * time.Second}
	}
	if t.maxRetries <= 0 {
		t.maxRetries = 3
	}
	if t.baseDelay <= 0 {
		t.baseDelay = 500 * time.Millisecond
	}
	return t
}

// newUnixSocketHTTPClient creates an HTTP client that connects via Unix socket.
func newUnixSocketHTTPClient(socketPath string) *http.Client {
	return &http.Client{
		Timeout: 30 * time.Second,
		Transport: &http.Transport{
			DialContext: func(_ context.Context, _, _ string) (net.Conn, error) {
				return net.Dial("unix", socketPath)
			},
		},
	}
}

// resourceURL builds the full URL for a resource path.
// Example: baseURL + "/slurm/v0.0.40" + "/jobs/123"
func (t *Transport) resourceURL(resourcePath string) string {
	if !strings.HasPrefix(resourcePath, "/") {
		resourcePath = "/" + resourcePath
	}
	return fmt.Sprintf("%s/%s/%s%s", t.baseURL, t.domain, t.version, resourcePath)
}

// Get performs a GET request and decodes the response into result.
func (t *Transport) Get(ctx context.Context, resourcePath string, result any) error {
	return t.do(ctx, http.MethodGet, resourcePath, nil, result)
}

// Post performs a POST request with a JSON body and decodes the response.
func (t *Transport) Post(ctx context.Context, resourcePath string, body, result any) error {
	return t.do(ctx, http.MethodPost, resourcePath, body, result)
}

// Put performs a PUT request with a JSON body and decodes the response.
func (t *Transport) Put(ctx context.Context, resourcePath string, body, result any) error {
	return t.do(ctx, http.MethodPut, resourcePath, body, result)
}

// Delete performs a DELETE request and optionally decodes the response.
func (t *Transport) Delete(ctx context.Context, resourcePath string, result any) error {
	return t.do(ctx, http.MethodDelete, resourcePath, nil, result)
}

func (t *Transport) do(ctx context.Context, method, resourcePath string, body, result any) error {
	url := t.resourceURL(resourcePath)

	var bodyReader io.Reader
	if body != nil {
		data, err := json.Marshal(body)
		if err != nil {
			return fmt.Errorf("slurm: marshal request body: %w", err)
		}
		bodyReader = bytes.NewReader(data)
	}

	var lastErr error
	for attempt := 0; attempt <= t.maxRetries; attempt++ {
		if attempt > 0 {
			delay := t.baseDelay * time.Duration(math.Pow(2, float64(attempt-1)))
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.After(delay):
			}
			t.log(ctx, slog.LevelWarn, "retrying request",
				"method", method, "url", url,
				"attempt", attempt, "max_retries", t.maxRetries)
		}

		// Reset body reader for retries
		if body != nil {
			data, _ := json.Marshal(body)
			bodyReader = bytes.NewReader(data)
		}

		req, err := http.NewRequestWithContext(ctx, method, url, bodyReader)
		if err != nil {
			return fmt.Errorf("slurm: create request: %w", err)
		}
		req.Header.Set("Authorization", "Bearer "+t.authToken)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Accept", "application/json")

		t.log(ctx, slog.LevelDebug, "sending request",
			"method", method, "url", url)

		resp, err := t.httpClient.Do(req)
		if err != nil {
			lastErr = fmt.Errorf("slurm: %s %s: %w", method, resourcePath, err)
			if isRetryable(0, err) {
				continue
			}
			return lastErr
		}

		respBody, err := io.ReadAll(resp.Body)
		resp.Body.Close()

		t.log(ctx, slog.LevelDebug, "received response",
			"method", method, "url", url,
			"status", resp.StatusCode, "body_len", len(respBody))

		if err != nil {
			lastErr = fmt.Errorf("slurm: read response body: %w", err)
			continue
		}

		if isRetryable(resp.StatusCode, nil) {
			lastErr = parseSlurmError(resp.StatusCode, respBody)
			continue
		}

		if resp.StatusCode == http.StatusUnauthorized {
			return ErrUnauthorized
		}
		if resp.StatusCode == http.StatusNotFound {
			return ErrNotFound
		}
		if resp.StatusCode >= 400 {
			return parseSlurmError(resp.StatusCode, respBody)
		}

		// Decode response
		if result != nil && len(respBody) > 0 {
			if err := json.Unmarshal(respBody, result); err != nil {
				return fmt.Errorf("slurm: decode response: %w", err)
			}
		}
		return nil
	}

	return &RetryExhaustedError{Attempts: t.maxRetries, Last: lastErr}
}

// isRetryable returns true if the request should be retried.
func isRetryable(statusCode int, err error) bool {
	if statusCode == http.StatusServiceUnavailable ||
		statusCode == http.StatusTooManyRequests ||
		statusCode == http.StatusBadGateway ||
		statusCode == http.StatusGatewayTimeout {
		return true
	}
	if err != nil {
		// Retry on connection reset, timeout, etc.
		var netErr net.Error
		if errors.As(err, &netErr) {
			return true
		}
	}
	return false
}

// parseSlurmError parses the Slurm API error response body.
func parseSlurmError(statusCode int, body []byte) *SlurmError {
	se := &SlurmError{StatusCode: statusCode}

	var resp struct {
		Errors   []struct{ Error *string `json:"error"` }   `json:"errors"`
		Warnings []struct{ Warning *string `json:"warning"` } `json:"warnings"`
	}
	if json.Unmarshal(body, &resp) == nil {
		for _, e := range resp.Errors {
			if e.Error != nil {
				se.Errors = append(se.Errors, *e.Error)
			}
		}
		for _, w := range resp.Warnings {
			if w.Warning != nil {
				se.Warnings = append(se.Warnings, *w.Warning)
			}
		}
	}
	return se
}

// log writes a structured log entry if a logger is configured.
// Authorization headers are always masked.
func (t *Transport) log(ctx context.Context, level slog.Level, msg string, args ...any) {
	if t.logger == nil {
		return
	}
	t.logger.Log(ctx, level, msg, args...)
}
