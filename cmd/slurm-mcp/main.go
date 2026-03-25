// slurm-mcp is an MCP (Model Context Protocol) server for Slurm.
// It enables AI agents to manage Slurm clusters via MCP tools.
//
// Single-cluster mode (backward compatible):
//
//	SLURM_ENDPOINT=http://slurmctld:6820 SLURM_TOKEN=jwt slurm-mcp
//
// Multi-cluster mode:
//
//	slurm-mcp --config clusters.yaml
package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"

	slurmclient "github.com/supergate-hub/slurm-client"
)

// clusterResolver resolves a cluster name from an MCP request to a *Client.
// In single-cluster mode, the cluster parameter is ignored.
// In multi-cluster mode, it routes to the appropriate cluster.
type clusterResolver struct {
	single  *slurmclient.Client  // non-nil in single-cluster mode
	manager *slurmclient.Manager // non-nil in multi-cluster mode
}

func (cr *clusterResolver) resolve(req mcp.CallToolRequest) (*slurmclient.Client, error) {
	if cr.single != nil {
		return cr.single, nil
	}
	var cluster string
	if args, ok := req.Params.Arguments.(map[string]any); ok {
		cluster, _ = args["cluster"].(string)
	}
	if cluster == "" {
		names := cr.manager.ClusterNames()
		if len(names) == 1 {
			return cr.manager.Cluster(names[0])
		}
		return nil, fmt.Errorf("cluster parameter is required in multi-cluster mode (available: %v)", names)
	}
	return cr.manager.Cluster(cluster)
}

func (cr *clusterResolver) isMultiCluster() bool {
	return cr.manager != nil
}

func main() {
	configPath := flag.String("config", "", "Path to clusters.yaml for multi-cluster mode")
	flag.Parse()

	ctx := context.Background()
	var resolver clusterResolver
	var mgr *slurmclient.Manager

	if *configPath != "" {
		// Multi-cluster mode
		opts, err := slurmclient.ParseConfig(*configPath)
		if err != nil {
			log.Fatalf("Failed to parse config: %v", err)
		}
		mgr, err = slurmclient.NewManager(ctx, *opts)
		if err != nil {
			log.Fatalf("Failed to create manager: %v", err)
		}
		resolver.manager = mgr
	} else {
		// Single-cluster mode (backward compatible)
		endpoint := os.Getenv("SLURM_ENDPOINT")
		token := os.Getenv("SLURM_TOKEN")
		version := os.Getenv("SLURM_VERSION")
		socket := os.Getenv("SLURM_UNIX_SOCKET")

		if endpoint == "" && socket == "" {
			log.Fatal("SLURM_ENDPOINT or SLURM_UNIX_SOCKET is required (or use --config for multi-cluster)")
		}
		if token == "" {
			log.Fatal("SLURM_TOKEN is required")
		}

		opts := slurmclient.AuthOpts{
			Endpoint:   endpoint,
			AuthToken:  token,
			UnixSocket: socket,
		}
		if version != "" {
			opts.Version = slurmclient.Version(version)
		}

		client, err := slurmclient.NewClient(ctx, opts)
		if err != nil {
			log.Fatalf("Failed to create Slurm client: %v", err)
		}
		resolver.single = client
	}

	s := server.NewMCPServer(
		"slurm-mcp",
		"3.0.0",
		server.WithToolCapabilities(false),
		server.WithResourceCapabilities(true, false),
		server.WithPromptCapabilities(false),
	)

	registerSlurmTools(s, &resolver)
	registerSlurmdbTools(s, &resolver)
	registerResources(s, &resolver, mgr)
	registerPrompts(s, &resolver, mgr)

	if resolver.isMultiCluster() {
		registerMultiClusterTools(s, mgr)
	}

	if err := server.ServeStdio(s); err != nil {
		fmt.Fprintf(os.Stderr, "Server error: %v\n", err)
		os.Exit(1)
	}
}
