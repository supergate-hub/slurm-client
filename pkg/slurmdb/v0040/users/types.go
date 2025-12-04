// SPDX-FileCopyrightText: Copyright (C) SchedMD LLC.
// SPDX-License-Identifier: Apache-2.0

package users

import (
	api "github.com/supergate-hub/slurm-client/api/v0040"
)

// User represents a Slurm user.
type User = api.V0040User

// ListOpts specifies the options for listing users.
type ListOpts struct {
	// WithDeleted includes deleted users in the response.
	WithDeleted *bool
	// WithAssociations includes association information.
	WithAssociations *bool
	// WithCoordinators includes coordinator information.
	WithCoordinators *bool
}

// GetOpts specifies the options for getting a single user.
type GetOpts struct {
	// WithDeleted includes deleted users.
	WithDeleted *bool
	// WithAssociations includes association information.
	WithAssociations *bool
	// WithCoordinators includes coordinator information.
	WithCoordinators *bool
}

// CreateOpts specifies the options for creating users.
type CreateOpts struct {
	// Users to create.
	Users []User
}

