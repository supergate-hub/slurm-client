package reservations

import (
	"context"
	"fmt"
	"net/url"

	api "github.com/supergate-hub/slurm-client/api/v0040"
	"github.com/supergate-hub/slurm-client/pkg/slurmclient"
)

// List retrieves all reservations from the Slurm cluster.
func List(ctx context.Context, client *slurmclient.ServiceClient, opts ListOpts) ([]Reservation, error) {
	query := url.Values{}
	if opts.UpdateTime != nil {
		query.Set("update_time", fmt.Sprintf("%d", *opts.UpdateTime))
	}

	resp, err := client.GetWithQuery(ctx, listURL(client), query)
	if err != nil {
		return nil, err
	}

	var result api.V0040OpenapiReservationResp
	if err := slurmclient.DecodeResponse(resp, &result); err != nil {
		return nil, err
	}

	if len(result.Reservations) == 0 {
		return []Reservation{}, nil
	}
	return result.Reservations, nil
}

// Get retrieves a single reservation by name.
func Get(ctx context.Context, client *slurmclient.ServiceClient, reservationName string, opts ...GetOpts) (*Reservation, error) {
	query := url.Values{}
	if len(opts) > 0 && opts[0].UpdateTime != nil {
		query.Set("update_time", fmt.Sprintf("%d", *opts[0].UpdateTime))
	}

	resp, err := client.GetWithQuery(ctx, getURL(client, reservationName), query)
	if err != nil {
		return nil, err
	}

	var result api.V0040OpenapiReservationResp
	if err := slurmclient.DecodeResponse(resp, &result); err != nil {
		return nil, err
	}

	if len(result.Reservations) == 0 {
		return nil, fmt.Errorf("reservation %s not found", reservationName)
	}

	return &result.Reservations[0], nil
}
