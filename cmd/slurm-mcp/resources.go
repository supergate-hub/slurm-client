package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"

	slurmclient "github.com/supergate-hub/slurm-client"
)

func registerResources(s *server.MCPServer, cr *clusterResolver, mgr *slurmclient.Manager) {
	if mgr != nil {
		// Multi-cluster: expose each cluster's data as a resource
		for _, name := range mgr.ClusterNames() {
			clusterName := name
			c, _ := mgr.Cluster(clusterName)

			s.AddResource(mcp.Resource{
				URI:         fmt.Sprintf("slurm://%s/nodes", clusterName),
				Name:        fmt.Sprintf("Nodes (%s)", clusterName),
				Description: fmt.Sprintf("All nodes in cluster %s with their states and resources", clusterName),
				MIMEType:    "application/json",
			}, makeNodeResourceHandler(c))

			s.AddResource(mcp.Resource{
				URI:         fmt.Sprintf("slurm://%s/jobs", clusterName),
				Name:        fmt.Sprintf("Jobs (%s)", clusterName),
				Description: fmt.Sprintf("All active jobs in cluster %s", clusterName),
				MIMEType:    "application/json",
			}, makeJobResourceHandler(c))

			s.AddResource(mcp.Resource{
				URI:         fmt.Sprintf("slurm://%s/partitions", clusterName),
				Name:        fmt.Sprintf("Partitions (%s)", clusterName),
				Description: fmt.Sprintf("All partitions in cluster %s with their configurations", clusterName),
				MIMEType:    "application/json",
			}, makePartitionResourceHandler(c))
		}
	} else if cr.single != nil {
		// Single-cluster mode
		s.AddResource(mcp.Resource{
			URI:         "slurm://default/nodes",
			Name:        "Nodes",
			Description: "All nodes in the Slurm cluster with their states and resources",
			MIMEType:    "application/json",
		}, makeNodeResourceHandler(cr.single))

		s.AddResource(mcp.Resource{
			URI:         "slurm://default/jobs",
			Name:        "Jobs",
			Description: "All active jobs in the Slurm cluster",
			MIMEType:    "application/json",
		}, makeJobResourceHandler(cr.single))

		s.AddResource(mcp.Resource{
			URI:         "slurm://default/partitions",
			Name:        "Partitions",
			Description: "All partitions in the Slurm cluster with their configurations",
			MIMEType:    "application/json",
		}, makePartitionResourceHandler(cr.single))
	}
}

func makeNodeResourceHandler(c *slurmclient.Client) server.ResourceHandlerFunc {
	return func(ctx context.Context, req mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
		nodes, err := c.Slurm.Nodes().List(ctx)
		if err != nil {
			return nil, err
		}
		data, _ := json.MarshalIndent(nodes, "", "  ")
		return []mcp.ResourceContents{
			mcp.TextResourceContents{
				URI:      req.Params.URI,
				MIMEType: "application/json",
				Text:     string(data),
			},
		}, nil
	}
}

func makeJobResourceHandler(c *slurmclient.Client) server.ResourceHandlerFunc {
	return func(ctx context.Context, req mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
		jobs, err := c.Slurm.Jobs().List(ctx)
		if err != nil {
			return nil, err
		}
		data, _ := json.MarshalIndent(jobs, "", "  ")
		return []mcp.ResourceContents{
			mcp.TextResourceContents{
				URI:      req.Params.URI,
				MIMEType: "application/json",
				Text:     string(data),
			},
		}, nil
	}
}

func makePartitionResourceHandler(c *slurmclient.Client) server.ResourceHandlerFunc {
	return func(ctx context.Context, req mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
		parts, err := c.Slurm.Partitions().List(ctx)
		if err != nil {
			return nil, err
		}
		data, _ := json.MarshalIndent(parts, "", "  ")
		return []mcp.ResourceContents{
			mcp.TextResourceContents{
				URI:      req.Params.URI,
				MIMEType: "application/json",
				Text:     string(data),
			},
		}, nil
	}
}
