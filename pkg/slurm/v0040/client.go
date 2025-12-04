// SPDX-FileCopyrightText: Copyright (C) SchedMD LLC.
// SPDX-License-Identifier: Apache-2.0

// Package v0040 provides a client-go style interface for Slurm REST API v0.0.40.
//
// Usage:
//
//	cfg := slurmclient.Config{
//	    Endpoint:  "http://slurmrestd:6820",
//	    AuthToken: "<jwt token>",
//	}
//
//	cli, err := v0040.NewClient(cfg)
//	if err != nil {
//	    return err
//	}
//
//	// Use domain-specific interfaces
//	jobs, _ := cli.Jobs.List(ctx, jobs.ListOpts{})
//	nodes, _ := cli.Nodes.List(ctx, nodes.ListOpts{})
//	pings, _ := cli.Ping.Ping(ctx)
package v0040

import (
	"github.com/supergate-hub/slurm-client/pkg/slurmclient"
)

// Client is the main entry point for Slurm REST API v0.0.40.
// It provides access to all Slurm resources through domain-specific interfaces.
type Client struct {
	// serviceClient is the underlying HTTP client.
	serviceClient *slurmclient.ServiceClient

	// Jobs provides access to job operations.
	Jobs JobsInterface

	// Nodes provides access to node operations.
	Nodes NodesInterface

	// Partitions provides access to partition operations.
	Partitions PartitionsInterface

	// Reservations provides access to reservation operations.
	Reservations ReservationsInterface

	// Licenses provides access to license operations.
	Licenses LicensesInterface

	// Shares provides access to fair share operations.
	Shares SharesInterface

	// Diag provides access to diagnostic operations.
	Diag DiagInterface

	// Ping provides access to controller ping operations.
	Ping PingInterface

	// Reconfigure provides access to reconfigure operations.
	Reconfigure ReconfigureInterface
}

// NewClient creates a new Slurm client for API version v0.0.40.
func NewClient(cfg slurmclient.Config) (*Client, error) {
	serviceClient, err := slurmclient.NewSlurmClient(cfg, "v0.0.40")
	if err != nil {
		return nil, err
	}

	c := &Client{serviceClient: serviceClient}

	// Initialize all domain interfaces
	c.Jobs = &jobsClient{client: serviceClient}
	c.Nodes = &nodesClient{client: serviceClient}
	c.Partitions = &partitionsClient{client: serviceClient}
	c.Reservations = &reservationsClient{client: serviceClient}
	c.Licenses = &licensesClient{client: serviceClient}
	c.Shares = &sharesClient{client: serviceClient}
	c.Diag = &diagClient{client: serviceClient}
	c.Ping = &pingClient{client: serviceClient}
	c.Reconfigure = &reconfigureClient{client: serviceClient}

	return c, nil
}

// ServiceClient returns the underlying ServiceClient for advanced use cases.
func (c *Client) ServiceClient() *slurmclient.ServiceClient {
	return c.serviceClient
}


