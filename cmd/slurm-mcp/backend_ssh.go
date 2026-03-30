package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"golang.org/x/crypto/ssh"
)

// SSHConfig holds configuration for the SSH backend.
type SSHConfig struct {
	Host    string        // hostname:port
	User    string        // SSH user
	KeyPath string        // path to SSH private key (default: ~/.ssh/id_rsa)
	Timeout time.Duration // connection timeout (default: 10s)
}

// SSHError represents a command execution failure over SSH.
type SSHError struct {
	Command  string
	ExitCode int
	Stderr   string
}

func (e *SSHError) Error() string {
	return fmt.Sprintf("SSH command %q failed (exit %d): %s", e.Command, e.ExitCode, e.Stderr)
}

// SSHBackend implements MCPBackend by running Slurm CLI commands over SSH.
type SSHBackend struct {
	config   SSHConfig
	mu       sync.Mutex
	client   *ssh.Client
	slurmVer string // detected version, e.g. "24.05.2"
	hasJSON  bool   // true if Slurm >= 21.08 (supports --json flag)
}

// NewSSHBackend creates and connects an SSH backend.
func NewSSHBackend(ctx context.Context, cfg SSHConfig) (*SSHBackend, error) {
	if cfg.KeyPath == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			return nil, fmt.Errorf("cannot determine home directory: %w", err)
		}
		cfg.KeyPath = home + "/.ssh/id_rsa"
	}
	if cfg.Timeout == 0 {
		cfg.Timeout = 10 * time.Second
	}
	if !strings.Contains(cfg.Host, ":") {
		cfg.Host = cfg.Host + ":22"
	}

	b := &SSHBackend{config: cfg}
	if err := b.connect(); err != nil {
		return nil, err
	}
	if err := b.detectVersion(ctx); err != nil {
		return nil, fmt.Errorf("failed to detect Slurm version: %w", err)
	}
	return b, nil
}

func (b *SSHBackend) connect() error {
	keyBytes, err := os.ReadFile(b.config.KeyPath)
	if err != nil {
		return fmt.Errorf("read SSH key %s: %w", b.config.KeyPath, err)
	}
	signer, err := ssh.ParsePrivateKey(keyBytes)
	if err != nil {
		return fmt.Errorf("parse SSH key: %w", err)
	}

	clientCfg := &ssh.ClientConfig{
		User:            b.config.User,
		Auth:            []ssh.AuthMethod{ssh.PublicKeys(signer)},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(), // #nosec - cluster environments typically use known hosts
		Timeout:         b.config.Timeout,
	}

	client, err := ssh.Dial("tcp", b.config.Host, clientCfg)
	if err != nil {
		return fmt.Errorf("SSH dial %s: %w", b.config.Host, err)
	}
	b.client = client
	return nil
}

// ensureConnected reconnects if the connection is broken.
func (b *SSHBackend) ensureConnected() error {
	b.mu.Lock()
	defer b.mu.Unlock()

	if b.client != nil {
		// Check if connection is still alive
		_, _, err := b.client.SendRequest("keepalive@openssh.com", true, nil)
		if err == nil {
			return nil
		}
		_ = b.client.Close()
		b.client = nil
	}
	return b.connect()
}

// runCommand executes a command over SSH with a 30-second timeout.
func (b *SSHBackend) runCommand(ctx context.Context, command string) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	if err := b.ensureConnected(); err != nil {
		return "", err
	}

	session, err := b.client.NewSession()
	if err != nil {
		return "", fmt.Errorf("create SSH session: %w", err)
	}
	defer session.Close()

	var stdout, stderr bytes.Buffer
	session.Stdout = &stdout
	session.Stderr = &stderr

	done := make(chan error, 1)
	go func() {
		done <- session.Run(command)
	}()

	select {
	case <-ctx.Done():
		_ = session.Signal(ssh.SIGTERM)
		return "", ctx.Err()
	case err := <-done:
		if err != nil {
			exitCode := 1
			if exitErr, ok := err.(*ssh.ExitError); ok {
				exitCode = exitErr.ExitStatus()
			}
			return "", &SSHError{
				Command:  command,
				ExitCode: exitCode,
				Stderr:   strings.TrimSpace(stderr.String()),
			}
		}
	}

	return stdout.String(), nil
}

// runCommandWithStdin executes a command over SSH with stdin data.
func (b *SSHBackend) runCommandWithStdin(ctx context.Context, command, stdin string) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	if err := b.ensureConnected(); err != nil {
		return "", err
	}

	session, err := b.client.NewSession()
	if err != nil {
		return "", fmt.Errorf("create SSH session: %w", err)
	}
	defer session.Close()

	session.Stdin = strings.NewReader(stdin)

	var stdout, stderr bytes.Buffer
	session.Stdout = &stdout
	session.Stderr = &stderr

	done := make(chan error, 1)
	go func() {
		done <- session.Run(command)
	}()

	select {
	case <-ctx.Done():
		_ = session.Signal(ssh.SIGTERM)
		return "", ctx.Err()
	case err := <-done:
		if err != nil {
			exitCode := 1
			if exitErr, ok := err.(*ssh.ExitError); ok {
				exitCode = exitErr.ExitStatus()
			}
			return "", &SSHError{
				Command:  command,
				ExitCode: exitCode,
				Stderr:   strings.TrimSpace(stderr.String()),
			}
		}
	}

	return stdout.String(), nil
}

// detectVersion runs squeue --version and parses the Slurm version.
func (b *SSHBackend) detectVersion(ctx context.Context) error {
	output, err := b.runCommand(ctx, "squeue --version")
	if err != nil {
		return err
	}
	b.slurmVer = parseSlurmVersion(output)
	b.hasJSON = versionSupportsJSON(b.slurmVer)
	return nil
}

// parseSlurmVersion extracts the version string from squeue --version output.
// Expected format: "slurm 24.05.2" or "slurm 21.08.1-2"
func parseSlurmVersion(output string) string {
	re := regexp.MustCompile(`(\d+\.\d+\.\d+)`)
	match := re.FindString(strings.TrimSpace(output))
	if match == "" {
		return "0.0.0"
	}
	return match
}

// versionSupportsJSON returns true if version >= 21.08.
func versionSupportsJSON(version string) bool {
	parts := strings.SplitN(version, ".", 3)
	if len(parts) < 2 {
		return false
	}
	major, err := strconv.Atoi(parts[0])
	if err != nil {
		return false
	}
	minor, err := strconv.Atoi(parts[1])
	if err != nil {
		return false
	}
	return major > 21 || (major == 21 && minor >= 8)
}

// --- MCPBackend implementation ---

func (b *SSHBackend) ListJobs(ctx context.Context) (string, error) {
	if b.hasJSON {
		return b.runCommand(ctx, "squeue --json")
	}
	output, err := b.runCommand(ctx, `squeue -o "%A|%j|%T|%u|%P|%M|%l|%D|%R" --noheader`)
	if err != nil {
		return "", err
	}
	return parseDelimited(output, []string{"JobID", "Name", "State", "User", "Partition", "Time", "TimeLimit", "Nodes", "Reason"})
}

func (b *SSHBackend) GetJob(ctx context.Context, jobID string) (string, error) {
	if b.hasJSON {
		return b.runCommand(ctx, fmt.Sprintf("scontrol show job %s --json", sanitizeID(jobID)))
	}
	output, err := b.runCommand(ctx, fmt.Sprintf("scontrol show job %s", sanitizeID(jobID)))
	if err != nil {
		return "", err
	}
	return parseScontrolShow(output)
}

func (b *SSHBackend) SubmitJob(ctx context.Context, script string) (string, error) {
	output, err := b.runCommandWithStdin(ctx, "sbatch --parsable", script)
	if err != nil {
		return "", err
	}
	jobID := strings.TrimSpace(output)
	result := map[string]string{"job_id": jobID}
	data, err := json.Marshal(result)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func (b *SSHBackend) CancelJob(ctx context.Context, jobID string) error {
	_, err := b.runCommand(ctx, fmt.Sprintf("scancel %s", sanitizeID(jobID)))
	return err
}

func (b *SSHBackend) ListNodes(ctx context.Context) (string, error) {
	if b.hasJSON {
		return b.runCommand(ctx, "sinfo --json -N")
	}
	output, err := b.runCommand(ctx, `sinfo -N -o "%N|%T|%P|%c|%m|%G" --noheader`)
	if err != nil {
		return "", err
	}
	return parseDelimited(output, []string{"NodeName", "State", "Partition", "CPUs", "Memory", "GRES"})
}

func (b *SSHBackend) GetNode(ctx context.Context, nodeName string) (string, error) {
	if b.hasJSON {
		return b.runCommand(ctx, fmt.Sprintf("scontrol show node %s --json", sanitizeName(nodeName)))
	}
	output, err := b.runCommand(ctx, fmt.Sprintf("scontrol show node %s", sanitizeName(nodeName)))
	if err != nil {
		return "", err
	}
	return parseScontrolShow(output)
}

func (b *SSHBackend) ListPartitions(ctx context.Context) (string, error) {
	if b.hasJSON {
		return b.runCommand(ctx, "sinfo --json -s")
	}
	output, err := b.runCommand(ctx, `sinfo -o "%P|%a|%l|%D|%T|%c|%m" --noheader`)
	if err != nil {
		return "", err
	}
	return parseDelimited(output, []string{"Partition", "Avail", "TimeLimit", "Nodes", "State", "CPUs", "Memory"})
}

func (b *SSHBackend) GetPartition(ctx context.Context, partitionName string) (string, error) {
	if b.hasJSON {
		return b.runCommand(ctx, fmt.Sprintf("scontrol show partition %s --json", sanitizeName(partitionName)))
	}
	output, err := b.runCommand(ctx, fmt.Sprintf("scontrol show partition %s", sanitizeName(partitionName)))
	if err != nil {
		return "", err
	}
	return parseScontrolShow(output)
}

func (b *SSHBackend) ListReservations(ctx context.Context) (string, error) {
	if b.hasJSON {
		return b.runCommand(ctx, "scontrol show reservation --json")
	}
	output, err := b.runCommand(ctx, "scontrol show reservation")
	if err != nil {
		return "", err
	}
	return parseScontrolShowMulti(output)
}

func (b *SSHBackend) ListLicenses(ctx context.Context) (string, error) {
	if b.hasJSON {
		return b.runCommand(ctx, "scontrol show license --json")
	}
	output, err := b.runCommand(ctx, "scontrol show license")
	if err != nil {
		return "", err
	}
	return parseScontrolShowMulti(output)
}

func (b *SSHBackend) GetDiag(ctx context.Context) (string, error) {
	output, err := b.runCommand(ctx, "sdiag")
	if err != nil {
		return "", err
	}
	return parseSdiag(output)
}

func (b *SSHBackend) Ping(ctx context.Context) (string, error) {
	output, err := b.runCommand(ctx, "scontrol ping")
	if err != nil {
		return "", err
	}
	alive := strings.Contains(output, "is alive")
	result := map[string]any{
		"status":  "ok",
		"alive":   alive,
		"message": strings.TrimSpace(output),
	}
	if !alive {
		result["status"] = "error"
	}
	data, err := json.Marshal(result)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// --- SlurmDB operations (via sacctmgr) ---

func (b *SSHBackend) ListAccounts(ctx context.Context) (string, error) {
	output, err := b.runCommand(ctx, `sacctmgr show account -P -n format=Account,Description,Organization`)
	if err != nil {
		return "", err
	}
	return parseDelimited(output, []string{"Account", "Description", "Organization"})
}

func (b *SSHBackend) GetAccount(ctx context.Context, name string) (string, error) {
	output, err := b.runCommand(ctx, fmt.Sprintf(`sacctmgr show account %s -P -n format=Account,Description,Organization`, sanitizeName(name)))
	if err != nil {
		return "", err
	}
	return parseDelimited(output, []string{"Account", "Description", "Organization"})
}

func (b *SSHBackend) ListUsers(ctx context.Context) (string, error) {
	output, err := b.runCommand(ctx, `sacctmgr show user -P -n format=User,DefaultAccount,AdminLevel`)
	if err != nil {
		return "", err
	}
	return parseDelimited(output, []string{"User", "DefaultAccount", "AdminLevel"})
}

func (b *SSHBackend) GetUser(ctx context.Context, name string) (string, error) {
	output, err := b.runCommand(ctx, fmt.Sprintf(`sacctmgr show user %s -P -n format=User,DefaultAccount,AdminLevel`, sanitizeName(name)))
	if err != nil {
		return "", err
	}
	return parseDelimited(output, []string{"User", "DefaultAccount", "AdminLevel"})
}

func (b *SSHBackend) ListSlurmdbJobs(ctx context.Context) (string, error) {
	output, err := b.runCommand(ctx, `sacct -P -n --format=JobID,JobName,State,User,Partition,Elapsed,Start,End,ExitCode`)
	if err != nil {
		return "", err
	}
	return parseDelimited(output, []string{"JobID", "JobName", "State", "User", "Partition", "Elapsed", "Start", "End", "ExitCode"})
}

func (b *SSHBackend) GetSlurmdbJob(ctx context.Context, jobID string) (string, error) {
	output, err := b.runCommand(ctx, fmt.Sprintf(`sacct -P -n -j %s --format=JobID,JobName,State,User,Partition,Elapsed,Start,End,ExitCode,MaxRSS,MaxVMSize`, sanitizeID(jobID)))
	if err != nil {
		return "", err
	}
	return parseDelimited(output, []string{"JobID", "JobName", "State", "User", "Partition", "Elapsed", "Start", "End", "ExitCode", "MaxRSS", "MaxVMSize"})
}

func (b *SSHBackend) ListClusters(ctx context.Context) (string, error) {
	output, err := b.runCommand(ctx, `sacctmgr show cluster -P -n format=Cluster,ControlHost,ControlPort,RPC`)
	if err != nil {
		return "", err
	}
	return parseDelimited(output, []string{"Cluster", "ControlHost", "ControlPort", "RPC"})
}

func (b *SSHBackend) ListQOS(ctx context.Context) (string, error) {
	output, err := b.runCommand(ctx, `sacctmgr show qos -P -n format=Name,Priority,GrpTRES,MaxTRES,MaxWall`)
	if err != nil {
		return "", err
	}
	return parseDelimited(output, []string{"Name", "Priority", "GrpTRES", "MaxTRES", "MaxWall"})
}

func (b *SSHBackend) ListAssociations(ctx context.Context) (string, error) {
	output, err := b.runCommand(ctx, `sacctmgr show association -P -n format=Cluster,Account,User,Partition,QOS`)
	if err != nil {
		return "", err
	}
	return parseDelimited(output, []string{"Cluster", "Account", "User", "Partition", "QOS"})
}

func (b *SSHBackend) ListTRES(ctx context.Context) (string, error) {
	output, err := b.runCommand(ctx, `sacctmgr show tres -P -n format=Type,Name,ID`)
	if err != nil {
		return "", err
	}
	return parseDelimited(output, []string{"Type", "Name", "ID"})
}

func (b *SSHBackend) ListWCKeys(ctx context.Context) (string, error) {
	output, err := b.runCommand(ctx, `sacctmgr show wckey -P -n format=WCKey,Cluster,User`)
	if err != nil {
		return "", err
	}
	return parseDelimited(output, []string{"WCKey", "Cluster", "User"})
}

func (b *SSHBackend) GetSlurmdbDiag(ctx context.Context) (string, error) {
	output, err := b.runCommand(ctx, "sdiag")
	if err != nil {
		return "", err
	}
	return parseSdiag(output)
}

// --- Parsing helpers ---

// parseDelimited parses pipe-delimited output into a JSON array of objects.
func parseDelimited(output string, fields []string) (string, error) {
	var results []map[string]string
	for _, line := range strings.Split(strings.TrimSpace(output), "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		parts := strings.Split(line, "|")
		if len(parts) != len(fields) {
			continue // skip malformed lines
		}
		row := make(map[string]string)
		for i, f := range fields {
			row[f] = strings.TrimSpace(parts[i])
		}
		results = append(results, row)
	}
	if results == nil {
		results = []map[string]string{}
	}
	data, err := json.Marshal(results)
	return string(data), err
}

// parseScontrolShow parses a single scontrol show record (key=value pairs).
func parseScontrolShow(output string) (string, error) {
	result := parseScontrolRecord(output)
	data, err := json.Marshal(result)
	return string(data), err
}

// parseScontrolShowMulti parses multiple scontrol show records separated by blank lines.
func parseScontrolShowMulti(output string) (string, error) {
	blocks := splitScontrolBlocks(output)
	var results []map[string]string
	for _, block := range blocks {
		record := parseScontrolRecord(block)
		if len(record) > 0 {
			results = append(results, record)
		}
	}
	if results == nil {
		results = []map[string]string{}
	}
	data, err := json.Marshal(results)
	return string(data), err
}

// parseScontrolRecord parses a single scontrol record text block.
// scontrol output format: "Key1=Value1 Key2=Value2\n   Key3=Value3"
// Values can contain spaces if they are part of a path or description,
// so we split on key=value boundaries.
func parseScontrolRecord(text string) map[string]string {
	result := make(map[string]string)

	// Normalize whitespace: join lines, collapse multiple spaces
	text = strings.ReplaceAll(text, "\n", " ")
	text = strings.TrimSpace(text)
	if text == "" {
		return result
	}

	// Match Key=Value pairs. Keys are alphanumeric (plus underscores).
	// A new key starts where we see " Key=" pattern.
	re := regexp.MustCompile(`([A-Za-z_][A-Za-z0-9_]*)=`)
	indices := re.FindAllStringIndex(text, -1)

	for i, idx := range indices {
		keyMatch := re.FindStringSubmatch(text[idx[0]:idx[1]])
		key := keyMatch[1]
		valueStart := idx[1]
		var valueEnd int
		if i+1 < len(indices) {
			// Value extends to just before the next key's space separator
			valueEnd = indices[i+1][0]
			// Trim trailing space before next key
			for valueEnd > valueStart && text[valueEnd-1] == ' ' {
				valueEnd--
			}
		} else {
			valueEnd = len(text)
		}
		result[key] = strings.TrimSpace(text[valueStart:valueEnd])
	}

	return result
}

// splitScontrolBlocks splits scontrol output into blocks separated by blank lines.
func splitScontrolBlocks(output string) []string {
	var blocks []string
	var current strings.Builder
	for _, line := range strings.Split(output, "\n") {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" {
			if current.Len() > 0 {
				blocks = append(blocks, current.String())
				current.Reset()
			}
			continue
		}
		if current.Len() > 0 {
			current.WriteByte('\n')
		}
		current.WriteString(line)
	}
	if current.Len() > 0 {
		blocks = append(blocks, current.String())
	}
	return blocks
}

// parseSdiag parses sdiag text output into a JSON object.
func parseSdiag(output string) (string, error) {
	result := make(map[string]string)
	for _, line := range strings.Split(output, "\n") {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "*") {
			continue
		}
		if idx := strings.Index(line, ":"); idx > 0 {
			key := strings.TrimSpace(line[:idx])
			value := strings.TrimSpace(line[idx+1:])
			// Normalize key: replace spaces with underscores
			key = strings.ReplaceAll(key, " ", "_")
			result[key] = value
		}
	}
	if len(result) == 0 {
		result["raw"] = strings.TrimSpace(output)
	}
	data, err := json.Marshal(result)
	return string(data), err
}

// --- Input sanitization ---

// sanitizeID ensures a job ID contains only digits, letters, periods, and underscores
// (for array job IDs like "12345_1" and step names like "12345.batch").
func sanitizeID(id string) string {
	re := regexp.MustCompile(`[^a-zA-Z0-9_.]`)
	return re.ReplaceAllString(id, "")
}

// sanitizeName ensures a node/partition name contains only safe characters.
func sanitizeName(name string) string {
	re := regexp.MustCompile(`[^a-zA-Z0-9_\-.]`)
	return re.ReplaceAllString(name, "")
}

// Compile-time interface check.
var _ MCPBackend = (*SSHBackend)(nil)

// Close closes the SSH connection.
func (b *SSHBackend) Close() error {
	b.mu.Lock()
	defer b.mu.Unlock()
	if b.client != nil {
		return b.client.Close()
	}
	return nil
}
