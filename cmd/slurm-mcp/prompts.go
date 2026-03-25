package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"

	slurmclient "github.com/supergate-hub/slurm-client"
)

func registerPrompts(s *server.MCPServer, cr *clusterResolver, mgr *slurmclient.Manager) {
	// --- Cluster Status Summary ---
	s.AddPrompt(mcp.Prompt{
		Name:        "cluster-status-summary",
		Description: "Generate a human-readable status summary of the Slurm cluster(s): node health, active jobs, and partition utilization.",
		Arguments: []mcp.PromptArgument{
			{Name: "cluster", Description: "Target cluster name (optional in single-cluster mode)", Required: false},
		},
	}, func(ctx context.Context, req mcp.GetPromptRequest) (*mcp.GetPromptResult, error) {
		data, err := gatherClusterStatus(ctx, cr, mgr, req.Params.Arguments["cluster"])
		if err != nil {
			return nil, err
		}
		return &mcp.GetPromptResult{
			Description: "Cluster status summary",
			Messages: []mcp.PromptMessage{
				mcp.NewPromptMessage(mcp.RoleUser, mcp.NewTextContent(
					fmt.Sprintf("Here is the current Slurm cluster status data. Please summarize it in a clear, human-readable format highlighting any issues:\n\n%s", data),
				)),
			},
		}, nil
	})

	// --- Job Queue Analysis ---
	s.AddPrompt(mcp.Prompt{
		Name:        "job-queue-analysis",
		Description: "Analyze the current job queue: pending jobs, running jobs, resource utilization, and potential bottlenecks.",
		Arguments: []mcp.PromptArgument{
			{Name: "cluster", Description: "Target cluster name (optional in single-cluster mode)", Required: false},
		},
	}, func(ctx context.Context, req mcp.GetPromptRequest) (*mcp.GetPromptResult, error) {
		client, err := resolveFromArgs(cr, mgr, req.Params.Arguments["cluster"])
		if err != nil {
			return nil, err
		}
		jobs, err := client.Slurm.Jobs().List(ctx)
		if err != nil {
			return nil, err
		}
		data, _ := json.MarshalIndent(jobs, "", "  ")
		return &mcp.GetPromptResult{
			Description: "Job queue analysis",
			Messages: []mcp.PromptMessage{
				mcp.NewPromptMessage(mcp.RoleUser, mcp.NewTextContent(
					fmt.Sprintf("Analyze this Slurm job queue data. Identify:\n1. How many jobs are pending vs running\n2. Which users have the most jobs\n3. Resource utilization patterns\n4. Any potential bottlenecks\n\n%s", string(data)),
				)),
			},
		}, nil
	})

	// --- Node Troubleshooting ---
	s.AddPrompt(mcp.Prompt{
		Name:        "node-troubleshooting",
		Description: "Diagnose node issues: identify down/drained/error nodes and suggest remediation steps.",
		Arguments: []mcp.PromptArgument{
			{Name: "cluster", Description: "Target cluster name (optional in single-cluster mode)", Required: false},
			{Name: "node_name", Description: "Specific node to investigate (optional, analyzes all if omitted)", Required: false},
		},
	}, func(ctx context.Context, req mcp.GetPromptRequest) (*mcp.GetPromptResult, error) {
		client, err := resolveFromArgs(cr, mgr, req.Params.Arguments["cluster"])
		if err != nil {
			return nil, err
		}

		nodeName := req.Params.Arguments["node_name"]
		var dataStr string
		if nodeName != "" {
			node, err := client.Slurm.Nodes().Get(ctx, nodeName)
			if err != nil {
				return nil, err
			}
			d, _ := json.MarshalIndent(node, "", "  ")
			dataStr = string(d)
		} else {
			nodes, err := client.Slurm.Nodes().List(ctx)
			if err != nil {
				return nil, err
			}
			d, _ := json.MarshalIndent(nodes, "", "  ")
			dataStr = string(d)
		}

		return &mcp.GetPromptResult{
			Description: "Node troubleshooting",
			Messages: []mcp.PromptMessage{
				mcp.NewPromptMessage(mcp.RoleUser, mcp.NewTextContent(
					fmt.Sprintf("Diagnose these Slurm node(s). For any node that is down, drained, or in an error state:\n1. What is the likely cause\n2. What remediation steps should be taken\n3. What impact does this have on running jobs\n\n%s", dataStr),
				)),
			},
		}, nil
	})
}

// gatherClusterStatus collects status data for the prompt
func gatherClusterStatus(ctx context.Context, cr *clusterResolver, mgr *slurmclient.Manager, clusterArg string) (string, error) {
	if mgr != nil && clusterArg == "" {
		// Multi-cluster: gather all
		type summary struct {
			Cluster string `json:"cluster"`
			Nodes   int    `json:"nodes"`
			Jobs    int    `json:"jobs"`
			Error   string `json:"error,omitempty"`
		}
		var summaries []summary
		for _, name := range mgr.ClusterNames() {
			c, _ := mgr.Cluster(name)
			s := summary{Cluster: name}
			nodes, err := c.Slurm.Nodes().List(ctx)
			if err != nil {
				s.Error = err.Error()
			} else {
				s.Nodes = len(nodes)
			}
			jobs, err := c.Slurm.Jobs().List(ctx)
			if err == nil {
				s.Jobs = len(jobs)
			}
			summaries = append(summaries, s)
		}
		data, _ := json.MarshalIndent(summaries, "", "  ")
		return string(data), nil
	}

	client, err := resolveFromArgs(cr, mgr, clusterArg)
	if err != nil {
		return "", err
	}
	type singleSummary struct {
		Nodes      any `json:"nodes"`
		Jobs       any `json:"jobs"`
		Partitions any `json:"partitions"`
	}
	var ss singleSummary
	ss.Nodes, _ = client.Slurm.Nodes().List(ctx)
	ss.Jobs, _ = client.Slurm.Jobs().List(ctx)
	ss.Partitions, _ = client.Slurm.Partitions().List(ctx)
	data, _ := json.MarshalIndent(ss, "", "  ")
	return string(data), nil
}

func resolveFromArgs(cr *clusterResolver, mgr *slurmclient.Manager, clusterArg string) (*slurmclient.Client, error) {
	if cr.single != nil {
		return cr.single, nil
	}
	if clusterArg == "" {
		names := mgr.ClusterNames()
		if len(names) == 1 {
			return mgr.Cluster(names[0])
		}
		return nil, fmt.Errorf("cluster argument is required in multi-cluster mode (available: %v)", names)
	}
	return mgr.Cluster(clusterArg)
}
