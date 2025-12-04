// SPDX-FileCopyrightText: Copyright (C) SchedMD LLC.
// SPDX-License-Identifier: Apache-2.0

package config

import "github.com/supergate-hub/slurm-client/pkg/slurmclient"

const resourcePath = "config"

func getURL(c *slurmclient.ServiceClient) string {
	return c.ServiceURL(resourcePath)
}

func updateURL(c *slurmclient.ServiceClient) string {
	return c.ServiceURL(resourcePath)
}



