package clusters

import "github.com/supergate-hub/slurm-client/pkg/slurmclient"

const (
	resourcePath  = "cluster"  // singular
	resourcesPath = "clusters" // plural
)

func listURL(c *slurmclient.ServiceClient) string {
	return c.ServiceURL(resourcesPath)
}

func getURL(c *slurmclient.ServiceClient, clusterName string) string {
	return c.ServiceURL(resourcePath, clusterName)
}

func createURL(c *slurmclient.ServiceClient) string {
	return c.ServiceURL(resourcesPath)
}

func deleteURL(c *slurmclient.ServiceClient, clusterName string) string {
	return c.ServiceURL(resourcePath, clusterName)
}
