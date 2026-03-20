package slurmtest_test

import (
	"context"
	"testing"

	slurmclient "github.com/supergate-hub/slurm-client"
	"github.com/supergate-hub/slurm-client/slurmtest"
)

func TestMockServer_EndToEnd(t *testing.T) {
	mock := slurmtest.NewServer()
	defer mock.Close()

	mock.AddJob(slurmtest.MockJob{ID: 100, Name: "job-alpha"})
	mock.AddJob(slurmtest.MockJob{ID: 200, Name: "job-beta"})
	mock.AddNode(slurmtest.MockNode{Name: "node01"})
	mock.AddPartition(slurmtest.MockPartition{Name: "gpu"})
	mock.AddAccount(slurmtest.MockAccount{Name: "research"})

	ctx := context.Background()
	client, err := slurmclient.NewClient(ctx, slurmclient.AuthOpts{
		Endpoint:  mock.URL(),
		AuthToken: "mock-token",
		Version:   slurmclient.V0040,
	})
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}

	// Jobs
	jobs, err := client.Slurm.Jobs().List(ctx)
	if err != nil {
		t.Fatalf("Jobs.List: %v", err)
	}
	if len(jobs) != 2 {
		t.Errorf("len(jobs) = %d, want 2", len(jobs))
	}

	// Nodes
	nodes, err := client.Slurm.Nodes().List(ctx)
	if err != nil {
		t.Fatalf("Nodes.List: %v", err)
	}
	if len(nodes) != 1 || *nodes[0].Name != "node01" {
		t.Errorf("unexpected nodes: %+v", nodes)
	}

	// Partitions
	parts, err := client.Slurm.Partitions().List(ctx)
	if err != nil {
		t.Fatalf("Partitions.List: %v", err)
	}
	if len(parts) != 1 || *parts[0].Name != "gpu" {
		t.Errorf("unexpected partitions: %+v", parts)
	}

	// Accounts (SlurmDB)
	accts, err := client.Slurmdb.Accounts().List(ctx)
	if err != nil {
		t.Fatalf("Accounts.List: %v", err)
	}
	if len(accts) != 1 || accts[0].Name != "research" {
		t.Errorf("unexpected accounts: %+v", accts)
	}

	// Submit
	result, err := client.Slurm.Jobs().Submit(ctx, slurmclient.JobSubmitOpts{
		Script: "#!/bin/bash\necho test",
	})
	if err != nil {
		t.Fatalf("Jobs.Submit: %v", err)
	}
	if result.JobID == 0 {
		t.Error("expected non-zero JobID")
	}

	// Ping
	pings, err := client.Slurm.Ping().Get(ctx)
	if err != nil {
		t.Fatalf("Ping.Get: %v", err)
	}
	if len(pings) != 1 {
		t.Errorf("len(pings) = %d, want 1", len(pings))
	}

	// Version discovery
	client2, err := slurmclient.NewClient(ctx, slurmclient.AuthOpts{
		Endpoint:  mock.URL(),
		AuthToken: "mock-token",
		// Version omitted
	})
	if err != nil {
		t.Fatalf("NewClient (auto-detect): %v", err)
	}
	if client2.Version() != slurmclient.V0040 {
		t.Errorf("auto-detected version = %q, want v0.0.40", client2.Version())
	}
}
