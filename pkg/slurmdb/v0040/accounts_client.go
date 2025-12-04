// SPDX-FileCopyrightText: Copyright (C) SchedMD LLC.
// SPDX-License-Identifier: Apache-2.0

package v0040

import (
	"context"

	"github.com/supergate-hub/slurm-client/pkg/slurmclient"
	"github.com/supergate-hub/slurm-client/pkg/slurmdb/v0040/accounts"
)

// accountsClient implements AccountsInterface.
type accountsClient struct {
	client *slurmclient.ServiceClient
}

var _ AccountsInterface = (*accountsClient)(nil)

// List retrieves all accounts.
func (c *accountsClient) List(ctx context.Context, opts accounts.ListOpts) ([]Account, error) {
	return accounts.List(ctx, c.client, opts)
}

// Get retrieves a single account by name.
func (c *accountsClient) Get(ctx context.Context, accountName string, opts ...accounts.GetOpts) (*Account, error) {
	return accounts.Get(ctx, c.client, accountName, opts...)
}

// Create creates one or more accounts.
func (c *accountsClient) Create(ctx context.Context, opts accounts.CreateOpts) error {
	return accounts.Create(ctx, c.client, opts)
}

// Delete removes an account.
func (c *accountsClient) Delete(ctx context.Context, accountName string) error {
	return accounts.Delete(ctx, c.client, accountName)
}



