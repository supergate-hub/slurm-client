package slurmclient

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strings"
	"time"
)

// AuthOpts configures authentication and connection to slurmrestd.
type AuthOpts struct {
	// Endpoint is the slurmrestd HTTP endpoint (e.g., "http://slurmctld:6820").
	// Required unless UnixSocket is set.
	Endpoint string

	// AuthToken is the JWT token for slurmrestd authentication.
	AuthToken string

	// UnixSocket is the path to the slurmrestd Unix socket.
	// When set, Endpoint is used only for URL construction (default: "http://localhost").
	UnixSocket string

	// Version selects the Slurm REST API version.
	// If empty, the SDK auto-detects by querying slurmrestd's /openapi endpoint.
	Version Version

	// HTTPClient is an optional custom HTTP client.
	// When UnixSocket is set, this is ignored.
	HTTPClient *http.Client

	// Logger enables structured logging of HTTP requests/responses.
	// Authorization headers are always masked in logs.
	Logger *slog.Logger

	// MaxRetries is the maximum number of retry attempts for transient failures.
	// Defaults to 3.
	MaxRetries int

	// RetryBaseDelay is the base delay for exponential backoff.
	// Defaults to 500ms.
	RetryBaseDelay time.Duration
}

// Client is the top-level SDK client for the Slurm REST API.
// Use NewClient to create one.
//
// Usage:
//
//	client, err := slurmclient.NewClient(ctx, slurmclient.AuthOpts{
//	    Endpoint:  "http://slurmctld:6820",
//	    AuthToken: "jwt-token",
//	    Version:   slurmclient.V0040,
//	})
//	jobs, err := client.Slurm.Jobs().List(ctx)
type Client struct {
	// Slurm provides access to the Slurm API resources (jobs, nodes, partitions, etc.).
	Slurm SlurmProxy

	// Slurmdb provides access to the SlurmDB API resources (accounts, users, etc.).
	Slurmdb SlurmdbProxy

	slurmTransport   *Transport
	slurmdbTransport *Transport
	version          Version
}

// NewClient creates a new Slurm SDK client.
//
// If Version is not set, the SDK auto-detects the API version by querying
// slurmrestd's /openapi endpoint.
func NewClient(ctx context.Context, opts AuthOpts) (*Client, error) {
	if err := validateOpts(opts); err != nil {
		return nil, err
	}

	endpoint := opts.Endpoint
	httpClient := opts.HTTPClient

	if opts.UnixSocket != "" {
		httpClient = newUnixSocketHTTPClient(opts.UnixSocket)
		if endpoint == "" {
			endpoint = "http://localhost"
		}
	}

	version := opts.Version
	if version == "" {
		discovered, err := discoverVersion(ctx, endpoint, opts.AuthToken, httpClient, opts.Logger)
		if err != nil {
			return nil, err
		}
		version = discovered
	}

	slurmT := newTransport(TransportConfig{
		BaseURL:    endpoint,
		Domain:     "slurm",
		Version:    version,
		AuthToken:  opts.AuthToken,
		HTTPClient: httpClient,
		Logger:     opts.Logger,
		MaxRetries: opts.MaxRetries,
		BaseDelay:  opts.RetryBaseDelay,
	})

	slurmdbT := newTransport(TransportConfig{
		BaseURL:    endpoint,
		Domain:     "slurmdb",
		Version:    version,
		AuthToken:  opts.AuthToken,
		HTTPClient: httpClient,
		Logger:     opts.Logger,
		MaxRetries: opts.MaxRetries,
		BaseDelay:  opts.RetryBaseDelay,
	})

	c := &Client{
		slurmTransport:   slurmT,
		slurmdbTransport: slurmdbT,
		version:          version,
	}

	c.Slurm = newSlurmProxy(slurmT)
	c.Slurmdb = newSlurmdbProxy(slurmdbT)

	return c, nil
}

// Version returns the API version this client is using.
func (c *Client) Version() Version {
	return c.version
}

func validateOpts(opts AuthOpts) error {
	if opts.Endpoint == "" && opts.UnixSocket == "" {
		return fmt.Errorf("%w: endpoint or unix socket is required", ErrValidation)
	}
	if opts.AuthToken == "" {
		return fmt.Errorf("%w: auth token is required", ErrValidation)
	}
	if opts.Version != "" && !opts.Version.IsValid() {
		return fmt.Errorf("%w: unsupported version %q", ErrValidation, opts.Version)
	}
	return nil
}

// discoverVersion queries slurmrestd to find the best supported API version.
func discoverVersion(ctx context.Context, endpoint, token string, httpClient *http.Client, logger *slog.Logger) (Version, error) {
	if httpClient == nil {
		httpClient = &http.Client{Timeout: 10 * time.Second}
	}

	url := strings.TrimRight(endpoint, "/") + "/openapi/v3"
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return "", &VersionDiscoveryError{Err: err}
	}
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Accept", "application/json")

	resp, err := httpClient.Do(req)
	if err != nil {
		return "", &VersionDiscoveryError{Err: fmt.Errorf("request to %s failed: %w", url, err)}
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", &VersionDiscoveryError{
			Err: fmt.Errorf("GET %s returned %d", url, resp.StatusCode),
		}
	}

	var spec struct {
		Paths map[string]any `json:"paths"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&spec); err != nil {
		return "", &VersionDiscoveryError{Err: fmt.Errorf("parse openapi spec: %w", err)}
	}

	// Scan paths for version prefixes like /slurm/v0.0.40/
	found := make(map[Version]bool)
	for path := range spec.Paths {
		for _, v := range supportedVersions {
			if strings.HasPrefix(path, "/slurm/"+string(v)+"/") ||
				strings.HasPrefix(path, "/slurmdb/"+string(v)+"/") {
				found[v] = true
			}
		}
	}

	// Return the newest supported version
	for _, v := range supportedVersions {
		if found[v] {
			if logger != nil {
				logger.Info("auto-detected slurm API version", "version", v)
			}
			return v, nil
		}
	}

	return "", &VersionDiscoveryError{Err: fmt.Errorf("no supported API version found in slurmrestd")}
}
