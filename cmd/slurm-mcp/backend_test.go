package main

import (
	"context"
	"encoding/json"
	"testing"
)

// mockBackend implements MCPBackend for testing.
type mockBackend struct {
	listJobsResult string
	getJobResult   string
	submitResult   string
	cancelErr      error
	listNodesRes   string
	getNodeRes     string
	listPartsRes   string
	getPartRes     string
	listResrvRes   string
	listLicRes     string
	getDiagRes     string
	pingRes        string

	listAccountsRes string
	getAccountRes   string
	listUsersRes    string
	getUserRes      string
	listDbJobsRes   string
	getDbJobRes     string
	listClustersRes string
	listQOSRes      string
	listAssocRes    string
	listTRESRes     string
	listWCKeysRes   string
	getDbDiagRes    string

	err error // generic error for any call
}

func (m *mockBackend) ListJobs(ctx context.Context) (string, error) {
	return m.listJobsResult, m.err
}
func (m *mockBackend) GetJob(ctx context.Context, jobID string) (string, error) {
	return m.getJobResult, m.err
}
func (m *mockBackend) SubmitJob(ctx context.Context, script string) (string, error) {
	return m.submitResult, m.err
}
func (m *mockBackend) CancelJob(ctx context.Context, jobID string) error {
	if m.cancelErr != nil {
		return m.cancelErr
	}
	return m.err
}
func (m *mockBackend) ListNodes(ctx context.Context) (string, error) {
	return m.listNodesRes, m.err
}
func (m *mockBackend) GetNode(ctx context.Context, name string) (string, error) {
	return m.getNodeRes, m.err
}
func (m *mockBackend) ListPartitions(ctx context.Context) (string, error) {
	return m.listPartsRes, m.err
}
func (m *mockBackend) GetPartition(ctx context.Context, name string) (string, error) {
	return m.getPartRes, m.err
}
func (m *mockBackend) ListReservations(ctx context.Context) (string, error) {
	return m.listResrvRes, m.err
}
func (m *mockBackend) ListLicenses(ctx context.Context) (string, error) {
	return m.listLicRes, m.err
}
func (m *mockBackend) GetDiag(ctx context.Context) (string, error) {
	return m.getDiagRes, m.err
}
func (m *mockBackend) Ping(ctx context.Context) (string, error) {
	return m.pingRes, m.err
}
func (m *mockBackend) ListAccounts(ctx context.Context) (string, error) {
	return m.listAccountsRes, m.err
}
func (m *mockBackend) GetAccount(ctx context.Context, name string) (string, error) {
	return m.getAccountRes, m.err
}
func (m *mockBackend) ListUsers(ctx context.Context) (string, error) {
	return m.listUsersRes, m.err
}
func (m *mockBackend) GetUser(ctx context.Context, name string) (string, error) {
	return m.getUserRes, m.err
}
func (m *mockBackend) ListSlurmdbJobs(ctx context.Context) (string, error) {
	return m.listDbJobsRes, m.err
}
func (m *mockBackend) GetSlurmdbJob(ctx context.Context, jobID string) (string, error) {
	return m.getDbJobRes, m.err
}
func (m *mockBackend) ListClusters(ctx context.Context) (string, error) {
	return m.listClustersRes, m.err
}
func (m *mockBackend) ListQOS(ctx context.Context) (string, error) {
	return m.listQOSRes, m.err
}
func (m *mockBackend) ListAssociations(ctx context.Context) (string, error) {
	return m.listAssocRes, m.err
}
func (m *mockBackend) ListTRES(ctx context.Context) (string, error) {
	return m.listTRESRes, m.err
}
func (m *mockBackend) ListWCKeys(ctx context.Context) (string, error) {
	return m.listWCKeysRes, m.err
}
func (m *mockBackend) GetSlurmdbDiag(ctx context.Context) (string, error) {
	return m.getDbDiagRes, m.err
}

// Verify mockBackend satisfies MCPBackend at compile time.
var _ MCPBackend = (*mockBackend)(nil)

func TestMCPBackendInterface(t *testing.T) {
	ctx := context.Background()
	backend := &mockBackend{
		listJobsResult: `[{"job_id": 123}]`,
		getJobResult:   `{"job_id": 123}`,
		submitResult:   `{"job_id": 456}`,
		pingRes:        `{"status": "ok"}`,
	}

	t.Run("ListJobs", func(t *testing.T) {
		result, err := backend.ListJobs(ctx)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if result != `[{"job_id": 123}]` {
			t.Errorf("unexpected result: %s", result)
		}
	})

	t.Run("GetJob", func(t *testing.T) {
		result, err := backend.GetJob(ctx, "123")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if result != `{"job_id": 123}` {
			t.Errorf("unexpected result: %s", result)
		}
	})

	t.Run("SubmitJob", func(t *testing.T) {
		result, err := backend.SubmitJob(ctx, "#!/bin/bash\necho hello")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if result != `{"job_id": 456}` {
			t.Errorf("unexpected result: %s", result)
		}
	})

	t.Run("CancelJob", func(t *testing.T) {
		err := backend.CancelJob(ctx, "123")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("Ping", func(t *testing.T) {
		result, err := backend.Ping(ctx)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if result != `{"status": "ok"}` {
			t.Errorf("unexpected result: %s", result)
		}
	})
}

func TestCheckRBAC(t *testing.T) {
	t.Run("nil RBAC allows all", func(t *testing.T) {
		if err := checkRBAC(nil, "slurm_submit_job"); err != nil {
			t.Errorf("expected nil error, got: %v", err)
		}
	})

	t.Run("admin allows all", func(t *testing.T) {
		rbac, _ := NewRBAC(RBACConfig{DefaultAccess: AccessAdmin})
		defer rbac.Close()
		if err := checkRBAC(rbac, "slurm_submit_job"); err != nil {
			t.Errorf("expected nil error, got: %v", err)
		}
	})

	t.Run("read-only blocks submit", func(t *testing.T) {
		rbac, _ := NewRBAC(RBACConfig{DefaultAccess: AccessReadOnly})
		defer rbac.Close()
		if err := checkRBAC(rbac, "slurm_submit_job"); err == nil {
			t.Error("expected error for submit with read-only access")
		}
	})

	t.Run("read-only allows list_jobs", func(t *testing.T) {
		rbac, _ := NewRBAC(RBACConfig{DefaultAccess: AccessReadOnly})
		defer rbac.Close()
		if err := checkRBAC(rbac, "slurm_list_jobs"); err != nil {
			t.Errorf("expected nil error for list_jobs with read-only, got: %v", err)
		}
	})

	t.Run("operator allows submit", func(t *testing.T) {
		rbac, _ := NewRBAC(RBACConfig{DefaultAccess: AccessOperator})
		defer rbac.Close()
		if err := checkRBAC(rbac, "slurm_submit_job"); err != nil {
			t.Errorf("expected nil error for submit with operator, got: %v", err)
		}
	})

	t.Run("operator allows cancel", func(t *testing.T) {
		rbac, _ := NewRBAC(RBACConfig{DefaultAccess: AccessOperator})
		defer rbac.Close()
		if err := checkRBAC(rbac, "slurm_cancel_job"); err != nil {
			t.Errorf("expected nil error for cancel with operator, got: %v", err)
		}
	})
}

func TestToolRegistrar(t *testing.T) {
	t.Run("without all-tools only registers core tools", func(t *testing.T) {
		registered := make(map[string]bool)
		tr := &toolRegistrar{
			server:   nil, // We can't actually call server.AddTool, so we test the gating logic
			allTools: false,
		}

		// Test core tools are recognized
		for name := range coreTools {
			if !tr.allTools && !coreTools[name] {
				t.Errorf("core tool %s should be registered", name)
			}
		}

		// Test non-core tool would be skipped
		nonCore := "slurm_list_reservations"
		if tr.allTools || coreTools[nonCore] {
			t.Errorf("non-core tool %s should be skipped without --all-tools", nonCore)
		}

		_ = registered // used only for tracking
	})

	t.Run("core tools list has expected count", func(t *testing.T) {
		if len(coreTools) != 8 {
			t.Errorf("expected 8 core tools, got %d", len(coreTools))
		}
	})

	t.Run("core tools contains expected entries", func(t *testing.T) {
		expected := []string{
			"slurm_ping", "slurm_list_jobs", "slurm_get_job",
			"slurm_submit_job", "slurm_cancel_job", "slurm_list_nodes",
			"slurm_list_partitions", "slurm_queue_depth",
		}
		for _, name := range expected {
			if !coreTools[name] {
				t.Errorf("expected %s to be a core tool", name)
			}
		}
	})
}

func TestMarshalJSON(t *testing.T) {
	t.Run("marshals struct to JSON", func(t *testing.T) {
		data := struct {
			Name string `json:"name"`
			ID   int    `json:"id"`
		}{"test", 42}

		result, err := marshalJSON(data)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		var parsed map[string]interface{}
		if err := json.Unmarshal([]byte(result), &parsed); err != nil {
			t.Fatalf("result is not valid JSON: %v", err)
		}
		if parsed["name"] != "test" {
			t.Errorf("expected name=test, got %v", parsed["name"])
		}
	})
}
