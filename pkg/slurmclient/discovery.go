package slurmclient

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"
	"time"
)

// DiscoveredVersions contains the discovered API versions for each domain.
type DiscoveredVersions struct {
	// Slurm contains discovered versions for the slurmctld domain.
	// e.g. []string{"v0.0.39", "v0.0.40", "v0.0.44"}
	Slurm []string

	// Slurmdb contains discovered versions for the slurmdbd domain.
	Slurmdb []string
}

// DiscoverVersions attempts to discover the supported Slurm REST API versions
// from the given endpoint. This function makes HTTP requests to probe the API.
//
// Discovery strategies:
// 1. Try the OpenAPI spec endpoint to extract version information
// 2. Try known version endpoints (ping for slurm, diag for slurmdb)
func DiscoverVersions(ctx context.Context, cfg Config) (*DiscoveredVersions, error) {
	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	// Create HTTP client
	httpClient := cfg.HTTPClient
	if httpClient == nil {
		httpClient = &http.Client{
			Timeout: 10 * time.Second,
		}
	}

	discovered := &DiscoveredVersions{
		Slurm:   []string{},
		Slurmdb: []string{},
	}

	baseEndpoint := strings.TrimSuffix(cfg.Endpoint, "/")

	// Try to discover from OpenAPI spec
	openapiVersions, err := discoverFromOpenAPI(ctx, httpClient, baseEndpoint, cfg.AuthToken)
	if err == nil && openapiVersions != nil {
		discovered.Slurm = openapiVersions.Slurm
		discovered.Slurmdb = openapiVersions.Slurmdb
		return discovered, nil
	}

	// Fall back to probing known versions
	knownVersions := []string{"v0.0.44", "v0.0.43", "v0.0.42", "v0.0.41", "v0.0.40", "v0.0.39"}

	for _, version := range knownVersions {
		// Check slurm domain
		if checkSlurmVersion(ctx, httpClient, baseEndpoint, version, cfg.AuthToken) {
			discovered.Slurm = append(discovered.Slurm, version)
		}

		// Check slurmdb domain
		if checkSlurmdbVersion(ctx, httpClient, baseEndpoint, version, cfg.AuthToken) {
			discovered.Slurmdb = append(discovered.Slurmdb, version)
		}
	}

	return discovered, nil
}

// discoverFromOpenAPI tries to discover versions from the OpenAPI specification.
func discoverFromOpenAPI(ctx context.Context, client *http.Client, baseEndpoint, authToken string) (*DiscoveredVersions, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, baseEndpoint+"/openapi/v3", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+authToken)
	req.Header.Set("Accept", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("openapi endpoint returned status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// Parse OpenAPI spec to extract paths
	var spec map[string]any
	if err := json.Unmarshal(body, &spec); err != nil {
		return nil, err
	}

	discovered := &DiscoveredVersions{
		Slurm:   []string{},
		Slurmdb: []string{},
	}

	paths, ok := spec["paths"].(map[string]any)
	if !ok {
		return nil, fmt.Errorf("invalid openapi spec: missing paths")
	}

	slurmVersions := make(map[string]bool)
	slurmdbVersions := make(map[string]bool)

	// Regex to match version patterns like /slurm/v0.0.40/... or /slurmdb/v0.0.40/...
	slurmPattern := regexp.MustCompile(`^/slurm/(v0\.0\.\d+)/`)
	slurmdbPattern := regexp.MustCompile(`^/slurmdb/(v0\.0\.\d+)/`)

	for path := range paths {
		if matches := slurmPattern.FindStringSubmatch(path); len(matches) > 1 {
			slurmVersions[matches[1]] = true
		}
		if matches := slurmdbPattern.FindStringSubmatch(path); len(matches) > 1 {
			slurmdbVersions[matches[1]] = true
		}
	}

	for v := range slurmVersions {
		discovered.Slurm = append(discovered.Slurm, v)
	}
	for v := range slurmdbVersions {
		discovered.Slurmdb = append(discovered.Slurmdb, v)
	}

	return discovered, nil
}

// checkSlurmVersion checks if a specific slurm version is available by trying the ping endpoint.
func checkSlurmVersion(ctx context.Context, client *http.Client, baseEndpoint, version, authToken string) bool {
	url := fmt.Sprintf("%s/slurm/%s/ping", baseEndpoint, version)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return false
	}
	req.Header.Set("Authorization", "Bearer "+authToken)

	resp, err := client.Do(req)
	if err != nil {
		return false
	}
	defer resp.Body.Close()

	return resp.StatusCode >= 200 && resp.StatusCode < 300
}

// checkSlurmdbVersion checks if a specific slurmdb version is available by trying the diag endpoint.
func checkSlurmdbVersion(ctx context.Context, client *http.Client, baseEndpoint, version, authToken string) bool {
	url := fmt.Sprintf("%s/slurmdb/%s/diag", baseEndpoint, version)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return false
	}
	req.Header.Set("Authorization", "Bearer "+authToken)

	resp, err := client.Do(req)
	if err != nil {
		return false
	}
	defer resp.Body.Close()

	return resp.StatusCode >= 200 && resp.StatusCode < 300
}

// LatestVersion returns the latest (highest) version from a list of versions.
func LatestVersion(versions []string) string {
	if len(versions) == 0 {
		return ""
	}

	// Simple string comparison works for version format v0.0.XX
	latest := versions[0]
	for _, v := range versions[1:] {
		if v > latest {
			latest = v
		}
	}
	return latest
}

