package slurmclient

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"strings"
	"time"
)

// Domain represents the Slurm REST API domain.
type Domain string

const (
	// DomainSlurm represents the slurmctld domain (/slurm/v0.0.XX/...).
	DomainSlurm Domain = "slurm"

	// DomainSlurmdb represents the slurmdbd domain (/slurmdb/v0.0.XX/...).
	DomainSlurmdb Domain = "slurmdb"
)

// Config holds the configuration for connecting to a Slurm REST API.
type Config struct {
	// Endpoint is the base URL of slurmrestd, e.g. "http://host:6820".
	// Required.
	Endpoint string

	// AuthToken is the Slurm JWT token for authentication.
	// Required.
	AuthToken string

	// TLSConfig is optional TLS configuration for HTTPS endpoints.
	TLSConfig *tls.Config

	// HTTPClient is an optional custom HTTP client (for timeouts, proxies, etc).
	// If nil, a default client with reasonable timeouts will be created.
	HTTPClient *http.Client
}

// Validate checks if the config has all required fields.
func (c *Config) Validate() error {
	if c == nil {
		return fmt.Errorf("config cannot be nil")
	}
	if c.Endpoint == "" {
		return fmt.Errorf("endpoint cannot be empty")
	}
	if c.AuthToken == "" {
		return fmt.Errorf("authToken cannot be empty")
	}
	// Normalize endpoint - remove trailing slash
	c.Endpoint = strings.TrimSuffix(c.Endpoint, "/")
	return nil
}

// defaultHTTPClient creates a default HTTP client with reasonable settings.
func defaultHTTPClient(tlsConfig *tls.Config) *http.Client {
	transport := &http.Transport{
		TLSClientConfig: tlsConfig,
	}
	return &http.Client{
		Transport: transport,
		Timeout:   30 * time.Second,
	}
}

// normalizeVersion ensures the version string is in the format "v0.0.XX".
func normalizeVersion(version string) string {
	// Strip leading "v" if present for processing
	v := strings.TrimPrefix(version, "v")

	// If already in format "0.0.XX", add "v" prefix
	if strings.HasPrefix(v, "0.0.") {
		return "v" + v
	}

	// Handle shorthand like "40" or "0040" -> "v0.0.40"
	v = strings.TrimLeft(v, "0")
	if v == "" {
		v = "0"
	}
	return fmt.Sprintf("v0.0.%s", v)
}

// ServiceClient represents a pre-configured client for a specific Slurm domain and version.
// This is provided for backward compatibility with the package-function style API.
// For new code, prefer using slurm.Client or slurmdb.Client directly.
type ServiceClient struct {
	// Endpoint is the full base URL for API calls, e.g. "http://host:6820/slurm/v0.0.40".
	Endpoint string

	// HTTPClient is the HTTP client used for making requests.
	HTTPClient *http.Client

	// AuthToken is the Slurm JWT token for authentication.
	AuthToken string

	// Version is the API version, e.g. "v0.0.40".
	Version string

	// Domain is either "slurm" or "slurmdb".
	Domain Domain
}

// NewSlurmClient creates a ServiceClient for the Slurm (slurmctld) domain.
// Deprecated: Use slurm.NewClient() instead for the new interface-based API.
func NewSlurmClient(cfg Config, version string) (*ServiceClient, error) {
	return newServiceClient(cfg, version, DomainSlurm)
}

// NewSlurmdbClient creates a ServiceClient for the Slurmdb (slurmdbd) domain.
// Deprecated: Use slurmdb.NewClient() instead for the new interface-based API.
func NewSlurmdbClient(cfg Config, version string) (*ServiceClient, error) {
	return newServiceClient(cfg, version, DomainSlurmdb)
}

func newServiceClient(cfg Config, version string, domain Domain) (*ServiceClient, error) {
	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	normalizedVersion := normalizeVersion(version)
	endpoint := fmt.Sprintf("%s/%s/%s", cfg.Endpoint, domain, normalizedVersion)

	httpClient := cfg.HTTPClient
	if httpClient == nil {
		httpClient = defaultHTTPClient(cfg.TLSConfig)
	}

	return &ServiceClient{
		Endpoint:   endpoint,
		HTTPClient: httpClient,
		AuthToken:  cfg.AuthToken,
		Version:    normalizedVersion,
		Domain:     domain,
	}, nil
}

// ServiceURL constructs a full URL for the given resource path parts.
// Example: c.ServiceURL("job", "123") -> "http://host:6820/slurm/v0.0.40/job/123"
func (c *ServiceClient) ServiceURL(parts ...string) string {
	return c.Endpoint + "/" + strings.Join(parts, "/")
}
