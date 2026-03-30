package main

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"
)

// --- parseParsable2 tests ---

func TestParseParsable2_HappyPath(t *testing.T) {
	input := "Account|Description|Organization\nroot|default root account|root\nresearch|Research group|university"
	got, err := parseParsable2(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var records []map[string]string
	if err := json.Unmarshal([]byte(got), &records); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if len(records) != 2 {
		t.Fatalf("expected 2 records, got %d", len(records))
	}
	if records[0]["Account"] != "root" {
		t.Errorf("expected Account=root, got %q", records[0]["Account"])
	}
	if records[1]["Organization"] != "university" {
		t.Errorf("expected Organization=university, got %q", records[1]["Organization"])
	}
}

func TestParseParsable2_Empty(t *testing.T) {
	got, err := parseParsable2("")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "[]" {
		t.Errorf("expected [], got %q", got)
	}
}

func TestParseParsable2_HeaderOnly(t *testing.T) {
	got, err := parseParsable2("Account|Description|Organization")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "[]" {
		t.Errorf("expected [], got %q", got)
	}
}

func TestParseParsable2_FewerFieldsThanHeaders(t *testing.T) {
	input := "A|B|C\n1|2"
	got, err := parseParsable2(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var records []map[string]string
	if err := json.Unmarshal([]byte(got), &records); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if len(records) != 1 {
		t.Fatalf("expected 1 record, got %d", len(records))
	}
	if records[0]["A"] != "1" {
		t.Errorf("expected A=1, got %q", records[0]["A"])
	}
	if records[0]["B"] != "2" {
		t.Errorf("expected B=2, got %q", records[0]["B"])
	}
	// C should not be present since there were fewer fields
	if _, ok := records[0]["C"]; ok {
		t.Errorf("expected C to be absent")
	}
}

func TestParseParsable2_TrailingNewlines(t *testing.T) {
	input := "Name|Value\nfoo|bar\n\n"
	got, err := parseParsable2(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var records []map[string]string
	if err := json.Unmarshal([]byte(got), &records); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if len(records) != 1 {
		t.Fatalf("expected 1 record, got %d", len(records))
	}
}

// --- parseDiagStats tests ---

func TestParseDiagStats(t *testing.T) {
	input := "RPC Count = 42\nThread Count = 8\n"
	got, err := parseDiagStats(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var m map[string]string
	if err := json.Unmarshal([]byte(got), &m); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if m["RPC Count"] != "42" {
		t.Errorf("expected RPC Count=42, got %q", m["RPC Count"])
	}
	if m["Thread Count"] != "8" {
		t.Errorf("expected Thread Count=8, got %q", m["Thread Count"])
	}
}

func TestParseDiagStats_Empty(t *testing.T) {
	got, err := parseDiagStats("")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "{}" {
		t.Errorf("expected {}, got %q", got)
	}
}

// --- validateName tests ---

func TestValidateName_Valid(t *testing.T) {
	valid := []string{"root", "user-1", "research_group", "ABC123", "a-b_c"}
	for _, name := range valid {
		if err := validateName(name); err != nil {
			t.Errorf("expected %q to be valid, got error: %v", name, err)
		}
	}
}

func TestValidateName_Invalid(t *testing.T) {
	invalid := []string{
		"",
		"user name",
		"user;rm -rf",
		"$(whoami)",
		"user|cat",
		"name&cmd",
		"a`b`",
		"user\nname",
		"../etc",
	}
	for _, name := range invalid {
		if err := validateName(name); err == nil {
			t.Errorf("expected %q to be invalid, got nil error", name)
		}
	}
}

// --- mock runner helper ---

func mockRunner(output string, err error) SacctmgrRunner {
	return func(ctx context.Context, cmd string) (string, error) {
		return output, err
	}
}

// captureRunner records the command that was executed.
func captureRunner(output string) (SacctmgrRunner, *string) {
	var captured string
	runner := func(ctx context.Context, cmd string) (string, error) {
		captured = cmd
		return output, nil
	}
	return runner, &captured
}

// --- sacctmgr function tests ---

func TestSacctmgrListAccounts(t *testing.T) {
	run, cmd := captureRunner("Account|Description\nroot|default")
	got, err := sacctmgrListAccounts(context.Background(), run)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if *cmd != "sacctmgr list account --parsable2" {
		t.Errorf("unexpected command: %q", *cmd)
	}
	var records []map[string]string
	if err := json.Unmarshal([]byte(got), &records); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if len(records) != 1 || records[0]["Account"] != "root" {
		t.Errorf("unexpected result: %s", got)
	}
}

func TestSacctmgrGetAccount(t *testing.T) {
	run, cmd := captureRunner("Account|Description\nresearch|Research group")
	got, err := sacctmgrGetAccount(context.Background(), run, "research")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if *cmd != "sacctmgr list account where name=research --parsable2" {
		t.Errorf("unexpected command: %q", *cmd)
	}
	var records []map[string]string
	if err := json.Unmarshal([]byte(got), &records); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if len(records) != 1 || records[0]["Account"] != "research" {
		t.Errorf("unexpected result: %s", got)
	}
}

func TestSacctmgrGetAccount_Injection(t *testing.T) {
	run := mockRunner("", nil)
	_, err := sacctmgrGetAccount(context.Background(), run, "name;rm -rf /")
	if err == nil {
		t.Fatal("expected validation error for injection attempt")
	}
}

func TestSacctmgrListUsers(t *testing.T) {
	run, cmd := captureRunner("User|AdminLevel\njohn|None")
	got, err := sacctmgrListUsers(context.Background(), run)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if *cmd != "sacctmgr list user --parsable2" {
		t.Errorf("unexpected command: %q", *cmd)
	}
	var records []map[string]string
	if err := json.Unmarshal([]byte(got), &records); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if len(records) != 1 || records[0]["User"] != "john" {
		t.Errorf("unexpected result: %s", got)
	}
}

func TestSacctmgrGetUser_Injection(t *testing.T) {
	run := mockRunner("", nil)
	_, err := sacctmgrGetUser(context.Background(), run, "$(whoami)")
	if err == nil {
		t.Fatal("expected validation error for injection attempt")
	}
}

func TestSacctmgrListClusters(t *testing.T) {
	run, cmd := captureRunner("Cluster|ControlHost\nmycluster|10.0.0.1")
	_, err := sacctmgrListClusters(context.Background(), run)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if *cmd != "sacctmgr list cluster --parsable2" {
		t.Errorf("unexpected command: %q", *cmd)
	}
}

func TestSacctmgrListQOS(t *testing.T) {
	run, cmd := captureRunner("Name|Priority\nnormal|10")
	_, err := sacctmgrListQOS(context.Background(), run)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if *cmd != "sacctmgr list qos --parsable2" {
		t.Errorf("unexpected command: %q", *cmd)
	}
}

func TestSacctmgrListAssociations(t *testing.T) {
	run, cmd := captureRunner("Cluster|Account|User\nmycluster|root|")
	_, err := sacctmgrListAssociations(context.Background(), run)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if *cmd != "sacctmgr list assoc --parsable2" {
		t.Errorf("unexpected command: %q", *cmd)
	}
}

func TestSacctmgrListTRES(t *testing.T) {
	run, cmd := captureRunner("Type|Name|ID\ncpu||1")
	_, err := sacctmgrListTRES(context.Background(), run)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if *cmd != "sacctmgr list tres --parsable2" {
		t.Errorf("unexpected command: %q", *cmd)
	}
}

func TestSacctmgrListWCKeys(t *testing.T) {
	run, cmd := captureRunner("WCKey|Cluster|User\ndefault|mycluster|john")
	_, err := sacctmgrListWCKeys(context.Background(), run)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if *cmd != "sacctmgr list wckey --parsable2" {
		t.Errorf("unexpected command: %q", *cmd)
	}
}

// --- sacct function tests ---

func TestSacctListJobs(t *testing.T) {
	run, cmd := captureRunner("JobID|JobName|User|State|ExitCode|Start|End|Elapsed|Partition|NodeList\n123|test|john|COMPLETED|0:0|2024-01-01|2024-01-02|01:00:00|batch|node1")
	got, err := sacctListJobs(context.Background(), run)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if *cmd != fmt.Sprintf("sacct --parsable2 --format=%s", sacctDefaultFormat) {
		t.Errorf("unexpected command: %q", *cmd)
	}
	var records []map[string]string
	if err := json.Unmarshal([]byte(got), &records); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if len(records) != 1 || records[0]["JobID"] != "123" {
		t.Errorf("unexpected result: %s", got)
	}
}

func TestSacctGetJob(t *testing.T) {
	run, cmd := captureRunner("JobID|JobName\n456|myjob")
	_, err := sacctGetJob(context.Background(), run, "456")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if *cmd != fmt.Sprintf("sacct --parsable2 --format=%s -j 456", sacctDefaultFormat) {
		t.Errorf("unexpected command: %q", *cmd)
	}
}

func TestSacctGetJob_Injection(t *testing.T) {
	run := mockRunner("", nil)
	_, err := sacctGetJob(context.Background(), run, "123;rm -rf /")
	if err == nil {
		t.Fatal("expected validation error for injection attempt")
	}
}

// --- sacctmgr diag/config tests ---

func TestSacctmgrGetDiag(t *testing.T) {
	run, cmd := captureRunner("RPC Count = 42\nThread Count = 8")
	got, err := sacctmgrGetDiag(context.Background(), run)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if *cmd != "sacctmgr show stats" {
		t.Errorf("unexpected command: %q", *cmd)
	}
	var m map[string]string
	if err := json.Unmarshal([]byte(got), &m); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if m["RPC Count"] != "42" {
		t.Errorf("unexpected result: %s", got)
	}
}

func TestSacctmgrGetConfig(t *testing.T) {
	run, cmd := captureRunner("Parameter|Value\nStorageType|accounting_storage/mysql")
	_, err := sacctmgrGetConfig(context.Background(), run)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if *cmd != "sacctmgr show config --parsable2" {
		t.Errorf("unexpected command: %q", *cmd)
	}
}

// --- runner error propagation ---

func TestRunnerError(t *testing.T) {
	run := mockRunner("", fmt.Errorf("ssh connection failed"))
	_, err := sacctmgrListAccounts(context.Background(), run)
	if err == nil {
		t.Fatal("expected error to propagate")
	}
	if !contains(err.Error(), "ssh connection failed") {
		t.Errorf("error should contain original message, got: %v", err)
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > 0 && containsSubstr(s, substr))
}

func containsSubstr(s, sub string) bool {
	for i := 0; i <= len(s)-len(sub); i++ {
		if s[i:i+len(sub)] == sub {
			return true
		}
	}
	return false
}
