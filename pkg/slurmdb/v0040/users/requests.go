// SPDX-FileCopyrightText: Copyright (C) SchedMD LLC.
// SPDX-License-Identifier: Apache-2.0

package users

import (
	"context"
	"fmt"
	"net/url"

	api "github.com/supergate-hub/slurm-client/api/v0040"
	"github.com/supergate-hub/slurm-client/pkg/slurmclient"
)

// List retrieves all users from the Slurm database.
func List(ctx context.Context, client *slurmclient.ServiceClient, opts ListOpts) ([]User, error) {
	query := url.Values{}
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

	var result api.V0040OpenapiUsersResp
	if err := slurmclient.DecodeResponse(resp, &result); err != nil {
		return nil, err
	}

	if len(result.Users) == 0 {
		return []User{}, nil
	}
	return result.Users, nil
}

// Get retrieves a single user by name.
func Get(ctx context.Context, client *slurmclient.ServiceClient, userName string, opts ...GetOpts) (*User, error) {
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

	resp, err := client.GetWithQuery(ctx, getURL(client, userName), query)
	if err != nil {
		return nil, err
	}

	var result api.V0040OpenapiUsersResp
	if err := slurmclient.DecodeResponse(resp, &result); err != nil {
		return nil, err
	}

	if len(result.Users) == 0 {
		return nil, fmt.Errorf("user %s not found", userName)
	}

	return &result.Users[0], nil
}

// Create creates one or more users in the Slurm database.
func Create(ctx context.Context, client *slurmclient.ServiceClient, opts CreateOpts) error {
	reqBody := struct {
		Users []User `json:"users"`
	}{
		Users: opts.Users,
	}

	resp, err := client.Post(ctx, createURL(client), reqBody)
	if err != nil {
		return err
	}

	return slurmclient.CheckResponse(resp)
}

// Delete removes a user from the Slurm database.
func Delete(ctx context.Context, client *slurmclient.ServiceClient, userName string) error {
	resp, err := client.Delete(ctx, deleteURL(client, userName))
	if err != nil {
		return err
	}

	return slurmclient.CheckResponse(resp)
}
