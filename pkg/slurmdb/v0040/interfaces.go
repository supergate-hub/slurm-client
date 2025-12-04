// SPDX-FileCopyrightText: Copyright (C) SchedMD LLC.
// SPDX-License-Identifier: Apache-2.0

package v0040

import (
	"context"

	"github.com/supergate-hub/slurm-client/pkg/slurmdb/v0040/accounts"
	"github.com/supergate-hub/slurm-client/pkg/slurmdb/v0040/associations"
	"github.com/supergate-hub/slurm-client/pkg/slurmdb/v0040/clusters"
	"github.com/supergate-hub/slurm-client/pkg/slurmdb/v0040/config"
	"github.com/supergate-hub/slurm-client/pkg/slurmdb/v0040/diag"
	"github.com/supergate-hub/slurm-client/pkg/slurmdb/v0040/instances"
	"github.com/supergate-hub/slurm-client/pkg/slurmdb/v0040/jobs"
	"github.com/supergate-hub/slurm-client/pkg/slurmdb/v0040/qos"
	"github.com/supergate-hub/slurm-client/pkg/slurmdb/v0040/tres"
	"github.com/supergate-hub/slurm-client/pkg/slurmdb/v0040/users"
	"github.com/supergate-hub/slurm-client/pkg/slurmdb/v0040/wckeys"
)

// =========================================================================
// Type Re-exports from domain packages
// =========================================================================

// Accounts types
type Account = accounts.Account

// Users types
type User = users.User

// Associations types
type Association = associations.Association

// Clusters types
type Cluster = clusters.Cluster

// QOS types
type QOS = qos.QOS

// Jobs types
type Job = jobs.Job

// WCKeys types
type WCKey = wckeys.WCKey

// TRES types
type TRES = tres.TRES

// Instances types
type Instance = instances.Instance

// Config types
type Config = config.Config

// Diag types
type DiagInfo = diag.DiagInfo

// =========================================================================
// Accounts Interface
// =========================================================================

// AccountsInterface provides operations on Slurm accounts.
type AccountsInterface interface {
	// List retrieves all accounts.
	List(ctx context.Context, opts accounts.ListOpts) ([]Account, error)

	// Get retrieves a single account by name.
	Get(ctx context.Context, accountName string, opts ...accounts.GetOpts) (*Account, error)

	// Create creates one or more accounts.
	Create(ctx context.Context, opts accounts.CreateOpts) error

	// Delete removes an account.
	Delete(ctx context.Context, accountName string) error
}

// =========================================================================
// Users Interface
// =========================================================================

// UsersInterface provides operations on Slurm users.
type UsersInterface interface {
	// List retrieves all users.
	List(ctx context.Context, opts users.ListOpts) ([]User, error)

	// Get retrieves a single user by name.
	Get(ctx context.Context, userName string, opts ...users.GetOpts) (*User, error)

	// Create creates one or more users.
	Create(ctx context.Context, opts users.CreateOpts) error

	// Delete removes a user.
	Delete(ctx context.Context, userName string) error
}

// =========================================================================
// Associations Interface
// =========================================================================

// AssociationsInterface provides operations on Slurm associations.
type AssociationsInterface interface {
	// List retrieves all associations.
	List(ctx context.Context, opts associations.ListOpts) ([]Association, error)

	// Get retrieves a single association.
	Get(ctx context.Context, opts associations.GetOpts) (*Association, error)

	// Create creates one or more associations.
	Create(ctx context.Context, opts associations.CreateOpts) error

	// Delete removes an association.
	Delete(ctx context.Context, opts associations.DeleteOpts) error
}

// =========================================================================
// Clusters Interface
// =========================================================================

// ClustersInterface provides operations on Slurm clusters.
type ClustersInterface interface {
	// List retrieves all clusters.
	List(ctx context.Context, opts clusters.ListOpts) ([]Cluster, error)

	// Get retrieves a single cluster by name.
	Get(ctx context.Context, clusterName string, opts ...clusters.GetOpts) (*Cluster, error)

	// Create creates one or more clusters.
	Create(ctx context.Context, opts clusters.CreateOpts) error

	// Delete removes a cluster.
	Delete(ctx context.Context, clusterName string, opts ...clusters.DeleteOpts) error
}

// =========================================================================
// QOS Interface
// =========================================================================

// QOSInterface provides operations on Slurm QOS.
type QOSInterface interface {
	// List retrieves all QOS.
	List(ctx context.Context, opts qos.ListOpts) ([]QOS, error)

	// Get retrieves a single QOS by name.
	Get(ctx context.Context, qosName string, opts ...qos.GetOpts) (*QOS, error)

	// Create creates one or more QOS.
	Create(ctx context.Context, opts qos.CreateOpts) error

	// Delete removes a QOS.
	Delete(ctx context.Context, qosName string) error
}

// =========================================================================
// Jobs Interface
// =========================================================================

// JobsInterface provides operations on Slurm job history.
type JobsInterface interface {
	// List retrieves jobs from the database.
	List(ctx context.Context, opts jobs.ListOpts) ([]Job, error)

	// Get retrieves a single job by ID.
	Get(ctx context.Context, jobID string) (*Job, error)
}

// =========================================================================
// WCKeys Interface
// =========================================================================

// WCKeysInterface provides operations on Slurm WCKeys.
type WCKeysInterface interface {
	// List retrieves all WCKeys.
	List(ctx context.Context, opts wckeys.ListOpts) ([]WCKey, error)

	// Get retrieves a single WCKey by ID.
	Get(ctx context.Context, wckeyID string) (*WCKey, error)

	// Create creates one or more WCKeys.
	Create(ctx context.Context, opts wckeys.CreateOpts) error

	// Delete removes a WCKey.
	Delete(ctx context.Context, wckeyID string) error
}

// =========================================================================
// TRES Interface
// =========================================================================

// TRESInterface provides operations on Slurm TRES.
type TRESInterface interface {
	// List retrieves all TRES.
	List(ctx context.Context, opts tres.ListOpts) ([]TRES, error)

	// Create creates one or more TRES.
	Create(ctx context.Context, opts tres.CreateOpts) error
}

// =========================================================================
// Instances Interface
// =========================================================================

// InstancesInterface provides operations on Slurm instances.
type InstancesInterface interface {
	// List retrieves all instances.
	List(ctx context.Context, opts instances.ListOpts) ([]Instance, error)

	// Get retrieves a single instance.
	Get(ctx context.Context, opts instances.GetOpts) (*Instance, error)
}

// =========================================================================
// Config Interface
// =========================================================================

// ConfigInterface provides operations on Slurm database configuration.
type ConfigInterface interface {
	// Get retrieves the database configuration.
	Get(ctx context.Context) (*Config, error)

	// Update updates the database configuration.
	Update(ctx context.Context, opts config.UpdateOpts) error
}

// =========================================================================
// Diag Interface
// =========================================================================

// DiagInterface provides operations on Slurm database diagnostics.
type DiagInterface interface {
	// Get retrieves diagnostic information.
	Get(ctx context.Context) (*DiagInfo, error)
}



