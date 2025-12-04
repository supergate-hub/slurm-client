package licenses

import (
	"context"

	api "github.com/supergate-hub/slurm-client/api/v0040"
	"github.com/supergate-hub/slurm-client/pkg/slurmclient"
)

// List retrieves all licenses from the Slurm cluster.
func List(ctx context.Context, client *slurmclient.ServiceClient, opts ListOpts) ([]License, error) {
	resp, err := client.Get(ctx, listURL(client))
	if err != nil {
		return nil, err
	}

	var result api.V0040OpenapiLicensesResp
	if err := slurmclient.DecodeResponse(resp, &result); err != nil {
		return nil, err
	}

	if len(result.Licenses) == 0 {
		return []License{}, nil
	}
	return result.Licenses, nil
}
