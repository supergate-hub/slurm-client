package reconfigure

import (
	"context"

	"github.com/supergate-hub/slurm-client/pkg/slurmclient"
)

// Trigger initiates a Slurm reconfiguration.
// This causes slurmctld to re-read its configuration file.
func Trigger(ctx context.Context, client *slurmclient.ServiceClient) error {
	resp, err := client.Get(ctx, reconfigureURL(client))
	if err != nil {
		return err
	}

	return slurmclient.CheckResponse(resp)
}

