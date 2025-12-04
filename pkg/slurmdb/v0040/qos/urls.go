// SPDX-FileCopyrightText: Copyright (C) SchedMD LLC.
// SPDX-License-Identifier: Apache-2.0

package qos

import "github.com/supergate-hub/slurm-client/pkg/slurmclient"

const resourcePath = "qos"

func listURL(c *slurmclient.ServiceClient) string {
	return c.ServiceURL(resourcePath)
}

func getURL(c *slurmclient.ServiceClient, qosName string) string {
	return c.ServiceURL(resourcePath, qosName)
}

func createURL(c *slurmclient.ServiceClient) string {
	return c.ServiceURL(resourcePath)
}

func deleteURL(c *slurmclient.ServiceClient, qosName string) string {
	return c.ServiceURL(resourcePath, qosName)
}




