// SPDX-FileCopyrightText: Copyright (C) SchedMD LLC.
// SPDX-License-Identifier: Apache-2.0

package instances

import "github.com/supergate-hub/slurm-client/pkg/slurmclient"

const (
	resourcePath  = "instance"  // singular
	resourcesPath = "instances" // plural
)

func listURL(c *slurmclient.ServiceClient) string {
	return c.ServiceURL(resourcesPath)
}

func getURL(c *slurmclient.ServiceClient) string {
	return c.ServiceURL(resourcePath)
}



