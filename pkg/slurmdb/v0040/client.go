// SPDX-FileCopyrightText: Copyright (C) SchedMD LLC.
// SPDX-License-Identifier: Apache-2.0

// Package v0040 provides a client-go style interface for Slurmdb REST API v0.0.40.
//
// Usage:
//
//	cfg := slurmclient.Config{
//	    Endpoint:  "http://slurmdbd:6820",
//	    AuthToken: "<jwt token>",
//	}
//
//	cli, err := v0040.NewClient(cfg)
//	if err != nil {
//	    return err
//	}
//
//	// Use domain-specific interfaces
//	accounts, _ := cli.Accounts.List(ctx, accounts.ListOpts{})
//	users, _ := cli.Users.List(ctx, users.ListOpts{})
package v0040

import (
	"github.com/supergate-hub/slurm-client/pkg/slurmclient"
)

// Client is the main entry point for Slurmdb REST API v0.0.40.
// It provides access to all Slurmdb resources through domain-specific interfaces.
type Client struct {
	// serviceClient is the underlying HTTP client.
	serviceClient *slurmclient.ServiceClient

	// Accounts provides access to account operations.
	Accounts AccountsInterface

	// Users provides access to user operations.
	Users UsersInterface

	// Associations provides access to association operations.
	Associations AssociationsInterface

	// Clusters provides access to cluster operations.
	Clusters ClustersInterface

	// QOS provides access to QOS operations.
	QOS QOSInterface

	// Jobs provides access to job history operations.
	Jobs JobsInterface

	// WCKeys provides access to WCKey operations.
	WCKeys WCKeysInterface

	// TRES provides access to TRES operations.
	TRES TRESInterface

	// Instances provides access to instance operations.
	Instances InstancesInterface

	// Config provides access to configuration operations.
	Config ConfigInterface

	// Diag provides access to diagnostic operations.
	Diag DiagInterface
}

// NewClient creates a new Slurmdb client for API version v0.0.40.
func NewClient(cfg slurmclient.Config) (*Client, error) {
	serviceClient, err := slurmclient.NewSlurmdbClient(cfg, "v0.0.40")
	if err != nil {
		return nil, err
	}

	c := &Client{serviceClient: serviceClient}

	// Initialize all domain interfaces
	c.Accounts = &accountsClient{client: serviceClient}
	c.Users = &usersClient{client: serviceClient}
	c.Associations = &associationsClient{client: serviceClient}
	c.Clusters = &clustersClient{client: serviceClient}
	c.QOS = &qosClient{client: serviceClient}
	c.Jobs = &jobsClient{client: serviceClient}
	c.WCKeys = &wckeysClient{client: serviceClient}
	c.TRES = &tresClient{client: serviceClient}
	c.Instances = &instancesClient{client: serviceClient}
	c.Config = &configClient{client: serviceClient}
	c.Diag = &diagClient{client: serviceClient}

	return c, nil
}

// ServiceClient returns the underlying ServiceClient for advanced use cases.
func (c *Client) ServiceClient() *slurmclient.ServiceClient {
	return c.serviceClient
}



