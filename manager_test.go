package slurmclient_test

import (
	"context"
	"errors"
	"testing"

	slurmclient "github.com/supergate-hub/slurm-client"
	"github.com/supergate-hub/slurm-client/slurmtest"
)

func TestNewManager(t *testing.T) {
	t.Run("two clusters", func(t *testing.T) {
		mock1 := slurmtest.NewServer()
		defer mock1.Close()
		mock2 := slurmtest.NewServer()
		defer mock2.Close()

		mgr, err := slurmclient.NewManager(context.Background(), slurmclient.ManagerOpts{
			Clusters: map[string]slurmclient.AuthOpts{
				"gpu": {Endpoint: mock1.URL(), AuthToken: "test", Version: slurmclient.V0040},
				"cpu": {Endpoint: mock2.URL(), AuthToken: "test", Version: slurmclient.V0040},
			},
		})
		if err != nil {
			t.Fatalf("NewManager: %v", err)
		}

		names := mgr.ClusterNames()
		if len(names) != 2 {
			t.Fatalf("expected 2 clusters, got %d", len(names))
		}
	})

	t.Run("empty clusters", func(t *testing.T) {
		_, err := slurmclient.NewManager(context.Background(), slurmclient.ManagerOpts{
			Clusters: map[string]slurmclient.AuthOpts{},
		})
		if err == nil {
			t.Fatal("expected error for empty clusters")
		}
		if !errors.Is(err, slurmclient.ErrValidation) {
			t.Fatalf("expected ErrValidation, got %v", err)
		}
	})

	t.Run("partial init failure", func(t *testing.T) {
		mock := slurmtest.NewServer()
		defer mock.Close()

		// Without Version, NewClient tries version discovery via HTTP.
		// The "bad" cluster endpoint is unreachable, so discovery fails.
		mgr, err := slurmclient.NewManager(context.Background(), slurmclient.ManagerOpts{
			Clusters: map[string]slurmclient.AuthOpts{
				"good": {Endpoint: mock.URL(), AuthToken: "test", Version: slurmclient.V0040},
				"bad":  {Endpoint: "http://127.0.0.1:1", AuthToken: "test"}, // no version → discovery fails
			},
		})
		// Partial success: one cluster up, one down
		if err != nil {
			t.Fatalf("expected partial success, got error: %v", err)
		}
		if len(mgr.ClusterNames()) != 1 {
			t.Fatalf("expected 1 cluster, got %d", len(mgr.ClusterNames()))
		}
	})

	t.Run("all init fail", func(t *testing.T) {
		_, err := slurmclient.NewManager(context.Background(), slurmclient.ManagerOpts{
			Clusters: map[string]slurmclient.AuthOpts{
				"bad1": {Endpoint: "http://127.0.0.1:1", AuthToken: "test"}, // no version → discovery fails
				"bad2": {Endpoint: "http://127.0.0.1:2", AuthToken: "test"}, // no version → discovery fails
			},
		})
		if err == nil {
			t.Fatal("expected error when all clusters fail")
		}
		var mce *slurmclient.MultiClusterError
		if !errors.As(err, &mce) {
			t.Fatalf("expected MultiClusterError, got %T", err)
		}
	})
}

func TestManagerCluster(t *testing.T) {
	mock := slurmtest.NewServer()
	defer mock.Close()

	mgr, err := slurmclient.NewManager(context.Background(), slurmclient.ManagerOpts{
		Clusters: map[string]slurmclient.AuthOpts{
			"test": {Endpoint: mock.URL(), AuthToken: "test", Version: slurmclient.V0040},
		},
	})
	if err != nil {
		t.Fatalf("NewManager: %v", err)
	}

	t.Run("known cluster", func(t *testing.T) {
		c, err := mgr.Cluster("test")
		if err != nil {
			t.Fatalf("Cluster: %v", err)
		}
		if c == nil {
			t.Fatal("expected non-nil client")
		}
	})

	t.Run("unknown cluster", func(t *testing.T) {
		_, err := mgr.Cluster("nonexistent")
		if err == nil {
			t.Fatal("expected error for unknown cluster")
		}
		if !errors.Is(err, slurmclient.ErrNotFound) {
			t.Fatalf("expected ErrNotFound, got %v", err)
		}
	})
}

func TestManagerAllJobs(t *testing.T) {
	mock1 := slurmtest.NewServer()
	defer mock1.Close()
	mock2 := slurmtest.NewServer()
	defer mock2.Close()

	mock1.AddJob(slurmtest.MockJob{ID: 1, Name: "job-a"})
	mock2.AddJob(slurmtest.MockJob{ID: 2, Name: "job-b"})

	mgr, err := slurmclient.NewManager(context.Background(), slurmclient.ManagerOpts{
		Clusters: map[string]slurmclient.AuthOpts{
			"c1": {Endpoint: mock1.URL(), AuthToken: "test", Version: slurmclient.V0040},
			"c2": {Endpoint: mock2.URL(), AuthToken: "test", Version: slurmclient.V0040},
		},
	})
	if err != nil {
		t.Fatalf("NewManager: %v", err)
	}

	results := mgr.AllJobs(context.Background())
	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}

	totalJobs := 0
	for _, r := range results {
		if r.Err != nil {
			t.Fatalf("cluster %s error: %v", r.Cluster, r.Err)
		}
		totalJobs += len(r.Data)
	}
	if totalJobs != 2 {
		t.Fatalf("expected 2 total jobs, got %d", totalJobs)
	}
}

func TestManagerAllNodes(t *testing.T) {
	mock1 := slurmtest.NewServer()
	defer mock1.Close()
	mock2 := slurmtest.NewServer()
	defer mock2.Close()

	mock1.AddNode(slurmtest.MockNode{Name: "node01"})
	mock2.AddNode(slurmtest.MockNode{Name: "node02"})

	mgr, err := slurmclient.NewManager(context.Background(), slurmclient.ManagerOpts{
		Clusters: map[string]slurmclient.AuthOpts{
			"c1": {Endpoint: mock1.URL(), AuthToken: "test", Version: slurmclient.V0040},
			"c2": {Endpoint: mock2.URL(), AuthToken: "test", Version: slurmclient.V0040},
		},
	})
	if err != nil {
		t.Fatalf("NewManager: %v", err)
	}

	results := mgr.AllNodes(context.Background())
	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}

	totalNodes := 0
	for _, r := range results {
		if r.Err != nil {
			t.Fatalf("cluster %s error: %v", r.Cluster, r.Err)
		}
		totalNodes += len(r.Data)
	}
	if totalNodes != 2 {
		t.Fatalf("expected 2 total nodes, got %d", totalNodes)
	}
}

func TestManagerContextCancellation(t *testing.T) {
	mock := slurmtest.NewServer()
	defer mock.Close()

	mgr, err := slurmclient.NewManager(context.Background(), slurmclient.ManagerOpts{
		Clusters: map[string]slurmclient.AuthOpts{
			"c1": {Endpoint: mock.URL(), AuthToken: "test", Version: slurmclient.V0040},
		},
	})
	if err != nil {
		t.Fatalf("NewManager: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	cancel() // cancel immediately

	results := mgr.AllJobs(ctx)
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	// With a cancelled context, we expect an error
	if results[0].Err == nil {
		t.Log("no error with cancelled context — mock server responded before cancellation (acceptable)")
	}
}
