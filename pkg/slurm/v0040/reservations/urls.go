package reservations

import "github.com/supergate-hub/slurm-client/pkg/slurmclient"

const (
	resourcePath  = "reservation"  // singular
	resourcesPath = "reservations" // plural
)

func listURL(c *slurmclient.ServiceClient) string {
	return c.ServiceURL(resourcesPath)
}

func getURL(c *slurmclient.ServiceClient, reservationName string) string {
	return c.ServiceURL(resourcePath, reservationName)
}
