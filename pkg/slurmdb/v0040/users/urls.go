package users

import "github.com/supergate-hub/slurm-client/pkg/slurmclient"

const (
	resourcePath  = "user"  // singular
	resourcesPath = "users" // plural
)

func listURL(c *slurmclient.ServiceClient) string {
	return c.ServiceURL(resourcesPath)
}

func getURL(c *slurmclient.ServiceClient, userName string) string {
	return c.ServiceURL(resourcePath, userName)
}

func createURL(c *slurmclient.ServiceClient) string {
	return c.ServiceURL(resourcesPath)
}

func deleteURL(c *slurmclient.ServiceClient, userName string) string {
	return c.ServiceURL(resourcePath, userName)
}
