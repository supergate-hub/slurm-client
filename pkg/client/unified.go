package client

import (
	"context"
	"fmt"

	slurmv0040 "github.com/supergate-hub/slurm-client/pkg/slurm/v0040"
	slurmv0040jobs "github.com/supergate-hub/slurm-client/pkg/slurm/v0040/jobs"
	slurmv0040nodes "github.com/supergate-hub/slurm-client/pkg/slurm/v0040/nodes"
	"github.com/supergate-hub/slurm-client/pkg/slurmclient"
	slurmdbv0040 "github.com/supergate-hub/slurm-client/pkg/slurmdb/v0040"
	slurmdbv0040accounts "github.com/supergate-hub/slurm-client/pkg/slurmdb/v0040/accounts"
	// Add other versions here
)

// Version constants
const (
	V0_0_40 = "v0.0.40"
	V0_0_41 = "v0.0.41"
	V0_0_42 = "v0.0.42"
	V0_0_43 = "v0.0.43"
	V0_0_44 = "v0.0.44"
)

// UnifiedClient provides access to both Slurm and Slurmdb APIs for a specific version.
type UnifiedClient struct {
	Slurm   SlurmClientInterface
	Slurmdb SlurmdbClientInterface
	
	// Raw clients for when type assertion is needed
	rawSlurm   interface{}
	rawSlurmdb interface{}
}

// NewUnifiedClient creates a new client that wraps version-specific implementations.
func NewUnifiedClient(cfg Config) (*UnifiedClient, error) {
	if err := validate(&cfg); err != nil {
		return nil, err
	}

	// Convert pkg/client.Config to pkg/slurmclient.Config
	// Note: This assumes they are compatible or mapping is needed.
	// Since slurmclient.Config is in a different package, we map manually.
	scCfg := slurmclient.Config{
		Endpoint:  cfg.Server,
		AuthToken: cfg.AuthToken,
		HTTPClient: cfg.HTTPClient,
	}

	c := &UnifiedClient{}

	switch cfg.Version {
	case V0_0_40, "": // Default to v0040 if not specified for now, or handle error
		// Initialize Slurm Client
		sCli, err := slurmv0040.NewClient(scCfg)
		if err != nil {
			return nil, fmt.Errorf("failed to create slurm client: %w", err)
		}
		c.Slurm = &slurmv0040Wrapper{client: sCli}
		c.rawSlurm = sCli

		// Initialize Slurmdb Client
		dbCli, err := slurmdbv0040.NewClient(scCfg)
		if err != nil {
			return nil, fmt.Errorf("failed to create slurmdb client: %w", err)
		}
		c.Slurmdb = &slurmdbv0040Wrapper{client: dbCli}
		c.rawSlurmdb = dbCli

	// Add cases for V0_0_41, etc.

	default:
		return nil, fmt.Errorf("unsupported version: %s", cfg.Version)
	}

	return c, nil
}

// -----------------------------------------------------------------------------
// v0040 Wrappers (Adapters to satisfy SlurmClientInterface)
// -----------------------------------------------------------------------------

type slurmv0040Wrapper struct {
	client *slurmv0040.Client
}

func (w *slurmv0040Wrapper) Jobs() GenericJobService {
	return &slurmv0040JobsWrapper{service: w.client.Jobs}
}

func (w *slurmv0040Wrapper) Nodes() GenericNodeService {
	return &slurmv0040NodesWrapper{service: w.client.Nodes}
}

// Jobs Wrapper
type slurmv0040JobsWrapper struct {
	service slurmv0040.JobsInterface
}

func (w *slurmv0040JobsWrapper) List(ctx context.Context, opts interface{}) (interface{}, error) {
	var listOpts slurmv0040jobs.ListOpts
	if opts != nil {
		var ok bool
		listOpts, ok = opts.(slurmv0040jobs.ListOpts)
		if !ok {
			// Try pointer if passed as pointer
			if ptrOpts, ok := opts.(*slurmv0040jobs.ListOpts); ok {
				if ptrOpts != nil {
					listOpts = *ptrOpts
				}
			} else {
				// Optionally: support empty struct or map conversion here?
				// For now, enforce strict type or nil
				return nil, fmt.Errorf("invalid options type for v0040: expected %T, got %T", listOpts, opts)
			}
		}
	}
	return w.service.List(ctx, listOpts)
}

func (w *slurmv0040JobsWrapper) Get(ctx context.Context, id string) (interface{}, error) {
	// GetOpts handling (variadic in interface)
	// Since our Generic interface signature for Get is `Get(ctx, id)`, we pass no options.
	// If Generic interface needs opts, we should update it.
	return w.service.Get(ctx, id)
}

func (w *slurmv0040JobsWrapper) Delete(ctx context.Context, id string) error {
	return w.service.Delete(ctx, id)
}

// Nodes Wrapper
type slurmv0040NodesWrapper struct {
	service slurmv0040.NodesInterface
}

func (w *slurmv0040NodesWrapper) List(ctx context.Context, opts interface{}) (interface{}, error) {
	var listOpts slurmv0040nodes.ListOpts
	if opts != nil {
		if val, ok := opts.(slurmv0040nodes.ListOpts); ok {
			listOpts = val
		} else if ptr, ok := opts.(*slurmv0040nodes.ListOpts); ok && ptr != nil {
			listOpts = *ptr
		} else {
			return nil, fmt.Errorf("invalid options type for v0040 nodes: expected %T", listOpts)
		}
	}
	return w.service.List(ctx, listOpts)
}

func (w *slurmv0040NodesWrapper) Get(ctx context.Context, name string) (interface{}, error) {
	return w.service.Get(ctx, name)
}

// -----------------------------------------------------------------------------
// v0040 Slurmdb Wrappers
// -----------------------------------------------------------------------------

type slurmdbv0040Wrapper struct {
	client *slurmdbv0040.Client
}

func (w *slurmdbv0040Wrapper) Accounts() GenericAccountService {
	return &slurmdbv0040AccountsWrapper{service: w.client.Accounts}
}

type slurmdbv0040AccountsWrapper struct {
	service slurmdbv0040.AccountsInterface
}

func (w *slurmdbv0040AccountsWrapper) List(ctx context.Context, opts interface{}) (interface{}, error) {
	var listOpts slurmdbv0040accounts.ListOpts
	if opts != nil {
		if val, ok := opts.(slurmdbv0040accounts.ListOpts); ok {
			listOpts = val
		} else if ptr, ok := opts.(*slurmdbv0040accounts.ListOpts); ok && ptr != nil {
			listOpts = *ptr
		} else {
			return nil, fmt.Errorf("invalid options type for v0040 accounts: expected %T", listOpts)
		}
	}
	return w.service.List(ctx, listOpts)
}

func (w *slurmdbv0040AccountsWrapper) Get(ctx context.Context, name string) (interface{}, error) {
	return w.service.Get(ctx, name)
}
