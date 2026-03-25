// Example: single-cluster SDK usage.
//
// Usage:
//
//	SLURM_ENDPOINT=http://slurmctld:6820 SLURM_TOKEN=jwt go run .
package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"os"

	slurmclient "github.com/supergate-hub/slurm-client"
)

func main() {
	ctx := context.Background()

	client, err := slurmclient.NewClient(ctx, slurmclient.AuthOpts{
		Endpoint:  os.Getenv("SLURM_ENDPOINT"),
		AuthToken: os.Getenv("SLURM_TOKEN"),
		Logger:    slog.Default(),
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Connected (version: %s)\n\n", client.Version())

	// List jobs
	jobs, err := client.Slurm.Jobs().List(ctx)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Jobs: %d\n", len(jobs))
	for _, j := range jobs {
		fmt.Printf("  %d: %s\n", *j.JobId, *j.Name)
	}

	// List nodes
	nodes, err := client.Slurm.Nodes().List(ctx)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("\nNodes: %d\n", len(nodes))
	for _, n := range nodes {
		fmt.Printf("  %s (cpus=%d)\n", *n.Name, *n.Cpus)
	}

	// Submit a job
	result, err := client.Slurm.Jobs().Submit(ctx, slurmclient.JobSubmitOpts{
		Script: "#!/bin/bash\nsrun hostname",
	})
	if err != nil {
		log.Printf("Submit failed: %v", err)
	} else {
		fmt.Printf("\nSubmitted job: %d\n", result.JobID)
	}
}
