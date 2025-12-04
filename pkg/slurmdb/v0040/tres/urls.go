// SPDX-FileCopyrightText: Copyright (C) SchedMD LLC.
// SPDX-License-Identifier: Apache-2.0

package tres

import "github.com/supergate-hub/slurm-client/pkg/slurmclient"

const resourcePath = "tres"

func listURL(c *slurmclient.ServiceClient) string {
	return c.ServiceURL(resourcePath)
}

func createURL(c *slurmclient.ServiceClient) string {
	return c.ServiceURL(resourcePath)
}



