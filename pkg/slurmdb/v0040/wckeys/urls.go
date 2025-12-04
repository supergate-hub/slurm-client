// SPDX-FileCopyrightText: Copyright (C) SchedMD LLC.
// SPDX-License-Identifier: Apache-2.0

package wckeys

import "github.com/supergate-hub/slurm-client/pkg/slurmclient"

const (
	resourcePath  = "wckey"  // singular
	resourcesPath = "wckeys" // plural
)

func listURL(c *slurmclient.ServiceClient) string {
	return c.ServiceURL(resourcesPath)
}

func getURL(c *slurmclient.ServiceClient, wckeyID string) string {
	return c.ServiceURL(resourcePath, wckeyID)
}

func createURL(c *slurmclient.ServiceClient) string {
	return c.ServiceURL(resourcesPath)
}

func deleteURL(c *slurmclient.ServiceClient, wckeyID string) string {
	return c.ServiceURL(resourcePath, wckeyID)
}




