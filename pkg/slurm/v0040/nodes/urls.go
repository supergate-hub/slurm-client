package nodes

import "github.com/supergate-hub/slurm-client/pkg/slurmclient"

const (
	resourcePath  = "node"  // singular
	resourcesPath = "nodes" // plural
)

func listURL(c *slurmclient.ServiceClient) string {
	return c.ServiceURL(resourcesPath)
}

func getURL(c *slurmclient.ServiceClient, nodeName string) string {
	return c.ServiceURL(resourcePath, nodeName)
}

func updateURL(c *slurmclient.ServiceClient, nodeName string) string {
	return c.ServiceURL(resourcePath, nodeName)
}

func deleteURL(c *slurmclient.ServiceClient, nodeName string) string {
	return c.ServiceURL(resourcePath, nodeName)
}
