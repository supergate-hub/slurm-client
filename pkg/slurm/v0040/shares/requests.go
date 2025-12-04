package shares

import (
	"context"
	"net/url"

	api "github.com/supergate-hub/slurm-client/api/v0040"
	"github.com/supergate-hub/slurm-client/pkg/slurmclient"
)

// Get retrieves fair share information from the Slurm cluster.
func Get(ctx context.Context, client *slurmclient.ServiceClient, opts GetOpts) ([]ShareInfo, error) {
	query := url.Values{}
	for _, account := range opts.Accounts {
		query.Add("accounts", account)
	}
	for _, user := range opts.Users {
		query.Add("users", user)
	}

	resp, err := client.GetWithQuery(ctx, getURL(client), query)
	if err != nil {
		return nil, err
	}

	var result api.V0040OpenapiSharesResp
	if err := slurmclient.DecodeResponse(resp, &result); err != nil {
		return nil, err
	}

	if result.Shares.Shares == nil {
		return []ShareInfo{}, nil
	}
	return *result.Shares.Shares, nil
}
