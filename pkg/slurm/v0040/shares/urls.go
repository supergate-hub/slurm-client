package shares

import "github.com/supergate-hub/slurm-client/pkg/slurmclient"

const resourcesPath = "shares"

func getURL(c *slurmclient.ServiceClient) string {
	return c.ServiceURL(resourcesPath)
}
