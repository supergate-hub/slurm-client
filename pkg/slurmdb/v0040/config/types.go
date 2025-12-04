// SPDX-FileCopyrightText: Copyright (C) SchedMD LLC.
// SPDX-License-Identifier: Apache-2.0

package config

import (
	api "github.com/supergate-hub/slurm-client/api/v0040"
)

// Config represents the Slurm database configuration response.
type Config = api.V0040OpenapiSlurmdbdConfigResp

// UpdateOpts specifies the options for updating the configuration.
type UpdateOpts struct {
	// Accounts to update.
	Accounts []api.V0040Account
	// Associations to update.
	Associations []api.V0040Assoc
	// Clusters to update.
	Clusters []api.V0040ClusterRec
	// QOS to update.
	QOS []api.V0040Qos
	// TRES to update.
	TRES []api.V0040Tres
	// Users to update.
	Users []api.V0040User
	// WCKeys to update.
	WCKeys []api.V0040Wckey
}
