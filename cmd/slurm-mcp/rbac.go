package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"gopkg.in/yaml.v2"
)

// AccessLevel defines the permission tier for MCP operations.
type AccessLevel string

const (
	// AccessReadOnly allows only read operations (list, get, ping, diag).
	AccessReadOnly AccessLevel = "read-only"
	// AccessOperator allows read + job submit/cancel, but not node/cluster management.
	AccessOperator AccessLevel = "operator"
	// AccessAdmin allows all operations.
	AccessAdmin AccessLevel = "admin"
)

// RBACConfig defines access control rules for the MCP server.
type RBACConfig struct {
	// DefaultAccess is the access level when no specific rule matches.
	DefaultAccess AccessLevel `yaml:"default_access"`
	// AuditLog path for operation audit logging. Empty disables audit logging.
	AuditLog string `yaml:"audit_log"`
}

// RBAC enforces access control and audit logging for MCP tool calls.
type RBAC struct {
	config   RBACConfig
	auditLog *os.File
}

// NewRBAC creates an RBAC enforcer from config.
func NewRBAC(cfg RBACConfig) (*RBAC, error) {
	r := &RBAC{config: cfg}

	if cfg.AuditLog != "" {
		f, err := os.OpenFile(cfg.AuditLog, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return nil, fmt.Errorf("open audit log %s: %w", cfg.AuditLog, err)
		}
		r.auditLog = f
	}

	if cfg.DefaultAccess == "" {
		r.config.DefaultAccess = AccessAdmin
	}

	return r, nil
}

// Close closes the audit log file.
func (r *RBAC) Close() {
	if r.auditLog != nil {
		r.auditLog.Close()
	}
}

// Check verifies whether the given tool is allowed under the current access level.
// Returns nil if allowed, error if denied.
func (r *RBAC) Check(toolName string) error {
	required := toolAccessLevel(toolName)
	if !r.hasAccess(required) {
		return fmt.Errorf("access denied: tool %q requires %s access (current: %s)", toolName, required, r.config.DefaultAccess)
	}
	return nil
}

// Audit logs a tool invocation to the audit log.
func (r *RBAC) Audit(toolName string, args map[string]any, allowed bool) {
	if r.auditLog == nil {
		return
	}

	entry := map[string]any{
		"timestamp": time.Now().UTC().Format(time.RFC3339),
		"tool":      toolName,
		"allowed":   allowed,
		"access":    string(r.config.DefaultAccess),
	}

	// Log cluster parameter if present (but not sensitive args like scripts)
	if cluster, ok := args["cluster"]; ok {
		entry["cluster"] = cluster
	}
	if jobID, ok := args["job_id"]; ok {
		entry["job_id"] = jobID
	}

	data, _ := json.Marshal(entry)
	fmt.Fprintln(r.auditLog, string(data))
}

func (r *RBAC) hasAccess(required AccessLevel) bool {
	switch r.config.DefaultAccess {
	case AccessAdmin:
		return true
	case AccessOperator:
		return required == AccessReadOnly || required == AccessOperator
	case AccessReadOnly:
		return required == AccessReadOnly
	default:
		return false
	}
}

// toolAccessLevel maps tool names to their required access level.
func toolAccessLevel(name string) AccessLevel {
	// Write operations require operator or admin
	writeTools := map[string]AccessLevel{
		"slurm_submit_job": AccessOperator,
		"slurm_cancel_job": AccessOperator,
	}
	if level, ok := writeTools[name]; ok {
		return level
	}

	// Admin-only operations
	adminTools := []string{
		// Node management would go here if we add drain/resume
	}
	for _, t := range adminTools {
		if name == t {
			return AccessAdmin
		}
	}

	// Everything else is read-only
	return AccessReadOnly
}

// ParseRBACConfig reads RBAC config from a YAML file.
func ParseRBACConfig(path string) (*RBACConfig, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	expanded := os.ExpandEnv(string(data))

	var cfg struct {
		RBAC RBACConfig `yaml:"rbac"`
	}
	if err := yaml.Unmarshal([]byte(expanded), &cfg); err != nil {
		return nil, fmt.Errorf("parse rbac config: %w", err)
	}
	return &cfg.RBAC, nil
}

// ParseRBACFromClustersConfig extracts RBAC config from clusters.yaml if present.
func ParseRBACFromClustersConfig(path string) *RBACConfig {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil
	}
	expanded := os.ExpandEnv(string(data))

	var cfg struct {
		RBAC *RBACConfig `yaml:"rbac"`
	}
	if err := yaml.Unmarshal([]byte(expanded), &cfg); err != nil {
		return nil
	}
	return cfg.RBAC
}

// LogStartup logs the RBAC configuration at startup.
func (r *RBAC) LogStartup() {
	log.Printf("RBAC: access_level=%s audit_log=%s", r.config.DefaultAccess, r.auditLogPath())

	readTools, opTools, adminTools := 0, 0, 0
	for _, name := range allToolNames() {
		switch toolAccessLevel(name) {
		case AccessReadOnly:
			readTools++
		case AccessOperator:
			opTools++
		case AccessAdmin:
			adminTools++
		}
	}
	log.Printf("RBAC: %d read-only tools, %d operator tools, %d admin tools", readTools, opTools, adminTools)
}

func (r *RBAC) auditLogPath() string {
	if r.config.AuditLog == "" {
		return "(disabled)"
	}
	return r.config.AuditLog
}

func allToolNames() []string {
	// All registered tool names for access level summary
	return []string{
		"slurm_ping", "slurm_list_jobs", "slurm_get_job", "slurm_submit_job",
		"slurm_cancel_job", "slurm_list_nodes", "slurm_get_node",
		"slurm_list_partitions", "slurm_get_partition", "slurm_list_reservations",
		"slurm_list_licenses", "slurm_diag",
		"slurmdb_list_accounts", "slurmdb_get_account", "slurmdb_list_users",
		"slurmdb_get_user", "slurmdb_list_jobs", "slurmdb_get_job",
		"slurmdb_list_clusters", "slurmdb_list_qos", "slurmdb_list_associations",
		"slurmdb_list_tres", "slurmdb_list_wckeys", "slurmdb_diag",
		"slurm_list_clusters", "slurm_cluster_status", "slurm_cross_cluster_jobs",
		"slurm_node_health_summary", "slurm_resource_availability",
		"slurm_gpu_allocation", "slurm_queue_depth",
	}
}

// WrapToolDescription appends access level info to tool descriptions.
func AccessBadge(toolName string) string {
	switch toolAccessLevel(toolName) {
	case AccessOperator:
		return " [requires: operator]"
	case AccessAdmin:
		return " [requires: admin]"
	default:
		return ""
	}
}

// NeedsRBAC returns true if any RBAC restriction is configured (not admin mode).
func (r *RBAC) NeedsEnforcement() bool {
	return r.config.DefaultAccess != AccessAdmin && r.config.DefaultAccess != ""
}

// Suppress unused import warning
var _ = strings.Contains
