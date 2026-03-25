package slurmclient

import (
	"context"
	"fmt"
	"sync"

	v0040 "github.com/supergate-hub/slurm-client/api/v0040"
)

// ManagerOpts configures a multi-cluster Manager.
type ManagerOpts struct {
	// Clusters maps cluster names to their connection options.
	// At least one cluster is required.
	Clusters map[string]AuthOpts
}

// Manager manages multiple Slurm clusters. It wraps individual Clients
// and provides both single-cluster access and cross-cluster operations.
//
// Usage:
//
//	manager, err := slurmclient.NewManager(ctx, slurmclient.ManagerOpts{
//	    Clusters: map[string]slurmclient.AuthOpts{
//	        "gpu": {Endpoint: "http://gpu:6820", AuthToken: "t1"},
//	        "cpu": {Endpoint: "http://cpu:6820", AuthToken: "t2"},
//	    },
//	})
//	jobs, _ := manager.Cluster("gpu").Slurm.Jobs().List(ctx)
//	allJobs, _ := manager.AllJobs(ctx)
type Manager struct {
	clients map[string]*Client
	names   []string // ordered cluster names for deterministic iteration
}

// ClusterResult holds the result of an operation on a single cluster
// within a multi-cluster query.
type ClusterResult[T any] struct {
	// Cluster is the name of the cluster this result came from.
	Cluster string

	// Data holds the successful result. Zero value if Err is non-nil.
	Data T

	// Err is non-nil if this cluster's operation failed.
	// Other clusters' results are still returned (partial result pattern).
	Err error
}

// NewManager creates a Manager that connects to multiple Slurm clusters.
// Clients are initialized in parallel using errgroup for minimal latency.
func NewManager(ctx context.Context, opts ManagerOpts) (*Manager, error) {
	if len(opts.Clusters) == 0 {
		return nil, fmt.Errorf("%w: at least one cluster is required", ErrValidation)
	}

	type result struct {
		name   string
		client *Client
		err    error
	}

	results := make(chan result, len(opts.Clusters))
	var wg sync.WaitGroup

	for name, authOpts := range opts.Clusters {
		wg.Add(1)
		go func(n string, ao AuthOpts) {
			defer wg.Done()
			c, err := NewClient(ctx, ao)
			results <- result{name: n, client: c, err: err}
		}(name, authOpts)
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	clients := make(map[string]*Client, len(opts.Clusters))
	var names []string
	var initErrors []error

	for r := range results {
		if r.err != nil {
			initErrors = append(initErrors, fmt.Errorf("cluster %q: %w", r.name, r.err))
			continue
		}
		clients[r.name] = r.client
		names = append(names, r.name)
	}

	if len(clients) == 0 {
		return nil, &MultiClusterError{Errors: initErrors}
	}

	// Log partial failures but continue with available clusters
	m := &Manager{
		clients: clients,
		names:   names,
	}
	return m, nil
}

// Cluster returns the Client for a named cluster.
func (m *Manager) Cluster(name string) (*Client, error) {
	c, ok := m.clients[name]
	if !ok {
		return nil, fmt.Errorf("%w: cluster %q not found (available: %v)", ErrNotFound, name, m.names)
	}
	return c, nil
}

// ClusterNames returns the names of all connected clusters.
func (m *Manager) ClusterNames() []string {
	out := make([]string, len(m.names))
	copy(out, m.names)
	return out
}

// AllJobs queries jobs from all clusters in parallel.
func (m *Manager) AllJobs(ctx context.Context) []ClusterResult[[]v0040.V0040JobInfo] {
	return forEachCluster(ctx, m, func(c *Client) ([]v0040.V0040JobInfo, error) {
		return c.Slurm.Jobs().List(ctx)
	})
}

// AllNodes queries nodes from all clusters in parallel.
func (m *Manager) AllNodes(ctx context.Context) []ClusterResult[[]v0040.V0040Node] {
	return forEachCluster(ctx, m, func(c *Client) ([]v0040.V0040Node, error) {
		return c.Slurm.Nodes().List(ctx)
	})
}

// AllPartitions queries partitions from all clusters in parallel.
func (m *Manager) AllPartitions(ctx context.Context) []ClusterResult[[]v0040.V0040PartitionInfo] {
	return forEachCluster(ctx, m, func(c *Client) ([]v0040.V0040PartitionInfo, error) {
		return c.Slurm.Partitions().List(ctx)
	})
}

// ForEach runs fn on every cluster in parallel and collects results.
// Use the typed All* methods (AllJobs, AllNodes, etc.) when possible.
func (m *Manager) ForEach(ctx context.Context, fn func(name string, c *Client) (any, error)) []ClusterResult[any] {
	return forEachCluster(ctx, m, func(c *Client) (any, error) {
		name := m.clientName(c)
		return fn(name, c)
	})
}

func (m *Manager) clientName(c *Client) string {
	for n, cl := range m.clients {
		if cl == c {
			return n
		}
	}
	return ""
}

// forEachCluster is the generic parallel-dispatch helper.
func forEachCluster[T any](ctx context.Context, m *Manager, fn func(c *Client) (T, error)) []ClusterResult[T] {
	type indexed struct {
		idx    int
		result ClusterResult[T]
	}

	ch := make(chan indexed, len(m.names))
	var wg sync.WaitGroup

	for i, name := range m.names {
		wg.Add(1)
		go func(idx int, n string) {
			defer wg.Done()
			data, err := fn(m.clients[n])
			ch <- indexed{idx: idx, result: ClusterResult[T]{
				Cluster: n,
				Data:    data,
				Err:     err,
			}}
		}(i, name)
	}

	go func() {
		wg.Wait()
		close(ch)
	}()

	results := make([]ClusterResult[T], len(m.names))
	for r := range ch {
		results[r.idx] = r.result
	}
	return results
}

// MultiClusterError is returned when all clusters fail during an operation.
type MultiClusterError struct {
	Errors []error
}

func (e *MultiClusterError) Error() string {
	msgs := make([]string, len(e.Errors))
	for i, err := range e.Errors {
		msgs[i] = err.Error()
	}
	return fmt.Sprintf("slurm: all clusters failed: %v", msgs)
}
