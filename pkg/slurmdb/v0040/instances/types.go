// SPDX-FileCopyrightText: Copyright (C) SchedMD LLC.
// SPDX-License-Identifier: Apache-2.0

package instances

import (
	api "github.com/supergate-hub/slurm-client/api/v0040"
)

// Instance represents a Slurm instance.
type Instance = api.V0040Instance

// ListOpts specifies the options for listing instances.
type ListOpts struct {
	// Cluster to filter by.
	Cluster *string
	// Extra info string to filter by.
	Extra *string
	// Format of the response.
	Format *string
	// InstanceID to filter by.
	InstanceID *string
	// InstanceType to filter by.
	InstanceType *string
	// NodeList to filter by.
	NodeList *string
}

// GetOpts specifies the options for getting a single instance.
type GetOpts struct {
	// Cluster name.
	Cluster *string
	// InstanceID.
	InstanceID *string
	// InstanceType.
	InstanceType *string
	// NodeList.
	NodeList *string
}

