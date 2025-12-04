// SPDX-FileCopyrightText: Copyright (C) SchedMD LLC.
// SPDX-License-Identifier: Apache-2.0

package config

import (
	"context"

	api "github.com/supergate-hub/slurm-client/api/v0040"
	"github.com/supergate-hub/slurm-client/pkg/slurmclient"
)

// Get retrieves the Slurm database configuration.
func Get(ctx context.Context, client *slurmclient.ServiceClient) (*Config, error) {
	resp, err := client.Get(ctx, getURL(client))
	if err != nil {
		return nil, err
	}

	var result api.V0040OpenapiSlurmdbdConfigResp
	if err := slurmclient.DecodeResponse(resp, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// Update updates the Slurm database configuration.
func Update(ctx context.Context, client *slurmclient.ServiceClient, opts UpdateOpts) error {
	reqBody := api.V0040OpenapiSlurmdbdConfigResp{}

	if opts.Accounts != nil {
		accounts := api.V0040AccountList(opts.Accounts)
		reqBody.Accounts = &accounts
	}
	if opts.Associations != nil {
		assocs := api.V0040AssocList(opts.Associations)
		reqBody.Associations = &assocs
	}
	if opts.Clusters != nil {
		clusters := api.V0040ClusterRecList(opts.Clusters)
		reqBody.Clusters = &clusters
	}
	if opts.QOS != nil {
		qos := api.V0040QosList(opts.QOS)
		reqBody.Qos = &qos
	}
	if opts.TRES != nil {
		tres := api.V0040TresList(opts.TRES)
		reqBody.Tres = &tres
	}
	if opts.Users != nil {
		users := api.V0040UserList(opts.Users)
		reqBody.Users = &users
	}
	if opts.WCKeys != nil {
		wckeys := api.V0040WckeyList(opts.WCKeys)
		reqBody.Wckeys = &wckeys
	}

	resp, err := client.Post(ctx, updateURL(client), reqBody)
	if err != nil {
		return err
	}

	return slurmclient.CheckResponse(resp)
}
