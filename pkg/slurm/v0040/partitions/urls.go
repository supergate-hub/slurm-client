package partitions

import "github.com/supergate-hub/slurm-client/pkg/slurmclient"

const (
	resourcePath  = "partition"  // singular
	resourcesPath = "partitions" // plural
)

func listURL(c *slurmclient.ServiceClient) string {
	return c.ServiceURL(resourcesPath)
}

func getURL(c *slurmclient.ServiceClient, partitionName string) string {
	return c.ServiceURL(resourcePath, partitionName)
}
