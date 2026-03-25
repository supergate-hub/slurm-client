package main

import (
	"context"
	"fmt"
	"strings"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"

	slurmclient "github.com/supergate-hub/slurm-client"
)

func registerAnalysisTools(s *server.MCPServer, cr *clusterResolver, mgr *slurmclient.Manager) {
	clusterDesc := ""
	if cr.isMultiCluster() {
		clusterDesc = " Use the 'cluster' parameter to target a specific cluster."
	}

	// --- GPU/GRES Allocation ---
	s.AddTool(mcp.NewTool("slurm_gpu_allocation",
		mcp.WithDescription("Show GPU and GRES (Generic Resource) allocation status per node. "+
			"Reports total vs used GRES resources. "+
			"Note: this shows GRES allocation from Slurm's perspective, not real-time GPU utilization (nvidia-smi level)."+clusterDesc),
		mcp.WithString("cluster", mcp.Description("Target cluster name (required in multi-cluster mode)")),
	), func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		if mgr != nil && cr.isMultiCluster() {
			// Multi-cluster: aggregate across all clusters
			results := mgr.AllNodes(ctx)
			type clusterGres struct {
				Cluster string     `json:"cluster"`
				Nodes   []nodeGres `json:"nodes,omitempty"`
				Error   string     `json:"error,omitempty"`
			}
			var out []clusterGres
			for _, r := range results {
				cg := clusterGres{Cluster: r.Cluster}
				if r.Err != nil {
					cg.Error = r.Err.Error()
				} else {
					cg.Nodes = extractGresFromNodes(r.Data)
				}
				out = append(out, cg)
			}
			return jsonResult(out)
		}

		client, err := cr.resolve(req)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		nodes, err := client.Slurm.Nodes().List(ctx)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		return jsonResult(extractGresFromNodes(nodes))
	})

	// --- Queue Depth ---
	s.AddTool(mcp.NewTool("slurm_queue_depth",
		mcp.WithDescription("Analyze job queue depth: count of pending vs running jobs, grouped by partition and user. "+
			"Helps identify bottlenecks and resource contention."+clusterDesc),
		mcp.WithString("cluster", mcp.Description("Target cluster name (required in multi-cluster mode)")),
	), func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		if mgr != nil && cr.isMultiCluster() {
			results := mgr.AllJobs(ctx)
			type clusterQueue struct {
				Cluster string     `json:"cluster"`
				Queue   *queueInfo `json:"queue,omitempty"`
				Error   string     `json:"error,omitempty"`
			}
			var out []clusterQueue
			for _, r := range results {
				cq := clusterQueue{Cluster: r.Cluster}
				if r.Err != nil {
					cq.Error = r.Err.Error()
				} else {
					qi := analyzeQueue(r.Data)
					cq.Queue = &qi
				}
				out = append(out, cq)
			}
			return jsonResult(out)
		}

		client, err := cr.resolve(req)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		jobs, err := client.Slurm.Jobs().List(ctx)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		return jsonResult(analyzeQueue(jobs))
	})
}

type nodeGres struct {
	Name        string `json:"name"`
	Gres        string `json:"gres,omitempty"`
	GresUsed    string `json:"gres_used,omitempty"`
	GresDrained string `json:"gres_drained,omitempty"`
	Tres        string `json:"tres,omitempty"`
	TresUsed    string `json:"tres_used,omitempty"`
	HasGPU      bool   `json:"has_gpu"`
}

func extractGresFromNodes(nodes interface{}) []nodeGres {
	// Use type assertion since we accept v0040 nodes
	type nodeWithGres interface {
		GetName() string
	}

	// Work with the raw slice via reflection-free approach
	// Since we know these are v0040.V0040Node, access fields directly
	var result []nodeGres

	// Type-switch for the actual v0040 node slice
	switch ns := nodes.(type) {
	case []struct{ Name *string }:
		// fallback
		_ = ns
	default:
		// Use JSON round-trip for generic access
		data, _ := jsonMarshal(nodes)
		var rawNodes []map[string]interface{}
		jsonUnmarshal(data, &rawNodes)

		for _, n := range rawNodes {
			ng := nodeGres{}
			if name, ok := n["name"].(string); ok {
				ng.Name = name
			}
			if gres, ok := n["gres"].(string); ok {
				ng.Gres = gres
				ng.HasGPU = strings.Contains(strings.ToLower(gres), "gpu")
			}
			if gresUsed, ok := n["gres_used"].(string); ok {
				ng.GresUsed = gresUsed
			}
			if gresDrained, ok := n["gres_drained"].(string); ok {
				ng.GresDrained = gresDrained
			}
			if tres, ok := n["tres"].(string); ok {
				ng.Tres = tres
			}
			if tresUsed, ok := n["tres_used"].(string); ok {
				ng.TresUsed = tresUsed
			}
			result = append(result, ng)
		}
	}

	return result
}

type queueInfo struct {
	Total        int                       `json:"total_jobs"`
	Running      int                       `json:"running"`
	Pending      int                       `json:"pending"`
	Other        int                       `json:"other"`
	ByPartition  map[string]partitionQueue `json:"by_partition"`
	ByUser       map[string]userQueue      `json:"by_user"`
	TopPending   []pendingJob              `json:"top_pending,omitempty"`
}

type partitionQueue struct {
	Running int `json:"running"`
	Pending int `json:"pending"`
}

type userQueue struct {
	Running int `json:"running"`
	Pending int `json:"pending"`
}

type pendingJob struct {
	JobID  int32  `json:"job_id"`
	Name   string `json:"name"`
	User   string `json:"user"`
	Reason string `json:"reason"`
}

func analyzeQueue(jobs interface{}) queueInfo {
	qi := queueInfo{
		ByPartition: make(map[string]partitionQueue),
		ByUser:      make(map[string]userQueue),
	}

	data, _ := jsonMarshal(jobs)
	var rawJobs []map[string]interface{}
	jsonUnmarshal(data, &rawJobs)

	for _, j := range rawJobs {
		qi.Total++

		states, _ := j["job_state"].([]interface{})
		state := ""
		if len(states) > 0 {
			state = fmt.Sprintf("%v", states[0])
		}

		partition, _ := j["partition"].(string)
		userName, _ := j["user_name"].(string)

		pq := qi.ByPartition[partition]
		uq := qi.ByUser[userName]

		switch state {
		case "RUNNING":
			qi.Running++
			pq.Running++
			uq.Running++
		case "PENDING":
			qi.Pending++
			pq.Pending++
			uq.Pending++

			if len(qi.TopPending) < 10 {
				jobID := int32(0)
				if id, ok := j["job_id"].(float64); ok {
					jobID = int32(id)
				}
				name, _ := j["name"].(string)
				reason, _ := j["state_description"].(string)
				qi.TopPending = append(qi.TopPending, pendingJob{
					JobID: jobID, Name: name, User: userName, Reason: reason,
				})
			}
		default:
			qi.Other++
		}

		qi.ByPartition[partition] = pq
		qi.ByUser[userName] = uq
	}

	return qi
}
