package accounts

import "github.com/supergate-hub/slurm-client/pkg/slurmclient"

const (
	resourcePath  = "account"  // singular
	resourcesPath = "accounts" // plural
)

func listURL(c *slurmclient.ServiceClient) string {
	return c.ServiceURL(resourcesPath)
}

func getURL(c *slurmclient.ServiceClient, accountName string) string {
	return c.ServiceURL(resourcePath, accountName)
}

func createURL(c *slurmclient.ServiceClient) string {
	return c.ServiceURL(resourcesPath)
}

func deleteURL(c *slurmclient.ServiceClient, accountName string) string {
	return c.ServiceURL(resourcePath, accountName)
}
