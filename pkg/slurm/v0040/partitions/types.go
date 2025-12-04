package partitions

import (
	api "github.com/supergate-hub/slurm-client/api/v0040"
)

// Partition represents a Slurm partition.
type Partition = api.V0040PartitionInfo

// ListOpts specifies the options for listing partitions.
type ListOpts struct {
	// UpdateTime filters partitions updated after this time.
	UpdateTime *int64
}

// GetOpts specifies the options for getting a single partition.
type GetOpts struct {
	// UpdateTime filters if the partition was updated after this time.
	UpdateTime *int64
}

