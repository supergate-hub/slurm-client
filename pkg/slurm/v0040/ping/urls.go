package ping

import "github.com/supergate-hub/slurm-client/pkg/slurmclient"

const resourcePath = "ping"

func pingURL(c *slurmclient.ServiceClient) string {
	return c.ServiceURL(resourcePath)
}
