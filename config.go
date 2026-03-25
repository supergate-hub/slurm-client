package slurmclient

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v2"
)

// ClusterConfig represents a single cluster entry in the configuration file.
type ClusterConfig struct {
	Endpoint   string `yaml:"endpoint"`
	Token      string `yaml:"token"`
	Version    string `yaml:"version,omitempty"`
	UnixSocket string `yaml:"unix_socket,omitempty"`
}

// Config represents the multi-cluster configuration file format.
//
// Example clusters.yaml:
//
//	clusters:
//	  gpu:
//	    endpoint: http://gpu-cluster:6820
//	    token: ${GPU_SLURM_TOKEN}
//	  cpu:
//	    endpoint: http://cpu-cluster:6820
//	    token: ${CPU_SLURM_TOKEN}
type Config struct {
	Clusters map[string]ClusterConfig `yaml:"clusters"`
}

// ParseConfig reads a YAML configuration file and returns ManagerOpts.
// Environment variables in the format ${VAR} are expanded.
func ParseConfig(path string) (*ManagerOpts, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read config %s: %w", path, err)
	}

	// Expand environment variables before parsing
	expanded := os.ExpandEnv(string(data))

	var cfg Config
	if err := yaml.Unmarshal([]byte(expanded), &cfg); err != nil {
		return nil, fmt.Errorf("parse config %s: %w", path, err)
	}

	if len(cfg.Clusters) == 0 {
		return nil, fmt.Errorf("%w: no clusters defined in %s", ErrValidation, path)
	}

	opts := &ManagerOpts{
		Clusters: make(map[string]AuthOpts, len(cfg.Clusters)),
	}

	for name, cc := range cfg.Clusters {
		if cc.Endpoint == "" && cc.UnixSocket == "" {
			return nil, fmt.Errorf("%w: cluster %q: endpoint or unix_socket is required", ErrValidation, name)
		}
		if cc.Token == "" {
			return nil, fmt.Errorf("%w: cluster %q: token is required", ErrValidation, name)
		}

		ao := AuthOpts{
			Endpoint:   cc.Endpoint,
			AuthToken:  cc.Token,
			UnixSocket: cc.UnixSocket,
		}
		if cc.Version != "" {
			ao.Version = Version(cc.Version)
		}

		opts.Clusters[name] = ao
	}

	return opts, nil
}
