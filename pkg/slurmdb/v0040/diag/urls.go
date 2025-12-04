// SPDX-FileCopyrightText: Copyright (C) SchedMD LLC.
// SPDX-License-Identifier: Apache-2.0

package diag

import "github.com/supergate-hub/slurm-client/pkg/slurmclient"

const resourcePath = "diag"

func getURL(c *slurmclient.ServiceClient) string {
	return c.ServiceURL(resourcePath)
}



