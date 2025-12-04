// SPDX-FileCopyrightText: Copyright (C) SchedMD LLC.
// SPDX-License-Identifier: Apache-2.0

package v0040

import (
	"context"

	"github.com/supergate-hub/slurm-client/pkg/slurmclient"
	"github.com/supergate-hub/slurm-client/pkg/slurmdb/v0040/users"
)

// usersClient implements UsersInterface.
type usersClient struct {
	client *slurmclient.ServiceClient
}

var _ UsersInterface = (*usersClient)(nil)

// List retrieves all users.
func (c *usersClient) List(ctx context.Context, opts users.ListOpts) ([]User, error) {
	return users.List(ctx, c.client, opts)
}

// Get retrieves a single user by name.
func (c *usersClient) Get(ctx context.Context, userName string, opts ...users.GetOpts) (*User, error) {
	return users.Get(ctx, c.client, userName, opts...)
}

// Create creates one or more users.
func (c *usersClient) Create(ctx context.Context, opts users.CreateOpts) error {
	return users.Create(ctx, c.client, opts)
}

// Delete removes a user.
func (c *usersClient) Delete(ctx context.Context, userName string) error {
	return users.Delete(ctx, c.client, userName)
}



