package associations

import "github.com/supergate-hub/slurm-client/pkg/slurmclient"

const (
	resourcePath  = "association"  // singular
	resourcesPath = "associations" // plural
)

func listURL(c *slurmclient.ServiceClient) string {
	return c.ServiceURL(resourcesPath)
}

func getURL(c *slurmclient.ServiceClient) string {
	return c.ServiceURL(resourcePath)
}

func createURL(c *slurmclient.ServiceClient) string {
	return c.ServiceURL(resourcesPath)
}

func deleteURL(c *slurmclient.ServiceClient) string {
	return c.ServiceURL(resourcePath)
}
