package slurmclient

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	v0040 "github.com/supergate-hub/slurm-client/api/v0040"
)

func newTestServer(t *testing.T, handler http.HandlerFunc) *httptest.Server {
	t.Helper()
	return httptest.NewServer(handler)
}

func newTestClient(t *testing.T, serverURL string) *Client {
	t.Helper()
	c, err := NewClient(context.Background(), AuthOpts{
		Endpoint:  serverURL,
		AuthToken: "test-token",
		Version:   V0040,
	})
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}
	return c
}

// --- NewClient tests ---

func TestNewClient_Valid(t *testing.T) {
	srv := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {})
	defer srv.Close()

	c, err := NewClient(context.Background(), AuthOpts{
		Endpoint:  srv.URL,
		AuthToken: "token",
		Version:   V0040,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if c.Version() != V0040 {
		t.Errorf("version = %q, want %q", c.Version(), V0040)
	}
	if c.Slurm == nil {
		t.Error("Slurm proxy is nil")
	}
	if c.Slurmdb == nil {
		t.Error("Slurmdb proxy is nil")
	}
}

func TestNewClient_MissingEndpoint(t *testing.T) {
	_, err := NewClient(context.Background(), AuthOpts{
		AuthToken: "token",
		Version:   V0040,
	})
	if err == nil {
		t.Fatal("expected error for missing endpoint")
	}
}

func TestNewClient_MissingToken(t *testing.T) {
	_, err := NewClient(context.Background(), AuthOpts{
		Endpoint: "http://localhost:6820",
		Version:  V0040,
	})
	if err == nil {
		t.Fatal("expected error for missing token")
	}
}

func TestNewClient_InvalidVersion(t *testing.T) {
	_, err := NewClient(context.Background(), AuthOpts{
		Endpoint:  "http://localhost:6820",
		AuthToken: "token",
		Version:   "v9.9.99",
	})
	if err == nil {
		t.Fatal("expected error for invalid version")
	}
}

// --- Transport / Auth tests ---

func TestTransport_BearerToken(t *testing.T) {
	var gotAuth string
	srv := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		gotAuth = r.Header.Get("Authorization")
		w.WriteHeader(200)
		json.NewEncoder(w).Encode(map[string]any{"pings": []any{}})
	})
	defer srv.Close()

	c := newTestClient(t, srv.URL)
	_, _ = c.Slurm.Ping().Get(context.Background())

	if gotAuth != "Bearer test-token" {
		t.Errorf("Authorization = %q, want %q", gotAuth, "Bearer test-token")
	}
}

// --- Jobs tests ---

func TestJobs_List(t *testing.T) {
	srv := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/slurm/v0.0.40/jobs" {
			t.Errorf("path = %q, want /slurm/v0.0.40/jobs", r.URL.Path)
		}
		if r.Method != http.MethodGet {
			t.Errorf("method = %q, want GET", r.Method)
		}
		name := "test-job"
		json.NewEncoder(w).Encode(v0040.V0040OpenapiJobInfoResp{
			Jobs: []v0040.V0040JobInfo{{Name: &name}},
		})
	})
	defer srv.Close()

	c := newTestClient(t, srv.URL)
	jobs, err := c.Slurm.Jobs().List(context.Background())
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(jobs) != 1 {
		t.Fatalf("len(jobs) = %d, want 1", len(jobs))
	}
	if *jobs[0].Name != "test-job" {
		t.Errorf("name = %q, want %q", *jobs[0].Name, "test-job")
	}
}

func TestJobs_Get(t *testing.T) {
	srv := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/slurm/v0.0.40/job/123" {
			t.Errorf("path = %q, want /slurm/v0.0.40/job/123", r.URL.Path)
		}
		name := "my-job"
		json.NewEncoder(w).Encode(v0040.V0040OpenapiJobInfoResp{
			Jobs: []v0040.V0040JobInfo{{Name: &name}},
		})
	})
	defer srv.Close()

	c := newTestClient(t, srv.URL)
	job, err := c.Slurm.Jobs().Get(context.Background(), "123")
	if err != nil {
		t.Fatalf("Get: %v", err)
	}
	if *job.Name != "my-job" {
		t.Errorf("name = %q, want %q", *job.Name, "my-job")
	}
}

func TestJobs_Get_NotFound(t *testing.T) {
	srv := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(404)
	})
	defer srv.Close()

	c := newTestClient(t, srv.URL)
	_, err := c.Slurm.Jobs().Get(context.Background(), "999")
	if err != ErrNotFound {
		t.Errorf("err = %v, want ErrNotFound", err)
	}
}

func TestJobs_Submit(t *testing.T) {
	srv := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("method = %q, want POST", r.Method)
		}
		if r.URL.Path != "/slurm/v0.0.40/job/submit" {
			t.Errorf("path = %q, want /slurm/v0.0.40/job/submit", r.URL.Path)
		}
		var body v0040.V0040JobSubmitReq
		json.NewDecoder(r.Body).Decode(&body)
		if body.Script == nil || *body.Script != "#!/bin/bash\necho hello" {
			t.Errorf("script = %v", body.Script)
		}
		jobID := int32(42)
		json.NewEncoder(w).Encode(v0040.V0040OpenapiJobSubmitResponse{
			JobId: &jobID,
		})
	})
	defer srv.Close()

	c := newTestClient(t, srv.URL)
	result, err := c.Slurm.Jobs().Submit(context.Background(), JobSubmitOpts{
		Script: "#!/bin/bash\necho hello",
	})
	if err != nil {
		t.Fatalf("Submit: %v", err)
	}
	if result.JobID != 42 {
		t.Errorf("JobID = %d, want 42", result.JobID)
	}
}

func TestJobs_Submit_EmptyScript(t *testing.T) {
	srv := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		t.Fatal("server should not be called")
	})
	defer srv.Close()

	c := newTestClient(t, srv.URL)
	_, err := c.Slurm.Jobs().Submit(context.Background(), JobSubmitOpts{})
	if err == nil {
		t.Fatal("expected validation error")
	}
}

func TestJobs_Delete(t *testing.T) {
	srv := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("method = %q, want DELETE", r.Method)
		}
		w.WriteHeader(200)
	})
	defer srv.Close()

	c := newTestClient(t, srv.URL)
	err := c.Slurm.Jobs().Delete(context.Background(), "123")
	if err != nil {
		t.Fatalf("Delete: %v", err)
	}
}

// --- Nodes tests ---

func TestNodes_List(t *testing.T) {
	srv := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/slurm/v0.0.40/nodes" {
			t.Errorf("path = %q", r.URL.Path)
		}
		name := "node01"
		json.NewEncoder(w).Encode(v0040.V0040OpenapiNodesResp{
			Nodes: []v0040.V0040Node{{Name: &name}},
		})
	})
	defer srv.Close()

	c := newTestClient(t, srv.URL)
	nodes, err := c.Slurm.Nodes().List(context.Background())
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(nodes) != 1 || *nodes[0].Name != "node01" {
		t.Errorf("unexpected nodes: %v", nodes)
	}
}

// --- Partitions tests ---

func TestPartitions_List(t *testing.T) {
	srv := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/slurm/v0.0.40/partitions" {
			t.Errorf("path = %q", r.URL.Path)
		}
		name := "default"
		json.NewEncoder(w).Encode(v0040.V0040OpenapiPartitionResp{
			Partitions: []v0040.V0040PartitionInfo{{Name: &name}},
		})
	})
	defer srv.Close()

	c := newTestClient(t, srv.URL)
	parts, err := c.Slurm.Partitions().List(context.Background())
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(parts) != 1 || *parts[0].Name != "default" {
		t.Errorf("unexpected partitions: %v", parts)
	}
}

// --- SlurmDB Accounts tests ---

func TestAccounts_List(t *testing.T) {
	srv := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/slurmdb/v0.0.40/accounts" {
			t.Errorf("path = %q, want /slurmdb/v0.0.40/accounts", r.URL.Path)
		}
		json.NewEncoder(w).Encode(v0040.V0040OpenapiAccountsResp{
			Accounts: []v0040.V0040Account{{Name: "root"}},
		})
	})
	defer srv.Close()

	c := newTestClient(t, srv.URL)
	accts, err := c.Slurmdb.Accounts().List(context.Background())
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(accts) != 1 || accts[0].Name != "root" {
		t.Errorf("unexpected accounts: %v", accts)
	}
}

// --- Error handling tests ---

func TestTransport_Unauthorized(t *testing.T) {
	srv := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(401)
	})
	defer srv.Close()

	c := newTestClient(t, srv.URL)
	_, err := c.Slurm.Jobs().List(context.Background())
	if err != ErrUnauthorized {
		t.Errorf("err = %v, want ErrUnauthorized", err)
	}
}

func TestTransport_ServerError(t *testing.T) {
	srv := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(500)
		json.NewEncoder(w).Encode(map[string]any{
			"errors": []map[string]string{{"error": "internal failure"}},
		})
	})
	defer srv.Close()

	c := newTestClient(t, srv.URL)
	_, err := c.Slurm.Jobs().List(context.Background())
	if err == nil {
		t.Fatal("expected error")
	}
	se, ok := err.(*SlurmError)
	if !ok {
		t.Fatalf("err type = %T, want *SlurmError", err)
	}
	if se.StatusCode != 500 {
		t.Errorf("StatusCode = %d, want 500", se.StatusCode)
	}
	if len(se.Errors) == 0 || se.Errors[0] != "internal failure" {
		t.Errorf("Errors = %v", se.Errors)
	}
}

// --- Proxy caching test ---

func TestProxy_ReturnsCachedService(t *testing.T) {
	srv := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {})
	defer srv.Close()

	c := newTestClient(t, srv.URL)
	j1 := c.Slurm.Jobs()
	j2 := c.Slurm.Jobs()
	if j1 != j2 {
		t.Error("Jobs() returned different instances — expected cached")
	}

	n1 := c.Slurm.Nodes()
	n2 := c.Slurm.Nodes()
	if n1 != n2 {
		t.Error("Nodes() returned different instances")
	}

	a1 := c.Slurmdb.Accounts()
	a2 := c.Slurmdb.Accounts()
	if a1 != a2 {
		t.Error("Accounts() returned different instances")
	}
}

// --- Version tests ---

func TestVersion_IsValid(t *testing.T) {
	tests := []struct {
		v    Version
		want bool
	}{
		{V0039, true},
		{V0040, true},
		{V0044, true},
		{"v9.9.9", false},
		{"", false},
	}
	for _, tt := range tests {
		if got := tt.v.IsValid(); got != tt.want {
			t.Errorf("Version(%q).IsValid() = %v, want %v", tt.v, got, tt.want)
		}
	}
}

// --- Version discovery test ---

func TestNewClient_VersionDiscovery(t *testing.T) {
	srv := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/openapi/v3" {
			json.NewEncoder(w).Encode(map[string]any{
				"paths": map[string]any{
					"/slurm/v0.0.40/jobs":    nil,
					"/slurm/v0.0.40/nodes":   nil,
					"/slurmdb/v0.0.40/users": nil,
				},
			})
			return
		}
		w.WriteHeader(200)
	})
	defer srv.Close()

	c, err := NewClient(context.Background(), AuthOpts{
		Endpoint:  srv.URL,
		AuthToken: "token",
		// Version omitted — auto-detect
	})
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}
	if c.Version() != V0040 {
		t.Errorf("auto-detected version = %q, want %q", c.Version(), V0040)
	}
}

// --- Multi-version URL tests ---

func TestMultiVersion_URLPrefix(t *testing.T) {
	versions := []Version{V0039, V0040, V0041, V0042, V0043, V0044}

	for _, v := range versions {
		t.Run(string(v), func(t *testing.T) {
			var gotPath string
			srv := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
				gotPath = r.URL.Path
				json.NewEncoder(w).Encode(map[string]any{"jobs": []any{}})
			})
			defer srv.Close()

			c, err := NewClient(context.Background(), AuthOpts{
				Endpoint:  srv.URL,
				AuthToken: "token",
				Version:   v,
			})
			if err != nil {
				t.Fatalf("NewClient(%s): %v", v, err)
			}

			_, _ = c.Slurm.Jobs().List(context.Background())

			wantPrefix := "/slurm/" + string(v) + "/jobs"
			if gotPath != wantPrefix {
				t.Errorf("path = %q, want %q", gotPath, wantPrefix)
			}
		})
	}
}

func TestMultiVersion_SlurmdbURLPrefix(t *testing.T) {
	var gotPath string
	srv := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Path
		json.NewEncoder(w).Encode(map[string]any{"accounts": []any{}})
	})
	defer srv.Close()

	c, err := NewClient(context.Background(), AuthOpts{
		Endpoint:  srv.URL,
		AuthToken: "token",
		Version:   V0041,
	})
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}

	_, _ = c.Slurmdb.Accounts().List(context.Background())

	want := "/slurmdb/v0.0.41/accounts"
	if gotPath != want {
		t.Errorf("path = %q, want %q", gotPath, want)
	}
}

// --- Retry test ---

func TestTransport_Retry(t *testing.T) {
	attempts := 0
	srv := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		attempts++
		if attempts < 3 {
			w.WriteHeader(503)
			return
		}
		json.NewEncoder(w).Encode(map[string]any{"pings": []any{}})
	})
	defer srv.Close()

	c, err := NewClient(context.Background(), AuthOpts{
		Endpoint:  srv.URL,
		AuthToken: "token",
		Version:   V0040,
	})
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}

	_, err = c.Slurm.Ping().Get(context.Background())
	if err != nil {
		t.Fatalf("expected success after retries, got: %v", err)
	}
	if attempts != 3 {
		t.Errorf("attempts = %d, want 3", attempts)
	}
}

func TestTransport_RetryExhausted(t *testing.T) {
	srv := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(503)
	})
	defer srv.Close()

	c, err := NewClient(context.Background(), AuthOpts{
		Endpoint:       srv.URL,
		AuthToken:      "token",
		Version:        V0040,
		MaxRetries:     2,
		RetryBaseDelay: 1, // 1ns for fast test
	})
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}

	_, err = c.Slurm.Ping().Get(context.Background())
	if err == nil {
		t.Fatal("expected RetryExhaustedError")
	}
	var re *RetryExhaustedError
	if !errors.As(err, &re) {
		t.Fatalf("err type = %T, want *RetryExhaustedError", err)
	}
}

// --- Unix socket test (config only) ---

func TestNewClient_UnixSocket(t *testing.T) {
	// Can't test actual unix socket without a real socket file,
	// but verify that the client creates successfully with the option.
	_, err := NewClient(context.Background(), AuthOpts{
		UnixSocket: "/tmp/nonexistent.sock",
		AuthToken:  "token",
		Version:    V0040,
	})
	if err != nil {
		t.Fatalf("NewClient with UnixSocket: %v", err)
	}
}

// --- Logging test ---

func TestTransport_LoggingDoesNotPanic(t *testing.T) {
	srv := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]any{"pings": []any{}})
	})
	defer srv.Close()

	logger := slog.Default()
	c, err := NewClient(context.Background(), AuthOpts{
		Endpoint:  srv.URL,
		AuthToken: "secret-token",
		Version:   V0040,
		Logger:    logger,
	})
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}

	// Should not panic with logger enabled
	_, err = c.Slurm.Ping().Get(context.Background())
	if err != nil {
		t.Fatalf("Ping: %v", err)
	}
}
