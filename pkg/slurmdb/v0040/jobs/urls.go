// SPDX-FileCopyrightText: Copyright (C) SchedMD LLC.
// SPDX-License-Identifier: Apache-2.0

package jobs

import "github.com/supergate-hub/slurm-client/pkg/slurmclient"

const (
	resourcePath  = "job"  // singular
	resourcesPath = "jobs" // plural
)

func listURL(c *slurmclient.ServiceClient) string {
	return c.ServiceURL(resourcesPath)
}

func getURL(c *slurmclient.ServiceClient, jobID string) string {
	return c.ServiceURL(resourcePath, jobID)
}




