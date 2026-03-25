package slurmclient_test

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	slurmclient "github.com/supergate-hub/slurm-client"
)

func TestParseConfig(t *testing.T) {
	t.Run("valid config", func(t *testing.T) {
		dir := t.TempDir()
		path := filepath.Join(dir, "clusters.yaml")
		os.WriteFile(path, []byte(`
clusters:
  gpu:
    endpoint: http://gpu:6820
    token: token1
    version: v0.0.40
  cpu:
    endpoint: http://cpu:6820
    token: token2
`), 0644)

		opts, err := slurmclient.ParseConfig(path)
		if err != nil {
			t.Fatalf("ParseConfig: %v", err)
		}
		if len(opts.Clusters) != 2 {
			t.Fatalf("expected 2 clusters, got %d", len(opts.Clusters))
		}
		gpu := opts.Clusters["gpu"]
		if gpu.Endpoint != "http://gpu:6820" {
			t.Fatalf("expected http://gpu:6820, got %s", gpu.Endpoint)
		}
		if gpu.Version != "v0.0.40" {
			t.Fatalf("expected v0.0.40, got %s", gpu.Version)
		}
	})

	t.Run("env var expansion", func(t *testing.T) {
		t.Setenv("TEST_SLURM_TOKEN", "secret123")

		dir := t.TempDir()
		path := filepath.Join(dir, "clusters.yaml")
		os.WriteFile(path, []byte(`
clusters:
  test:
    endpoint: http://test:6820
    token: ${TEST_SLURM_TOKEN}
`), 0644)

		opts, err := slurmclient.ParseConfig(path)
		if err != nil {
			t.Fatalf("ParseConfig: %v", err)
		}
		if opts.Clusters["test"].AuthToken != "secret123" {
			t.Fatalf("expected expanded token, got %s", opts.Clusters["test"].AuthToken)
		}
	})

	t.Run("file not found", func(t *testing.T) {
		_, err := slurmclient.ParseConfig("/nonexistent/path.yaml")
		if err == nil {
			t.Fatal("expected error for missing file")
		}
	})

	t.Run("invalid yaml", func(t *testing.T) {
		dir := t.TempDir()
		path := filepath.Join(dir, "bad.yaml")
		os.WriteFile(path, []byte(`{{{invalid`), 0644)

		_, err := slurmclient.ParseConfig(path)
		if err == nil {
			t.Fatal("expected error for invalid YAML")
		}
	})

	t.Run("empty clusters", func(t *testing.T) {
		dir := t.TempDir()
		path := filepath.Join(dir, "empty.yaml")
		os.WriteFile(path, []byte(`clusters: {}`), 0644)

		_, err := slurmclient.ParseConfig(path)
		if !errors.Is(err, slurmclient.ErrValidation) {
			t.Fatalf("expected ErrValidation, got %v", err)
		}
	})

	t.Run("missing endpoint", func(t *testing.T) {
		dir := t.TempDir()
		path := filepath.Join(dir, "noep.yaml")
		os.WriteFile(path, []byte(`
clusters:
  bad:
    token: tok
`), 0644)

		_, err := slurmclient.ParseConfig(path)
		if !errors.Is(err, slurmclient.ErrValidation) {
			t.Fatalf("expected ErrValidation, got %v", err)
		}
	})

	t.Run("missing token", func(t *testing.T) {
		dir := t.TempDir()
		path := filepath.Join(dir, "notok.yaml")
		os.WriteFile(path, []byte(`
clusters:
  bad:
    endpoint: http://test:6820
`), 0644)

		_, err := slurmclient.ParseConfig(path)
		if !errors.Is(err, slurmclient.ErrValidation) {
			t.Fatalf("expected ErrValidation, got %v", err)
		}
	})

	t.Run("unix socket", func(t *testing.T) {
		dir := t.TempDir()
		path := filepath.Join(dir, "unix.yaml")
		os.WriteFile(path, []byte(`
clusters:
  local:
    unix_socket: /var/run/slurmrestd.sock
    token: tok
`), 0644)

		opts, err := slurmclient.ParseConfig(path)
		if err != nil {
			t.Fatalf("ParseConfig: %v", err)
		}
		if opts.Clusters["local"].UnixSocket != "/var/run/slurmrestd.sock" {
			t.Fatalf("expected unix socket, got %s", opts.Clusters["local"].UnixSocket)
		}
	})
}
