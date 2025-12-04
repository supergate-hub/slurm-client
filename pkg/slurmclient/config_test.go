// SPDX-FileCopyrightText: Copyright (C) SchedMD LLC.
// SPDX-License-Identifier: Apache-2.0

package slurmclient

import (
	"net/http"
	"testing"
)

func TestConfig_Validate(t *testing.T) {
	tests := []struct {
		name    string
		config  *Config
		wantErr bool
		errMsg  string
	}{
		{
			name:    "nil config",
			config:  nil,
			wantErr: true,
			errMsg:  "config cannot be nil",
		},
		{
			name: "empty endpoint",
			config: &Config{
				Endpoint:  "",
				AuthToken: "test-token",
			},
			wantErr: true,
			errMsg:  "endpoint cannot be empty",
		},
		{
			name: "empty auth token",
			config: &Config{
				Endpoint:  "http://localhost:6820",
				AuthToken: "",
			},
			wantErr: true,
			errMsg:  "authToken cannot be empty",
		},
		{
			name: "valid config",
			config: &Config{
				Endpoint:  "http://localhost:6820",
				AuthToken: "test-token",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Config.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestNormalizeVersion(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"v0.0.40", "v0.0.40"},
		{"0.0.40", "v0.0.40"},
		{"40", "v0.0.40"},
		{"0040", "v0.0.40"},
		{"v40", "v0.0.40"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := normalizeVersion(tt.input)
			if result != tt.expected {
				t.Errorf("normalizeVersion(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestNewSlurmClient(t *testing.T) {
	cfg := Config{
		Endpoint:  "http://localhost:6820",
		AuthToken: "test-token",
	}

	client, err := NewSlurmClient(cfg, "v0.0.40")
	if err != nil {
		t.Fatalf("NewSlurmClient() error = %v", err)
	}

	if client.Domain != DomainSlurm {
		t.Errorf("NewSlurmClient() domain = %v, want %v", client.Domain, DomainSlurm)
	}
}

func TestNewSlurmdbClient(t *testing.T) {
	cfg := Config{
		Endpoint:  "http://localhost:6820",
		AuthToken: "test-token",
	}

	client, err := NewSlurmdbClient(cfg, "v0.0.40")
	if err != nil {
		t.Fatalf("NewSlurmdbClient() error = %v", err)
	}

	if client.Domain != DomainSlurmdb {
		t.Errorf("NewSlurmdbClient() domain = %v, want %v", client.Domain, DomainSlurmdb)
	}
}

func TestDefaultHTTPClient(t *testing.T) {
	client := defaultHTTPClient(nil)
	if client == nil {
		t.Error("defaultHTTPClient() returned nil")
	}
}

func TestConfigWithCustomHTTPClient(t *testing.T) {
	customClient := &http.Client{}
	cfg := Config{
		Endpoint:   "http://localhost:6820",
		AuthToken:  "test-token",
		HTTPClient: customClient,
	}

	client, err := NewSlurmClient(cfg, "v0.0.40")
	if err != nil {
		t.Fatalf("NewSlurmClient() error = %v", err)
	}

	if client.HTTPClient != customClient {
		t.Error("NewSlurmClient() should use the provided custom HTTP client")
	}
}
