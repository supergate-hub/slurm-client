// SPDX-FileCopyrightText: Copyright (C) SchedMD LLC.
// SPDX-License-Identifier: Apache-2.0

package slurmclient

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestCoreClient_URL(t *testing.T) {
	core := &CoreClient{
		BaseURL:    "http://localhost:6820/slurm",
		APIVersion: "v0.0.40",
	}

	tests := []struct {
		path     string
		expected string
	}{
		{"/jobs", "http://localhost:6820/slurm/v0.0.40/jobs"},
		{"/job/123", "http://localhost:6820/slurm/v0.0.40/job/123"},
		{"", "http://localhost:6820/slurm/v0.0.40"},
	}

	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			result := core.URL(tt.path)
			if result != tt.expected {
				t.Errorf("CoreClient.URL(%q) = %q, want %q", tt.path, result, tt.expected)
			}
		})
	}
}

func TestCoreClient_Get(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("Expected GET, got %s", r.Method)
		}
		if r.Header.Get("Authorization") != "Bearer test-token" {
			t.Errorf("Expected Bearer test-token, got %s", r.Header.Get("Authorization"))
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
	}))
	defer server.Close()

	core := &CoreClient{
		HTTPClient: server.Client(),
		BaseURL:    server.URL,
		APIVersion: "v0.0.40",
		AuthToken:  "test-token",
	}

	var result map[string]string
	err := core.Get(context.Background(), "/test", nil, &result)
	if err != nil {
		t.Fatalf("CoreClient.Get() error = %v", err)
	}

	if result["status"] != "ok" {
		t.Errorf("CoreClient.Get() result = %v, want {status: ok}", result)
	}
}

func TestCoreClient_GetWithQuery(t *testing.T) {
	var receivedQuery url.Values
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		receivedQuery = r.URL.Query()
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
	}))
	defer server.Close()

	core := &CoreClient{
		HTTPClient: server.Client(),
		BaseURL:    server.URL,
		APIVersion: "v0.0.40",
		AuthToken:  "test-token",
	}

	query := url.Values{}
	query.Set("update_time", "12345")

	var result map[string]string
	err := core.Get(context.Background(), "/test", query, &result)
	if err != nil {
		t.Fatalf("CoreClient.Get() error = %v", err)
	}

	if receivedQuery.Get("update_time") != "12345" {
		t.Errorf("Expected query param update_time=12345, got %s", receivedQuery.Get("update_time"))
	}
}

func TestCoreClient_Post(t *testing.T) {
	var receivedBody map[string]string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("Expected POST, got %s", r.Method)
		}
		json.NewDecoder(r.Body).Decode(&receivedBody)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"result": "created"})
	}))
	defer server.Close()

	core := &CoreClient{
		HTTPClient: server.Client(),
		BaseURL:    server.URL,
		APIVersion: "v0.0.40",
		AuthToken:  "test-token",
	}

	reqBody := map[string]string{"script": "#!/bin/bash\necho hello"}
	var result map[string]string
	err := core.Post(context.Background(), "/job/submit", reqBody, &result)
	if err != nil {
		t.Fatalf("CoreClient.Post() error = %v", err)
	}

	if result["result"] != "created" {
		t.Errorf("CoreClient.Post() result = %v, want {result: created}", result)
	}
}

func TestCoreClient_Delete(t *testing.T) {
	var receivedPath string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("Expected DELETE, got %s", r.Method)
		}
		receivedPath = r.URL.Path
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	core := &CoreClient{
		HTTPClient: server.Client(),
		BaseURL:    server.URL,
		APIVersion: "v0.0.40",
		AuthToken:  "test-token",
	}

	err := core.Delete(context.Background(), "/job/123", nil)
	if err != nil {
		t.Fatalf("CoreClient.Delete() error = %v", err)
	}

	if receivedPath != "/v0.0.40/job/123" {
		t.Errorf("Expected path /v0.0.40/job/123, got %s", receivedPath)
	}
}

func TestCoreClient_HTTPError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"errors": [{"message": "job not found"}]}`))
	}))
	defer server.Close()

	core := &CoreClient{
		HTTPClient: server.Client(),
		BaseURL:    server.URL,
		APIVersion: "v0.0.40",
		AuthToken:  "test-token",
	}

	var result map[string]string
	err := core.Get(context.Background(), "/job/999", nil, &result)
	if err == nil {
		t.Fatal("CoreClient.Get() expected error for 404 response")
	}

	apiErr, ok := err.(*APIError)
	if !ok {
		t.Fatalf("Expected APIError, got %T", err)
	}
	if apiErr.StatusCode != http.StatusNotFound {
		t.Errorf("APIError.StatusCode = %d, want %d", apiErr.StatusCode, http.StatusNotFound)
	}
}

func TestAPIError_Error(t *testing.T) {
	err := &APIError{
		StatusCode: 404,
		Status:     "404 Not Found",
		Body:       `{"errors": [{"message": "not found"}]}`,
	}

	if err.Error() == "" {
		t.Error("APIError.Error() should return a non-empty string")
	}
}

func TestNewCoreClient(t *testing.T) {
	cfg := Config{
		Endpoint:  "http://localhost:6820",
		AuthToken: "test-token",
	}

	core, err := NewCoreClient(cfg, DomainSlurm, "v0.0.40")
	if err != nil {
		t.Fatalf("NewCoreClient() error = %v", err)
	}

	if core.BaseURL != "http://localhost:6820/slurm" {
		t.Errorf("CoreClient.BaseURL = %q, want %q", core.BaseURL, "http://localhost:6820/slurm")
	}
	if core.APIVersion != "v0.0.40" {
		t.Errorf("CoreClient.APIVersion = %q, want %q", core.APIVersion, "v0.0.40")
	}
}

func TestNewCoreClient_InvalidConfig(t *testing.T) {
	cfg := Config{
		Endpoint:  "",
		AuthToken: "test-token",
	}

	_, err := NewCoreClient(cfg, DomainSlurm, "v0.0.40")
	if err == nil {
		t.Fatal("NewCoreClient() expected error for invalid config")
	}
}
