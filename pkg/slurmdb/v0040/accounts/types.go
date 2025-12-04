// SPDX-FileCopyrightText: Copyright (C) SchedMD LLC.
// SPDX-License-Identifier: Apache-2.0

package accounts

import (
	api "github.com/supergate-hub/slurm-client/api/v0040"
)

// Account represents a Slurm account.
type Account = api.V0040Account

// ListOpts specifies the options for listing accounts.
type ListOpts struct {
	// Description to filter by.
	Description *string
	// WithDeleted includes deleted accounts in the response.
	WithDeleted *bool
	// WithAssociations includes association information.
	WithAssociations *bool
	// WithCoordinators includes coordinator information.
	WithCoordinators *bool
}

// GetOpts specifies the options for getting a single account.
type GetOpts struct {
	// WithDeleted includes deleted accounts.
	WithDeleted *bool
	// WithAssociations includes association information.
	WithAssociations *bool
	// WithCoordinators includes coordinator information.
	WithCoordinators *bool
}

// CreateOpts specifies the options for creating accounts.
type CreateOpts struct {
	// Accounts to create.
	Accounts []Account
}

