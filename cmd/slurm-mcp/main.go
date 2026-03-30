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
//
// SSE transport (remote/Docker):
//
//	slurm-mcp --transport=sse --port=8080 --bearer-token=secret
package main

import (
	"context"
	"crypto/subtle"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"

	slurmclient "github.com/supergate-hub/slurm-client"
)

// coreTools defines the 8 default tools exposed without --all-tools.
var coreTools = map[string]bool{
	"slurm_ping":            true,
	"slurm_list_jobs":       true,
	"slurm_get_job":         true,
	"slurm_submit_job":      true,
	"slurm_cancel_job":      true,
	"slurm_list_nodes":      true,
	"slurm_list_partitions": true,
	"slurm_queue_depth":     true,
}

// clusterResolver resolves a cluster name from an MCP request to a *Client.
// In single-cluster mode, the cluster parameter is ignored.
// In multi-cluster mode, it routes to the appropriate cluster.
type clusterResolver struct {
	single    *slurmclient.Client  // non-nil in single-cluster mode
	manager   *slurmclient.Manager // non-nil in multi-cluster mode
	backendFn func(*slurmclient.Client) MCPBackend
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

func (cr *clusterResolver) resolveBackend(req mcp.CallToolRequest) (MCPBackend, error) {
	client, err := cr.resolve(req)
	if err != nil {
		return nil, err
	}
	return cr.backendFn(client), nil
}

func (cr *clusterResolver) isMultiCluster() bool {
	return cr.manager != nil
}

// toolRegistrar wraps server.MCPServer to gate tool registration based on --all-tools.
type toolRegistrar struct {
	server   *server.MCPServer
	allTools bool
}

// AddTool registers a tool only if it's a core tool or --all-tools is enabled.
func (tr *toolRegistrar) AddTool(tool mcp.Tool, handler func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error)) {
	if tr.allTools || coreTools[tool.Name] {
		tr.server.AddTool(tool, handler)
	}
}

// checkRBAC verifies access before a tool handler runs.
// Returns nil if RBAC is not configured or if access is granted.
func checkRBAC(rbac *RBAC, toolName string) error {
	if rbac == nil || !rbac.NeedsEnforcement() {
		return nil
	}
	return rbac.Check(toolName)
}

func main() {
	configPath := flag.String("config", "", "Path to clusters.yaml for multi-cluster mode")
	transport := flag.String("transport", "stdio", "Transport mode: stdio or sse")
	port := flag.String("port", "8080", "Port for SSE transport")
	bearerToken := flag.String("bearer-token", "", "Bearer token for SSE transport authentication (required for sse)")
	allTools := flag.Bool("all-tools", false, "Expose all 31 tools (default: 8 core tools)")
	backend := flag.String("backend", "rest", "Backend type: rest or ssh")
	flag.Parse()

	// Validate backend flag
	switch *backend {
	case "rest":
		// OK
	case "ssh":
		fmt.Fprintln(os.Stderr, "SSH backend not yet available, coming in next release")
		os.Exit(1)
	default:
		log.Fatalf("Unknown backend: %s (supported: rest, ssh)", *backend)
	}

	ctx := context.Background()
	var resolver clusterResolver
	var mgr *slurmclient.Manager

	// Set the backend factory function
	resolver.backendFn = func(c *slurmclient.Client) MCPBackend {
		return NewRestBackend(c)
	}

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

	// RBAC: parse from clusters.yaml if present, or use default (admin)
	var rbacCfg RBACConfig
	if *configPath != "" {
		if cfg := ParseRBACFromClustersConfig(*configPath); cfg != nil {
			rbacCfg = *cfg
		}
	}
	if accessLevel := os.Getenv("MCP_ACCESS_LEVEL"); accessLevel != "" {
		rbacCfg.DefaultAccess = AccessLevel(accessLevel)
	}
	if auditLog := os.Getenv("MCP_AUDIT_LOG"); auditLog != "" {
		rbacCfg.AuditLog = auditLog
	}
	rbac, err := NewRBAC(rbacCfg)
	if err != nil {
		log.Fatalf("Failed to initialize RBAC: %v", err)
	}
	defer rbac.Close()
	if rbac.NeedsEnforcement() || rbacCfg.AuditLog != "" {
		rbac.LogStartup()
	}

	tr := &toolRegistrar{server: s, allTools: *allTools}

	registerSlurmTools(tr, &resolver, rbac)
	registerSlurmdbTools(tr, &resolver, rbac)
	registerResources(s, &resolver, mgr)
	registerPrompts(s, &resolver, mgr)
	registerAnalysisTools(tr, &resolver, mgr, rbac)

	if resolver.isMultiCluster() {
		registerMultiClusterTools(tr, mgr, rbac)
	}

	switch *transport {
	case "stdio":
		if err := server.ServeStdio(s); err != nil {
			fmt.Fprintf(os.Stderr, "Server error: %v\n", err)
			os.Exit(1)
		}

	case "sse":
		// Bearer token is required for network-facing transport
		token := *bearerToken
		if token == "" {
			token = os.Getenv("MCP_BEARER_TOKEN")
		}
		if token == "" {
			log.Fatal("--bearer-token or MCP_BEARER_TOKEN is required for SSE transport")
		}

		addr := ":" + *port
		sseServer := server.NewSSEServer(s,
			server.WithBaseURL(fmt.Sprintf("http://localhost%s", addr)),
		)

		// Wrap with bearer token auth middleware
		handler := bearerAuthMiddleware(token, sseServer)

		log.Printf("Starting SSE transport on %s (auth: bearer token)", addr)
		httpServer := &http.Server{Addr: addr, Handler: handler}
		if err := httpServer.ListenAndServe(); err != nil {
			log.Fatalf("SSE server error: %v", err)
		}

	default:
		log.Fatalf("Unknown transport: %s (supported: stdio, sse)", *transport)
	}
}

// bearerAuthMiddleware rejects requests without a valid Bearer token.
func bearerAuthMiddleware(expectedToken string, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		auth := r.Header.Get("Authorization")
		if len(auth) < 7 || auth[:7] != "Bearer " {
			http.Error(w, "Unauthorized: Bearer token required", http.StatusUnauthorized)
			return
		}
		provided := auth[7:]
		if subtle.ConstantTimeCompare([]byte(provided), []byte(expectedToken)) != 1 {
			http.Error(w, "Unauthorized: invalid token", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}
