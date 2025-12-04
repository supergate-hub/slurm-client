package ping

import (
	"context"

	api "github.com/supergate-hub/slurm-client/api/v0040"
	"github.com/supergate-hub/slurm-client/pkg/slurmclient"
)

// Ping checks the status of Slurm controllers.
func Ping(ctx context.Context, client *slurmclient.ServiceClient) ([]PingInfo, error) {
	resp, err := client.Get(ctx, pingURL(client))
	if err != nil {
		return nil, err
	}

	var result api.V0040OpenapiPingArrayResp
	if err := slurmclient.DecodeResponse(resp, &result); err != nil {
		return nil, err
	}

	if len(result.Pings) == 0 {
		return []PingInfo{}, nil
	}
	return result.Pings, nil
}
