package main

import (
	"context"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"

	slurmclient "github.com/supergate-hub/slurm-client"
)

func registerMultiClusterTools(s toolAdder, mgr *slurmclient.Manager, rbac *RBAC) {
	// --- Multi-Cluster: List Clusters ---
	s.AddTool(mcp.NewTool("slurm_list_clusters",
		mcp.WithDescription("List all connected Slurm clusters with their API versions."),
	), func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		if err := checkRBAC(rbac, "slurm_list_clusters"); err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		type clusterInfo struct {
			Name    string `json:"name"`
			Version string `json:"version"`
		}
		var clusters []clusterInfo
		for _, name := range mgr.ClusterNames() {
			c, _ := mgr.Cluster(name)
			clusters = append(clusters, clusterInfo{
				Name:    name,
				Version: string(c.Version()),
			})
		}
		return jsonResult(clusters)
	})

	// --- Multi-Cluster: Cluster Status ---
	s.AddTool(mcp.NewTool("slurm_cluster_status",
		mcp.WithDescription("Get a health summary for each connected cluster: ping status, node counts, and job counts."),
	), func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		if err := checkRBAC(rbac, "slurm_cluster_status"); err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		type clusterStatus struct {
			Name       string `json:"name"`
			Reachable  bool   `json:"reachable"`
			TotalNodes int    `json:"total_nodes,omitempty"`
			TotalJobs  int    `json:"total_jobs,omitempty"`
			Error      string `json:"error,omitempty"`
		}

		var statuses []clusterStatus
		for _, name := range mgr.ClusterNames() {
			c, _ := mgr.Cluster(name)
			st := clusterStatus{Name: name}

			// Check reachability via ping
			_, err := c.Slurm.Ping().Get(ctx)
			if err != nil {
				st.Error = err.Error()
				statuses = append(statuses, st)
				continue
			}
			st.Reachable = true

			// Node count
			nodes, err := c.Slurm.Nodes().List(ctx)
			if err == nil {
				st.TotalNodes = len(nodes)
			}

			// Job count
			jobs, err := c.Slurm.Jobs().List(ctx)
			if err == nil {
				st.TotalJobs = len(jobs)
			}

			statuses = append(statuses, st)
		}
		return jsonResult(statuses)
	})

	// --- Multi-Cluster: Cross-Cluster Jobs ---
	s.AddTool(mcp.NewTool("slurm_cross_cluster_jobs",
		mcp.WithDescription("List jobs from ALL connected clusters. Returns results grouped by cluster, including partial results if some clusters are unreachable."),
	), func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		if err := checkRBAC(rbac, "slurm_cross_cluster_jobs"); err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		results := mgr.AllJobs(ctx)

		type clusterJobs struct {
			Cluster  string `json:"cluster"`
			JobCount int    `json:"job_count"`
			Jobs     any    `json:"jobs,omitempty"`
			Error    string `json:"error,omitempty"`
		}

		var out []clusterJobs
		for _, r := range results {
			cj := clusterJobs{Cluster: r.Cluster}
			if r.Err != nil {
				cj.Error = r.Err.Error()
			} else {
				cj.JobCount = len(r.Data)
				cj.Jobs = r.Data
			}
			out = append(out, cj)
		}
		return jsonResult(out)
	})

	// --- Multi-Cluster: Node Health Summary ---
	s.AddTool(mcp.NewTool("slurm_node_health_summary",
		mcp.WithDescription("Get a summary of node health across all clusters: counts of idle, allocated, mixed, down, and drained nodes per cluster."),
	), func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		if err := checkRBAC(rbac, "slurm_node_health_summary"); err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		results := mgr.AllNodes(ctx)

		type nodeHealth struct {
			Cluster string         `json:"cluster"`
			Total   int            `json:"total"`
			States  map[string]int `json:"states"`
			Error   string         `json:"error,omitempty"`
		}

		var out []nodeHealth
		for _, r := range results {
			nh := nodeHealth{Cluster: r.Cluster}
			if r.Err != nil {
				nh.Error = r.Err.Error()
				out = append(out, nh)
				continue
			}
			nh.Total = len(r.Data)
			nh.States = make(map[string]int)
			for _, n := range r.Data {
				state := "unknown"
				if n.State != nil && len(*n.State) > 0 {
					state = fmt.Sprintf("%v", (*n.State)[0])
				}
				nh.States[state]++
			}
			out = append(out, nh)
		}
		return jsonResult(out)
	})

	// --- Multi-Cluster: Resource Availability ---
	s.AddTool(mcp.NewTool("slurm_resource_availability",
		mcp.WithDescription("Show available resources per cluster and partition: total/idle nodes, total/idle CPUs. AI agents use this to recommend which cluster and partition to submit jobs to."),
	), func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		if err := checkRBAC(rbac, "slurm_resource_availability"); err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		results := mgr.AllPartitions(ctx)

		type partitionInfo struct {
			Name      string `json:"name"`
			TotalCPUs int    `json:"total_cpus,omitempty"`
			TotalNodes int   `json:"total_nodes,omitempty"`
		}
		type clusterResources struct {
			Cluster    string          `json:"cluster"`
			Partitions []partitionInfo `json:"partitions,omitempty"`
			Error      string          `json:"error,omitempty"`
		}

		var out []clusterResources
		for _, r := range results {
			cr := clusterResources{Cluster: r.Cluster}
			if r.Err != nil {
				cr.Error = r.Err.Error()
				out = append(out, cr)
				continue
			}
			for _, p := range r.Data {
				pi := partitionInfo{}
				if p.Name != nil {
					pi.Name = *p.Name
				}
				if p.Cpus != nil {
					if p.Cpus.Total != nil {
						pi.TotalCPUs = int(*p.Cpus.Total)
					}
				}
				if p.Nodes != nil && p.Nodes.Total != nil {
					pi.TotalNodes = int(*p.Nodes.Total)
				}
				cr.Partitions = append(cr.Partitions, pi)
			}
			out = append(out, cr)
		}
		return jsonResult(out)
	})
}
