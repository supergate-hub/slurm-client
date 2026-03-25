// Example: multi-cluster Manager usage.
//
// Usage:
//
//	go run . --config clusters.yaml
//
// Or with inline config:
//
//	SLURM_TOKEN=jwt go run .
package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	slurmclient "github.com/supergate-hub/slurm-client"
)

func main() {
	configPath := flag.String("config", "", "Path to clusters.yaml")
	flag.Parse()

	ctx := context.Background()
	var mgr *slurmclient.Manager
	var err error

	if *configPath != "" {
		// Load from YAML config file
		opts, err := slurmclient.ParseConfig(*configPath)
		if err != nil {
			log.Fatal(err)
		}
		mgr, err = slurmclient.NewManager(ctx, *opts)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		// Inline config (demo: same endpoint as two "clusters")
		token := os.Getenv("SLURM_TOKEN")
		endpoint := os.Getenv("SLURM_ENDPOINT")
		if token == "" || endpoint == "" {
			log.Fatal("Use --config or set SLURM_ENDPOINT + SLURM_TOKEN")
		}
		mgr, err = slurmclient.NewManager(ctx, slurmclient.ManagerOpts{
			Clusters: map[string]slurmclient.AuthOpts{
				"cluster-a": {Endpoint: endpoint, AuthToken: token},
				"cluster-b": {Endpoint: endpoint, AuthToken: token},
			},
		})
		if err != nil {
			log.Fatal(err)
		}
	}

	fmt.Printf("Connected clusters: %v\n\n", mgr.ClusterNames())

	// Access a specific cluster
	for _, name := range mgr.ClusterNames() {
		c, _ := mgr.Cluster(name)
		fmt.Printf("[%s] version=%s\n", name, c.Version())
	}

	// Cross-cluster: list all jobs
	fmt.Println("\n--- AllJobs (parallel) ---")
	jobResults := mgr.AllJobs(ctx)
	for _, r := range jobResults {
		if r.Err != nil {
			fmt.Printf("[%s] ERROR: %v\n", r.Cluster, r.Err)
		} else {
			fmt.Printf("[%s] %d jobs\n", r.Cluster, len(r.Data))
		}
	}

	// Cross-cluster: list all nodes
	fmt.Println("\n--- AllNodes (parallel) ---")
	nodeResults := mgr.AllNodes(ctx)
	for _, r := range nodeResults {
		if r.Err != nil {
			fmt.Printf("[%s] ERROR: %v\n", r.Cluster, r.Err)
		} else {
			fmt.Printf("[%s] %d nodes\n", r.Cluster, len(r.Data))
			for _, n := range r.Data {
				fmt.Printf("  └─ %s\n", *n.Name)
			}
		}
	}
}
