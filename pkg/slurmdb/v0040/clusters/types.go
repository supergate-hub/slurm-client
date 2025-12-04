// SPDX-FileCopyrightText: Copyright (C) SchedMD LLC.
// SPDX-License-Identifier: Apache-2.0

package clusters

import (
	api "github.com/supergate-hub/slurm-client/api/v0040"
)

// Cluster represents a Slurm cluster.
type Cluster = api.V0040ClusterRec

// ListOpts specifies the options for listing clusters.
type ListOpts struct {
	// WithDeleted includes deleted clusters.
	WithDeleted *bool
}

// GetOpts specifies the options for getting a single cluster.
type GetOpts struct {
	// WithDeleted includes deleted clusters.
	WithDeleted *bool
}

// CreateOpts specifies the options for creating clusters.
type CreateOpts struct {
	// Clusters to create.
	Clusters []Cluster
}

// DeleteOpts specifies the options for deleting clusters.
type DeleteOpts struct {
	// Classification to filter by.
	Classification *string
}

