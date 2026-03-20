// slurm-mcp is an MCP (Model Context Protocol) server for Slurm.
// It enables AI agents to manage Slurm clusters via MCP tools.
//
// Usage:
//
//	SLURM_ENDPOINT=http://slurmctld:6820 SLURM_TOKEN=jwt slurm-mcp
//
// Claude Desktop configuration (claude_desktop_config.json):
//
//	{
//	  "mcpServers": {
//	    "slurm": {
//	      "command": "slurm-mcp",
//	      "env": {
//	        "SLURM_ENDPOINT": "http://slurmctld:6820",
//	        "SLURM_TOKEN": "your-jwt-token"
//	      }
//	    }
//	  }
//	}
package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/mark3labs/mcp-go/server"

	slurmclient "github.com/supergate-hub/slurm-client"
)

func main() {
	endpoint := os.Getenv("SLURM_ENDPOINT")
	token := os.Getenv("SLURM_TOKEN")
	version := os.Getenv("SLURM_VERSION")
	socket := os.Getenv("SLURM_UNIX_SOCKET")

	if endpoint == "" && socket == "" {
		log.Fatal("SLURM_ENDPOINT or SLURM_UNIX_SOCKET is required")
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

	client, err := slurmclient.NewClient(context.Background(), opts)
	if err != nil {
		log.Fatalf("Failed to create Slurm client: %v", err)
	}

	s := server.NewMCPServer(
		"slurm-mcp",
		"2.0.0",
		server.WithToolCapabilities(false),
	)

	registerSlurmTools(s, client)
	registerSlurmdbTools(s, client)

	if err := server.ServeStdio(s); err != nil {
		fmt.Fprintf(os.Stderr, "Server error: %v\n", err)
		os.Exit(1)
	}
}
