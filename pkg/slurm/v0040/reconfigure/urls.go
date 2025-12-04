package reconfigure

import "github.com/supergate-hub/slurm-client/pkg/slurmclient"

const resourcePath = "reconfigure"

func reconfigureURL(c *slurmclient.ServiceClient) string {
	return c.ServiceURL(resourcePath)
}
