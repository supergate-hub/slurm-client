package diag

import "github.com/supergate-hub/slurm-client/pkg/slurmclient"

const resourcePath = "diag"

func getURL(c *slurmclient.ServiceClient) string {
	return c.ServiceURL(resourcePath)
}
