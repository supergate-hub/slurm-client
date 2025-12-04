package diag

import (
	"context"

	api "github.com/supergate-hub/slurm-client/api/v0040"
	"github.com/supergate-hub/slurm-client/pkg/slurmclient"
)

// Get retrieves diagnostic information from the Slurm cluster.
func Get(ctx context.Context, client *slurmclient.ServiceClient) (*DiagInfo, error) {
	resp, err := client.Get(ctx, getURL(client))
	if err != nil {
		return nil, err
	}

	var result api.V0040OpenapiDiagResp
	if err := slurmclient.DecodeResponse(resp, &result); err != nil {
		return nil, err
	}

	return &result.Statistics, nil
}
