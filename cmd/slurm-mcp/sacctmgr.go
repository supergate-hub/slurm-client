package main

import (
	"context"
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
)

// SacctmgrRunner executes a sacctmgr/sacct command over SSH and returns stdout.
type SacctmgrRunner func(ctx context.Context, cmd string) (string, error)

var validNameRe = regexp.MustCompile(`^[a-zA-Z0-9_-]+$`)

func validateName(name string) error {
	if !validNameRe.MatchString(name) {
		return fmt.Errorf("invalid name: %q (only alphanumeric, hyphens, underscores allowed)", name)
	}
	return nil
}

// parseParsable2 parses sacctmgr/sacct --parsable2 output (pipe-delimited with header row)
// into a JSON array of objects.
func parseParsable2(output string) (string, error) {
	lines := strings.Split(strings.TrimSpace(output), "\n")
	if len(lines) < 1 || lines[0] == "" {
		return "[]", nil
	}

	headers := strings.Split(lines[0], "|")

	var results []map[string]string
	for _, line := range lines[1:] {
		if line == "" {
			continue
		}
		fields := strings.Split(line, "|")
		row := make(map[string]string)
		for i, h := range headers {
			if i < len(fields) {
				row[h] = fields[i]
			}
		}
		results = append(results, row)
	}

	if results == nil {
		return "[]", nil
	}
	data, err := json.Marshal(results)
	if err != nil {
		return "", fmt.Errorf("json marshal error: %w", err)
	}
	return string(data), nil
}

// parseDiagStats parses "sacctmgr show stats" text output (key=value lines).
func parseDiagStats(output string) (string, error) {
	result := make(map[string]string)
	for _, line := range strings.Split(strings.TrimSpace(output), "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		if idx := strings.IndexByte(line, '='); idx >= 0 {
			key := strings.TrimSpace(line[:idx])
			val := strings.TrimSpace(line[idx+1:])
			result[key] = val
		}
	}
	data, err := json.Marshal(result)
	if err != nil {
		return "", fmt.Errorf("json marshal error: %w", err)
	}
	return string(data), nil
}

// --- sacctmgr commands ---

func sacctmgrListAccounts(ctx context.Context, run SacctmgrRunner) (string, error) {
	out, err := run(ctx, "sacctmgr list account --parsable2")
	if err != nil {
		return "", fmt.Errorf("sacctmgr list account: %w", err)
	}
	return parseParsable2(out)
}

func sacctmgrGetAccount(ctx context.Context, run SacctmgrRunner, name string) (string, error) {
	if err := validateName(name); err != nil {
		return "", err
	}
	out, err := run(ctx, fmt.Sprintf("sacctmgr list account where name=%s --parsable2", name))
	if err != nil {
		return "", fmt.Errorf("sacctmgr list account where name=%s: %w", name, err)
	}
	return parseParsable2(out)
}

func sacctmgrListUsers(ctx context.Context, run SacctmgrRunner) (string, error) {
	out, err := run(ctx, "sacctmgr list user --parsable2")
	if err != nil {
		return "", fmt.Errorf("sacctmgr list user: %w", err)
	}
	return parseParsable2(out)
}

func sacctmgrGetUser(ctx context.Context, run SacctmgrRunner, name string) (string, error) {
	if err := validateName(name); err != nil {
		return "", err
	}
	out, err := run(ctx, fmt.Sprintf("sacctmgr list user where name=%s --parsable2", name))
	if err != nil {
		return "", fmt.Errorf("sacctmgr list user where name=%s: %w", name, err)
	}
	return parseParsable2(out)
}

func sacctmgrListClusters(ctx context.Context, run SacctmgrRunner) (string, error) {
	out, err := run(ctx, "sacctmgr list cluster --parsable2")
	if err != nil {
		return "", fmt.Errorf("sacctmgr list cluster: %w", err)
	}
	return parseParsable2(out)
}

func sacctmgrListQOS(ctx context.Context, run SacctmgrRunner) (string, error) {
	out, err := run(ctx, "sacctmgr list qos --parsable2")
	if err != nil {
		return "", fmt.Errorf("sacctmgr list qos: %w", err)
	}
	return parseParsable2(out)
}

func sacctmgrListAssociations(ctx context.Context, run SacctmgrRunner) (string, error) {
	out, err := run(ctx, "sacctmgr list assoc --parsable2")
	if err != nil {
		return "", fmt.Errorf("sacctmgr list assoc: %w", err)
	}
	return parseParsable2(out)
}

func sacctmgrListTRES(ctx context.Context, run SacctmgrRunner) (string, error) {
	out, err := run(ctx, "sacctmgr list tres --parsable2")
	if err != nil {
		return "", fmt.Errorf("sacctmgr list tres: %w", err)
	}
	return parseParsable2(out)
}

func sacctmgrListWCKeys(ctx context.Context, run SacctmgrRunner) (string, error) {
	out, err := run(ctx, "sacctmgr list wckey --parsable2")
	if err != nil {
		return "", fmt.Errorf("sacctmgr list wckey: %w", err)
	}
	return parseParsable2(out)
}

// --- sacct commands (SlurmDB jobs) ---

const sacctDefaultFormat = "JobID,JobName,User,State,ExitCode,Start,End,Elapsed,Partition,NodeList"

func sacctListJobs(ctx context.Context, run SacctmgrRunner) (string, error) {
	out, err := run(ctx, fmt.Sprintf("sacct --parsable2 --format=%s", sacctDefaultFormat))
	if err != nil {
		return "", fmt.Errorf("sacct: %w", err)
	}
	return parseParsable2(out)
}

func sacctGetJob(ctx context.Context, run SacctmgrRunner, jobID string) (string, error) {
	if err := validateName(jobID); err != nil {
		return "", err
	}
	out, err := run(ctx, fmt.Sprintf("sacct --parsable2 --format=%s -j %s", sacctDefaultFormat, jobID))
	if err != nil {
		return "", fmt.Errorf("sacct -j %s: %w", jobID, err)
	}
	return parseParsable2(out)
}

// --- sacctmgr diag and config ---

func sacctmgrGetDiag(ctx context.Context, run SacctmgrRunner) (string, error) {
	out, err := run(ctx, "sacctmgr show stats")
	if err != nil {
		return "", fmt.Errorf("sacctmgr show stats: %w", err)
	}
	return parseDiagStats(out)
}

func sacctmgrGetConfig(ctx context.Context, run SacctmgrRunner) (string, error) {
	out, err := run(ctx, "sacctmgr show config --parsable2")
	if err != nil {
		return "", fmt.Errorf("sacctmgr show config: %w", err)
	}
	return parseParsable2(out)
}
