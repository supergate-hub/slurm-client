package licenses

import "github.com/supergate-hub/slurm-client/pkg/slurmclient"

const resourcesPath = "licenses"

func listURL(c *slurmclient.ServiceClient) string {
	return c.ServiceURL(resourcesPath)
}
