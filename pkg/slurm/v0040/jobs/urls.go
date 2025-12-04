package jobs

import "github.com/supergate-hub/slurm-client/pkg/slurmclient"

const (
	resourcePath  = "job"  // singular (single resource)
	resourcesPath = "jobs" // plural (list)
)

func listURL(c *slurmclient.ServiceClient) string {
	return c.ServiceURL(resourcesPath)
}

func getURL(c *slurmclient.ServiceClient, jobID string) string {
	return c.ServiceURL(resourcePath, jobID)
}

func submitURL(c *slurmclient.ServiceClient) string {
	return c.ServiceURL(resourcePath, "submit")
}

func updateURL(c *slurmclient.ServiceClient, jobID string) string {
	return c.ServiceURL(resourcePath, jobID)
}

func deleteURL(c *slurmclient.ServiceClient, jobID string) string {
	return c.ServiceURL(resourcePath, jobID)
}

