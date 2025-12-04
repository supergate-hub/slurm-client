// SPDX-FileCopyrightText: Copyright (C) SchedMD LLC.
// SPDX-License-Identifier: Apache-2.0

package jobs

import (
	api "github.com/supergate-hub/slurm-client/api/v0040"
)

// Job represents a Slurm job from the database.
type Job = api.V0040Job

// ListOpts specifies the options for listing jobs from slurmdb.
type ListOpts struct {
	// Users to filter by.
	Users []string
	// Accounts to filter by.
	Accounts []string
	// Clusters to filter by.
	Clusters []string
	// StartTime filters jobs started after this time.
	StartTime *string
	// EndTime filters jobs ended before this time.
	EndTime *string
	// Flags for the query.
	Flags []string
}

