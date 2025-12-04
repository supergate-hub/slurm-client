// SPDX-FileCopyrightText: Copyright (C) SchedMD LLC.
// SPDX-License-Identifier: Apache-2.0

package tres

import (
	api "github.com/supergate-hub/slurm-client/api/v0040"
)

// TRES represents a Slurm Trackable Resource.
type TRES = api.V0040Tres

// ListOpts specifies the options for listing TRES.
type ListOpts struct {
	// Currently no options are supported.
}

// CreateOpts specifies the options for creating TRES.
type CreateOpts struct {
	// TRES list to create.
	TRES []TRES
}

