// SPDX-FileCopyrightText: Copyright (C) SchedMD LLC.
// SPDX-License-Identifier: Apache-2.0

package v0040

import (
	"context"

	"github.com/supergate-hub/slurm-client/pkg/slurmclient"
	"github.com/supergate-hub/slurm-client/pkg/slurmdb/v0040/associations"
	"github.com/supergate-hub/slurm-client/pkg/slurmdb/v0040/clusters"
	"github.com/supergate-hub/slurm-client/pkg/slurmdb/v0040/config"
	"github.com/supergate-hub/slurm-client/pkg/slurmdb/v0040/diag"
	"github.com/supergate-hub/slurm-client/pkg/slurmdb/v0040/instances"
	"github.com/supergate-hub/slurm-client/pkg/slurmdb/v0040/jobs"
	"github.com/supergate-hub/slurm-client/pkg/slurmdb/v0040/qos"
	"github.com/supergate-hub/slurm-client/pkg/slurmdb/v0040/tres"
	"github.com/supergate-hub/slurm-client/pkg/slurmdb/v0040/wckeys"
)

// =========================================================================
// Associations Client
// =========================================================================

// associationsClient implements AssociationsInterface.
type associationsClient struct {
	client *slurmclient.ServiceClient
}

var _ AssociationsInterface = (*associationsClient)(nil)

// List retrieves all associations.
func (c *associationsClient) List(ctx context.Context, opts associations.ListOpts) ([]Association, error) {
	return associations.List(ctx, c.client, opts)
}

// Get retrieves a single association.
func (c *associationsClient) Get(ctx context.Context, opts associations.GetOpts) (*Association, error) {
	return associations.Get(ctx, c.client, opts)
}

// Create creates one or more associations.
func (c *associationsClient) Create(ctx context.Context, opts associations.CreateOpts) error {
	return associations.Create(ctx, c.client, opts)
}

// Delete removes an association.
func (c *associationsClient) Delete(ctx context.Context, opts associations.DeleteOpts) error {
	return associations.Delete(ctx, c.client, opts)
}

// =========================================================================
// Clusters Client
// =========================================================================

// clustersClient implements ClustersInterface.
type clustersClient struct {
	client *slurmclient.ServiceClient
}

var _ ClustersInterface = (*clustersClient)(nil)

// List retrieves all clusters.
func (c *clustersClient) List(ctx context.Context, opts clusters.ListOpts) ([]Cluster, error) {
	return clusters.List(ctx, c.client, opts)
}

// Get retrieves a single cluster by name.
func (c *clustersClient) Get(ctx context.Context, clusterName string, opts ...clusters.GetOpts) (*Cluster, error) {
	return clusters.Get(ctx, c.client, clusterName, opts...)
}

// Create creates one or more clusters.
func (c *clustersClient) Create(ctx context.Context, opts clusters.CreateOpts) error {
	return clusters.Create(ctx, c.client, opts)
}

// Delete removes a cluster.
func (c *clustersClient) Delete(ctx context.Context, clusterName string, opts ...clusters.DeleteOpts) error {
	return clusters.Delete(ctx, c.client, clusterName, opts...)
}

// =========================================================================
// QOS Client
// =========================================================================

// qosClient implements QOSInterface.
type qosClient struct {
	client *slurmclient.ServiceClient
}

var _ QOSInterface = (*qosClient)(nil)

// List retrieves all QOS.
func (c *qosClient) List(ctx context.Context, opts qos.ListOpts) ([]QOS, error) {
	return qos.List(ctx, c.client, opts)
}

// Get retrieves a single QOS by name.
func (c *qosClient) Get(ctx context.Context, qosName string, opts ...qos.GetOpts) (*QOS, error) {
	return qos.Get(ctx, c.client, qosName, opts...)
}

// Create creates one or more QOS.
func (c *qosClient) Create(ctx context.Context, opts qos.CreateOpts) error {
	return qos.Create(ctx, c.client, opts)
}

// Delete removes a QOS.
func (c *qosClient) Delete(ctx context.Context, qosName string) error {
	return qos.Delete(ctx, c.client, qosName)
}

// =========================================================================
// Jobs Client
// =========================================================================

// jobsClient implements JobsInterface.
type jobsClient struct {
	client *slurmclient.ServiceClient
}

var _ JobsInterface = (*jobsClient)(nil)

// List retrieves jobs from the database.
func (c *jobsClient) List(ctx context.Context, opts jobs.ListOpts) ([]Job, error) {
	return jobs.List(ctx, c.client, opts)
}

// Get retrieves a single job by ID.
func (c *jobsClient) Get(ctx context.Context, jobID string) (*Job, error) {
	return jobs.Get(ctx, c.client, jobID)
}

// =========================================================================
// WCKeys Client
// =========================================================================

// wckeysClient implements WCKeysInterface.
type wckeysClient struct {
	client *slurmclient.ServiceClient
}

var _ WCKeysInterface = (*wckeysClient)(nil)

// List retrieves all WCKeys.
func (c *wckeysClient) List(ctx context.Context, opts wckeys.ListOpts) ([]WCKey, error) {
	return wckeys.List(ctx, c.client, opts)
}

// Get retrieves a single WCKey by ID.
func (c *wckeysClient) Get(ctx context.Context, wckeyID string) (*WCKey, error) {
	return wckeys.Get(ctx, c.client, wckeyID)
}

// Create creates one or more WCKeys.
func (c *wckeysClient) Create(ctx context.Context, opts wckeys.CreateOpts) error {
	return wckeys.Create(ctx, c.client, opts)
}

// Delete removes a WCKey.
func (c *wckeysClient) Delete(ctx context.Context, wckeyID string) error {
	return wckeys.Delete(ctx, c.client, wckeyID)
}

// =========================================================================
// TRES Client
// =========================================================================

// tresClient implements TRESInterface.
type tresClient struct {
	client *slurmclient.ServiceClient
}

var _ TRESInterface = (*tresClient)(nil)

// List retrieves all TRES.
func (c *tresClient) List(ctx context.Context, opts tres.ListOpts) ([]TRES, error) {
	return tres.List(ctx, c.client, opts)
}

// Create creates one or more TRES.
func (c *tresClient) Create(ctx context.Context, opts tres.CreateOpts) error {
	return tres.Create(ctx, c.client, opts)
}

// =========================================================================
// Instances Client
// =========================================================================

// instancesClient implements InstancesInterface.
type instancesClient struct {
	client *slurmclient.ServiceClient
}

var _ InstancesInterface = (*instancesClient)(nil)

// List retrieves all instances.
func (c *instancesClient) List(ctx context.Context, opts instances.ListOpts) ([]Instance, error) {
	return instances.List(ctx, c.client, opts)
}

// Get retrieves a single instance.
func (c *instancesClient) Get(ctx context.Context, opts instances.GetOpts) (*Instance, error) {
	return instances.Get(ctx, c.client, opts)
}

// =========================================================================
// Config Client
// =========================================================================

// configClient implements ConfigInterface.
type configClient struct {
	client *slurmclient.ServiceClient
}

var _ ConfigInterface = (*configClient)(nil)

// Get retrieves the database configuration.
func (c *configClient) Get(ctx context.Context) (*Config, error) {
	return config.Get(ctx, c.client)
}

// Update updates the database configuration.
func (c *configClient) Update(ctx context.Context, opts config.UpdateOpts) error {
	return config.Update(ctx, c.client, opts)
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



