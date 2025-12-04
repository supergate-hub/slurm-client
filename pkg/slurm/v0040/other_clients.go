// SPDX-FileCopyrightText: Copyright (C) SchedMD LLC.
// SPDX-License-Identifier: Apache-2.0

package v0040

import (
	"context"

	"github.com/supergate-hub/slurm-client/pkg/slurm/v0040/diag"
	"github.com/supergate-hub/slurm-client/pkg/slurm/v0040/licenses"
	"github.com/supergate-hub/slurm-client/pkg/slurm/v0040/ping"
	"github.com/supergate-hub/slurm-client/pkg/slurm/v0040/reconfigure"
	"github.com/supergate-hub/slurm-client/pkg/slurm/v0040/reservations"
	"github.com/supergate-hub/slurm-client/pkg/slurm/v0040/shares"
	"github.com/supergate-hub/slurm-client/pkg/slurmclient"
)

// =========================================================================
// Reservations Client
// =========================================================================

// reservationsClient implements ReservationsInterface.
type reservationsClient struct {
	client *slurmclient.ServiceClient
}

var _ ReservationsInterface = (*reservationsClient)(nil)

// List retrieves all reservations.
func (c *reservationsClient) List(ctx context.Context, opts reservations.ListOpts) ([]Reservation, error) {
	return reservations.List(ctx, c.client, opts)
}

// Get retrieves a single reservation by name.
func (c *reservationsClient) Get(ctx context.Context, reservationName string, opts ...reservations.GetOpts) (*Reservation, error) {
	return reservations.Get(ctx, c.client, reservationName, opts...)
}

// =========================================================================
// Licenses Client
// =========================================================================

// licensesClient implements LicensesInterface.
type licensesClient struct {
	client *slurmclient.ServiceClient
}

var _ LicensesInterface = (*licensesClient)(nil)

// List retrieves all licenses.
func (c *licensesClient) List(ctx context.Context, opts licenses.ListOpts) ([]License, error) {
	return licenses.List(ctx, c.client, opts)
}

// =========================================================================
// Shares Client
// =========================================================================

// sharesClient implements SharesInterface.
type sharesClient struct {
	client *slurmclient.ServiceClient
}

var _ SharesInterface = (*sharesClient)(nil)

// Get retrieves fair share information.
func (c *sharesClient) Get(ctx context.Context, opts shares.GetOpts) ([]ShareInfo, error) {
	return shares.Get(ctx, c.client, opts)
}

// =========================================================================
// Diag Client
// =========================================================================

// diagClient implements DiagInterface.
type diagClient struct {
	client *slurmclient.ServiceClient
}

var _ DiagInterface = (*diagClient)(nil)

// Get retrieves diagnostic information.
func (c *diagClient) Get(ctx context.Context) (*DiagInfo, error) {
	return diag.Get(ctx, c.client)
}

// =========================================================================
// Ping Client
// =========================================================================

// pingClient implements PingInterface.
type pingClient struct {
	client *slurmclient.ServiceClient
}

var _ PingInterface = (*pingClient)(nil)

// Ping checks controller status.
func (c *pingClient) Ping(ctx context.Context) ([]PingInfo, error) {
	return ping.Ping(ctx, c.client)
}

// =========================================================================
// Reconfigure Client
// =========================================================================

// reconfigureClient implements ReconfigureInterface.
type reconfigureClient struct {
	client *slurmclient.ServiceClient
}

var _ ReconfigureInterface = (*reconfigureClient)(nil)

// Trigger initiates a reconfiguration.
func (c *reconfigureClient) Trigger(ctx context.Context) error {
	return reconfigure.Trigger(ctx, c.client)
}



