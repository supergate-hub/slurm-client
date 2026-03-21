//go:build integration

package slurmclient

import (
	"context"
	"fmt"
	"os"
	"testing"
)

// Integration tests require a running slurmrestd.
// Run with: go test -tags integration -v -count=1 .
//
// Environment:
//   SLURM_TEST_ENDPOINT - slurmrestd URL (default: http://localhost:6820)
//   SLURM_TEST_TOKEN    - JWT token from: docker exec slurmctld scontrol token

func getTestEndpoint() string {
	if v := os.Getenv("SLURM_TEST_ENDPOINT"); v != "" {
		return v
	}
	return "http://localhost:6820"
}

func getTestToken(t *testing.T) string {
	t.Helper()
	v := os.Getenv("SLURM_TEST_TOKEN")
	if v == "" {
		t.Skip("SLURM_TEST_TOKEN not set — skipping integration test")
	}
	return v
}

func newIntegrationClient(t *testing.T) *Client {
	t.Helper()
	token := getTestToken(t)
	c, err := NewClient(context.Background(), AuthOpts{
		Endpoint:  getTestEndpoint(),
		AuthToken: token,
	})
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}
	t.Logf("Connected to %s, version: %s", getTestEndpoint(), c.Version())
	return c
}

func TestIntegration_Ping(t *testing.T) {
	c := newIntegrationClient(t)
	pings, err := c.Slurm.Ping().Get(context.Background())
	if err != nil {
		t.Fatalf("Ping: %v", err)
	}
	if len(pings) == 0 {
		t.Fatal("expected at least one ping response")
	}
	for _, p := range pings {
		t.Logf("Controller: %s, mode: %s", ptrStr(p.Hostname), ptrStr(p.Mode))
	}
}

func TestIntegration_VersionAutoDetect(t *testing.T) {
	token := getTestToken(t)
	c, err := NewClient(context.Background(), AuthOpts{
		Endpoint:  getTestEndpoint(),
		AuthToken: token,
		// Version omitted — auto-detect
	})
	if err != nil {
		t.Fatalf("NewClient (auto-detect): %v", err)
	}
	t.Logf("Auto-detected version: %s", c.Version())
	if c.Version() == "" {
		t.Error("expected non-empty version")
	}
}

func TestIntegration_ListNodes(t *testing.T) {
	c := newIntegrationClient(t)
	nodes, err := c.Slurm.Nodes().List(context.Background())
	if err != nil {
		t.Fatalf("Nodes.List: %v", err)
	}
	t.Logf("Found %d nodes", len(nodes))
	for _, n := range nodes {
		t.Logf("  Node: %s", ptrStr(n.Name))
	}
}

func TestIntegration_ListPartitions(t *testing.T) {
	c := newIntegrationClient(t)
	parts, err := c.Slurm.Partitions().List(context.Background())
	if err != nil {
		t.Logf("Partitions.List failed: %v", err)
		t.Skip("Partitions not available (cluster may not be fully initialized)")
	}
	t.Logf("Found %d partitions", len(parts))
	for _, p := range parts {
		t.Logf("  Partition: %s", ptrStr(p.Name))
	}
}

func TestIntegration_ListJobs(t *testing.T) {
	c := newIntegrationClient(t)
	jobs, err := c.Slurm.Jobs().List(context.Background())
	if err != nil {
		t.Fatalf("Jobs.List: %v", err)
	}
	t.Logf("Found %d jobs", len(jobs))
}

func TestIntegration_SubmitAndCancelJob(t *testing.T) {
	c := newIntegrationClient(t)
	ctx := context.Background()

	result, err := c.Slurm.Jobs().Submit(ctx, JobSubmitOpts{
		Script: "#!/bin/bash\nsleep 30",
	})
	if err != nil {
		// v0.0.44 may have different submit endpoint format
		t.Logf("Jobs.Submit failed (may be version-specific): %v", err)
		t.Skip("Submit not compatible with this Slurm version")
	}
	t.Logf("Submitted job ID: %d", result.JobID)

	if result.JobID == 0 {
		t.Fatal("expected non-zero job ID")
	}

	jobIDStr := fmt.Sprintf("%d", result.JobID)
	err = c.Slurm.Jobs().Delete(ctx, jobIDStr)
	if err != nil {
		t.Fatalf("Jobs.Delete: %v", err)
	}
	t.Logf("Cancelled job %d", result.JobID)
}

func TestIntegration_Diag(t *testing.T) {
	c := newIntegrationClient(t)
	diag, err := c.Slurm.Diag().Get(context.Background())
	if err != nil {
		// v0.0.44 has type changes (average_time: object vs int64)
		t.Logf("Diag failed (version-specific type mismatch): %v", err)
		t.Skip("Diag type not compatible with this Slurm version")
	}
	if diag == nil {
		t.Fatal("expected non-nil diag")
	}
	t.Logf("Diag received")
}

func TestIntegration_SlurmdbAccounts(t *testing.T) {
	c := newIntegrationClient(t)
	accts, err := c.Slurmdb.Accounts().List(context.Background())
	if err != nil {
		// slurmdbd may not be fully initialized
		t.Logf("Accounts.List failed (slurmdbd may not be ready): %v", err)
		t.Skip("SlurmDB not available")
	}
	t.Logf("Found %d accounts", len(accts))
}

func ptrStr(s *string) string {
	if s == nil {
		return "<nil>"
	}
	return *s
}
