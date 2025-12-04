// SPDX-FileCopyrightText: Copyright (C) SchedMD LLC.
// SPDX-License-Identifier: Apache-2.0

package accounts

import (
	"context"
	"fmt"
	"net/url"

	api "github.com/supergate-hub/slurm-client/api/v0040"
	"github.com/supergate-hub/slurm-client/pkg/slurmclient"
)

// List retrieves all accounts from the Slurm database.
func List(ctx context.Context, client *slurmclient.ServiceClient, opts ListOpts) ([]Account, error) {
	query := url.Values{}
	if opts.Description != nil {
		query.Set("description", *opts.Description)
	}
	if opts.WithDeleted != nil && *opts.WithDeleted {
		query.Set("with_deleted", "true")
	}
	if opts.WithAssociations != nil && *opts.WithAssociations {
		query.Set("with_assocs", "true")
	}
	if opts.WithCoordinators != nil && *opts.WithCoordinators {
		query.Set("with_coords", "true")
	}

	resp, err := client.GetWithQuery(ctx, listURL(client), query)
	if err != nil {
		return nil, err
	}

	var result api.V0040OpenapiAccountsResp
	if err := slurmclient.DecodeResponse(resp, &result); err != nil {
		return nil, err
	}

	if len(result.Accounts) == 0 {
		return []Account{}, nil
	}
	return result.Accounts, nil
}

// Get retrieves a single account by name.
func Get(ctx context.Context, client *slurmclient.ServiceClient, accountName string, opts ...GetOpts) (*Account, error) {
	query := url.Values{}
	if len(opts) > 0 {
		if opts[0].WithDeleted != nil && *opts[0].WithDeleted {
			query.Set("with_deleted", "true")
		}
		if opts[0].WithAssociations != nil && *opts[0].WithAssociations {
			query.Set("with_assocs", "true")
		}
		if opts[0].WithCoordinators != nil && *opts[0].WithCoordinators {
			query.Set("with_coords", "true")
		}
	}

	resp, err := client.GetWithQuery(ctx, getURL(client, accountName), query)
	if err != nil {
		return nil, err
	}

	var result api.V0040OpenapiAccountsResp
	if err := slurmclient.DecodeResponse(resp, &result); err != nil {
		return nil, err
	}

	if len(result.Accounts) == 0 {
		return nil, fmt.Errorf("account %s not found", accountName)
	}

	return &result.Accounts[0], nil
}

// Create creates one or more accounts in the Slurm database.
func Create(ctx context.Context, client *slurmclient.ServiceClient, opts CreateOpts) error {
	reqBody := struct {
		Accounts []Account `json:"accounts"`
	}{
		Accounts: opts.Accounts,
	}

	resp, err := client.Post(ctx, createURL(client), reqBody)
	if err != nil {
		return err
	}

	return slurmclient.CheckResponse(resp)
}

// Delete removes an account from the Slurm database.
func Delete(ctx context.Context, client *slurmclient.ServiceClient, accountName string) error {
	resp, err := client.Delete(ctx, deleteURL(client, accountName))
	if err != nil {
		return err
	}

	return slurmclient.CheckResponse(resp)
}
