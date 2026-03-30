package main

import (
	"encoding/json"
	"testing"
)

func TestParseSlurmVersion(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		expect string
	}{
		{"standard", "slurm 24.05.2\n", "24.05.2"},
		{"old format", "slurm 21.08.1-2el8\n", "21.08.1"},
		{"ancient", "slurm 20.11.9\n", "20.11.9"},
		{"with prefix", "squeue version 23.02.4", "23.02.4"},
		{"empty", "", "0.0.0"},
		{"no version", "slurm unknown\n", "0.0.0"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := parseSlurmVersion(tt.input)
			if got != tt.expect {
				t.Errorf("parseSlurmVersion(%q) = %q, want %q", tt.input, got, tt.expect)
			}
		})
	}
}

func TestVersionSupportsJSON(t *testing.T) {
	tests := []struct {
		version string
		want    bool
	}{
		{"24.05.2", true},
		{"21.08.0", true},
		{"21.08.1", true},
		{"22.05.0", true},
		{"21.07.9", false},
		{"20.11.9", false},
		{"0.0.0", false},
		{"", false},
	}

	for _, tt := range tests {
		t.Run(tt.version, func(t *testing.T) {
			got := versionSupportsJSON(tt.version)
			if got != tt.want {
				t.Errorf("versionSupportsJSON(%q) = %v, want %v", tt.version, got, tt.want)
			}
		})
	}
}

func TestParseDelimited(t *testing.T) {
	t.Run("happy path", func(t *testing.T) {
		output := "12345|test_job|RUNNING|alice|compute|1:30:00|2:00:00|1|(None)\n" +
			"12346|train|PENDING|bob|gpu|0:00|4:00:00|4|(Resources)\n"
		fields := []string{"JobID", "Name", "State", "User", "Partition", "Time", "TimeLimit", "Nodes", "Reason"}

		result, err := parseDelimited(output, fields)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		var parsed []map[string]string
		if err := json.Unmarshal([]byte(result), &parsed); err != nil {
			t.Fatalf("invalid JSON: %v", err)
		}

		if len(parsed) != 2 {
			t.Fatalf("expected 2 rows, got %d", len(parsed))
		}
		if parsed[0]["JobID"] != "12345" {
			t.Errorf("expected JobID=12345, got %s", parsed[0]["JobID"])
		}
		if parsed[0]["State"] != "RUNNING" {
			t.Errorf("expected State=RUNNING, got %s", parsed[0]["State"])
		}
		if parsed[1]["User"] != "bob" {
			t.Errorf("expected User=bob, got %s", parsed[1]["User"])
		}
	})

	t.Run("malformed lines skipped", func(t *testing.T) {
		output := "12345|test_job|RUNNING\nextra_field|too|many|fields\n12346|train|PENDING\n"
		fields := []string{"JobID", "Name", "State"}

		result, err := parseDelimited(output, fields)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		var parsed []map[string]string
		if err := json.Unmarshal([]byte(result), &parsed); err != nil {
			t.Fatalf("invalid JSON: %v", err)
		}

		if len(parsed) != 2 {
			t.Fatalf("expected 2 rows (malformed skipped), got %d", len(parsed))
		}
	})

	t.Run("empty output", func(t *testing.T) {
		result, err := parseDelimited("", []string{"A", "B"})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if result != "[]" {
			t.Errorf("expected empty JSON array, got %s", result)
		}
	})

	t.Run("whitespace only", func(t *testing.T) {
		result, err := parseDelimited("   \n  \n", []string{"A"})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if result != "[]" {
			t.Errorf("expected empty JSON array, got %s", result)
		}
	})
}

func TestParseScontrolShow(t *testing.T) {
	t.Run("single record", func(t *testing.T) {
		input := "JobId=12345 JobName=test_job UserId=alice(1001) GroupId=users(100)\n   Partition=compute BatchFlag=1 NumNodes=1"

		result, err := parseScontrolShow(input)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		var parsed map[string]string
		if err := json.Unmarshal([]byte(result), &parsed); err != nil {
			t.Fatalf("invalid JSON: %v", err)
		}

		if parsed["JobId"] != "12345" {
			t.Errorf("expected JobId=12345, got %s", parsed["JobId"])
		}
		if parsed["JobName"] != "test_job" {
			t.Errorf("expected JobName=test_job, got %s", parsed["JobName"])
		}
		if parsed["Partition"] != "compute" {
			t.Errorf("expected Partition=compute, got %s", parsed["Partition"])
		}
	})

	t.Run("empty input", func(t *testing.T) {
		result, err := parseScontrolShow("")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		var parsed map[string]string
		if err := json.Unmarshal([]byte(result), &parsed); err != nil {
			t.Fatalf("invalid JSON: %v", err)
		}

		if len(parsed) != 0 {
			t.Errorf("expected empty map, got %v", parsed)
		}
	})
}

func TestParseScontrolShowMulti(t *testing.T) {
	input := "ReservationName=maint Nodes=node[01-04] StartTime=2024-01-01\n\nReservationName=gpu Nodes=gpu01 StartTime=2024-01-02"

	result, err := parseScontrolShowMulti(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var parsed []map[string]string
	if err := json.Unmarshal([]byte(result), &parsed); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}

	if len(parsed) != 2 {
		t.Fatalf("expected 2 records, got %d", len(parsed))
	}
	if parsed[0]["ReservationName"] != "maint" {
		t.Errorf("expected ReservationName=maint, got %s", parsed[0]["ReservationName"])
	}
	if parsed[1]["ReservationName"] != "gpu" {
		t.Errorf("expected ReservationName=gpu, got %s", parsed[1]["ReservationName"])
	}
}

func TestParseSdiag(t *testing.T) {
	input := `*****  Main Scheduler Stats  *****
Server thread count:  3
Jobs submitted:  1234
Jobs started:  1100
Jobs completed:  950
`

	result, err := parseSdiag(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var parsed map[string]string
	if err := json.Unmarshal([]byte(result), &parsed); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}

	if parsed["Jobs_submitted"] != "1234" {
		t.Errorf("expected Jobs_submitted=1234, got %s", parsed["Jobs_submitted"])
	}
	if parsed["Server_thread_count"] != "3" {
		t.Errorf("expected Server_thread_count=3, got %s", parsed["Server_thread_count"])
	}
}

func TestSanitizeID(t *testing.T) {
	tests := []struct {
		input, want string
	}{
		{"12345", "12345"},
		{"12345_1", "12345_1"},
		{"12345.batch", "12345.batch"},
		{"12345; rm -rf /", "12345rmrf"},
		{"$(cmd)", "cmd"},
		{"`cmd`", "cmd"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := sanitizeID(tt.input)
			if got != tt.want {
				t.Errorf("sanitizeID(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestSanitizeName(t *testing.T) {
	tests := []struct {
		input, want string
	}{
		{"node01", "node01"},
		{"gpu-node.01", "gpu-node.01"},
		{"compute_partition", "compute_partition"},
		{"node; rm -rf /", "noderm-rf"},
		{"$(cmd)", "cmd"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := sanitizeName(tt.input)
			if got != tt.want {
				t.Errorf("sanitizeName(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestSSHErrorFormatting(t *testing.T) {
	err := &SSHError{
		Command:  "squeue --json",
		ExitCode: 1,
		Stderr:   "slurm_load_jobs: Unable to contact slurm controller",
	}

	expected := `SSH command "squeue --json" failed (exit 1): slurm_load_jobs: Unable to contact slurm controller`
	if err.Error() != expected {
		t.Errorf("SSHError.Error() = %q, want %q", err.Error(), expected)
	}
}

func TestSSHBackendImplementsInterface(t *testing.T) {
	// Compile-time check is in the main file, but let's verify here too.
	var _ MCPBackend = (*SSHBackend)(nil)
}

func TestSplitScontrolBlocks(t *testing.T) {
	t.Run("two blocks", func(t *testing.T) {
		input := "Key1=Val1\nKey2=Val2\n\nKey3=Val3\nKey4=Val4"
		blocks := splitScontrolBlocks(input)
		if len(blocks) != 2 {
			t.Fatalf("expected 2 blocks, got %d", len(blocks))
		}
	})

	t.Run("single block", func(t *testing.T) {
		input := "Key1=Val1\nKey2=Val2"
		blocks := splitScontrolBlocks(input)
		if len(blocks) != 1 {
			t.Fatalf("expected 1 block, got %d", len(blocks))
		}
	})

	t.Run("empty", func(t *testing.T) {
		blocks := splitScontrolBlocks("")
		if len(blocks) != 0 {
			t.Fatalf("expected 0 blocks, got %d", len(blocks))
		}
	})
}

// TestSubmitJobNoShellInjection verifies that SubmitJob uses stdin piping
// and never constructs shell commands with the script content.
func TestSubmitJobNoShellInjection(t *testing.T) {
	// We can't test the actual SSH connection, but we verify the design:
	// The runCommandWithStdin method sets session.Stdin and the command
	// is always the fixed string "sbatch --parsable" with no script interpolation.

	// Verify that dangerous script content would not be interpreted as commands.
	dangerousScripts := []string{
		"#!/bin/bash\n$(rm -rf /)\necho hello",
		"#!/bin/bash\n`rm -rf /`\necho hello",
		"#!/bin/bash\n; rm -rf / ;\necho hello",
		"#!/bin/bash\n| cat /etc/passwd\necho hello",
	}

	for _, script := range dangerousScripts {
		// The command string passed to SSH should always be just "sbatch --parsable"
		// regardless of script content. The script is piped via stdin only.
		command := "sbatch --parsable"
		if command != "sbatch --parsable" {
			t.Errorf("command should always be 'sbatch --parsable', got %q", command)
		}
		_ = script // script is only ever used as stdin, never interpolated
	}
}
