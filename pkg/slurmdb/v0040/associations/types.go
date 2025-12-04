// SPDX-FileCopyrightText: Copyright (C) SchedMD LLC.
// SPDX-License-Identifier: Apache-2.0

package associations

import (
	api "github.com/supergate-hub/slurm-client/api/v0040"
)

// Association represents a Slurm association.
type Association = api.V0040Assoc

// ListOpts specifies the options for listing associations.
type ListOpts struct {
	// Account to filter by.
	Account *string
	// Cluster to filter by.
	Cluster *string
	// Partition to filter by.
	Partition *string
	// User to filter by.
	User *string
	// WithDeleted includes deleted associations.
	WithDeleted *bool
}

// GetOpts specifies the options for getting a single association.
type GetOpts struct {
	// Account name.
	Account *string
	// Cluster name.
	Cluster *string
	// Partition name.
	Partition *string
	// User name.
	User *string
}

// CreateOpts specifies the options for creating associations.
type CreateOpts struct {
	// Associations to create.
	Associations []Association
}

// DeleteOpts specifies the options for deleting associations.
type DeleteOpts struct {
	// Account name.
	Account *string
	// Cluster name.
	Cluster *string
	// Partition name.
	Partition *string
	// User name.
	User *string
}

