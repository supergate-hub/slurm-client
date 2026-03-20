// Package slurmclient provides a Go SDK for the Slurm REST API.
//
// The SDK uses an openstacksdk-inspired proxy pattern that enables
// IDE-friendly auto-completion chains:
//
//	client, err := slurmclient.NewClient(ctx, slurmclient.AuthOpts{
//	    Endpoint:  "http://slurmctld:6820",
//	    AuthToken: "jwt-token",
//	    Version:   slurmclient.V0040,
//	})
//
//	// Auto-complete chain: client.Slurm.Jobs().List()
//	jobs, err := client.Slurm.Jobs().List(ctx)
//	node, err := client.Slurm.Nodes().Get(ctx, "node01")
//	accounts, err := client.Slurmdb.Accounts().List(ctx)
//
// # Version Compatibility
//
// The SDK supports Slurm REST API versions v0.0.39 through v0.0.44.
// Response parsing uses v0.0.40 type definitions, which are JSON-compatible
// with all supported versions:
//
//   - v0.0.40 to v0.0.44: fully compatible (identical JSON field names)
//   - v0.0.39: compatible for core fields; metadata fields like
//     last_backfill and last_update may be absent (zero-valued)
//
// If Version is omitted from AuthOpts, the SDK auto-detects the best
// available version by querying slurmrestd's /openapi endpoint.
//
// # Testing
//
// The slurmtest sub-package provides an in-memory mock slurmrestd server:
//
//	mock := slurmtest.NewServer()
//	defer mock.Close()
//	mock.AddJob(slurmtest.MockJob{ID: 123, Name: "test"})
//
//	client, _ := slurmclient.NewClient(ctx, slurmclient.AuthOpts{
//	    Endpoint:  mock.URL(),
//	    AuthToken: "test",
//	    Version:   slurmclient.V0040,
//	})
package slurmclient
