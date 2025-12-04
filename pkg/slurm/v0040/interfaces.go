// SPDX-FileCopyrightText: Copyright (C) SchedMD LLC.
// SPDX-License-Identifier: Apache-2.0

package v0040

import (
	"context"

	"github.com/supergate-hub/slurm-client/pkg/slurm/v0040/diag"
	"github.com/supergate-hub/slurm-client/pkg/slurm/v0040/jobs"
	"github.com/supergate-hub/slurm-client/pkg/slurm/v0040/licenses"
	"github.com/supergate-hub/slurm-client/pkg/slurm/v0040/nodes"
	"github.com/supergate-hub/slurm-client/pkg/slurm/v0040/partitions"
	"github.com/supergate-hub/slurm-client/pkg/slurm/v0040/ping"
	"github.com/supergate-hub/slurm-client/pkg/slurm/v0040/reservations"
	"github.com/supergate-hub/slurm-client/pkg/slurm/v0040/shares"
)

// =========================================================================
// Type Re-exports from domain packages
// =========================================================================

// Jobs types
type (
	Job               = jobs.Job
	JobSubmitResponse = jobs.JobSubmitResponse
	JobArrayResponse  = jobs.JobArrayResponse
)

// Nodes types
type Node = nodes.Node

// Partitions types
type Partition = partitions.Partition

// Reservations types
type Reservation = reservations.Reservation

// Licenses types
type License = licenses.License

// Shares types
type ShareInfo = shares.ShareInfo

// Diag types
type DiagInfo = diag.DiagInfo

// Ping types
type PingInfo = ping.PingInfo

// =========================================================================
// Jobs Interface
// =========================================================================

// JobsInterface provides operations on Slurm jobs.
type JobsInterface interface {
	// List retrieves all jobs.
	List(ctx context.Context, opts jobs.ListOpts) ([]Job, error)

	// Get retrieves a single job by ID.
	Get(ctx context.Context, jobID string, opts ...jobs.GetOpts) (*Job, error)

	// Submit submits a new job.
	Submit(ctx context.Context, opts jobs.SubmitOpts) (*JobSubmitResponse, error)

	// Update modifies an existing job.
	Update(ctx context.Context, jobID string, opts jobs.UpdateOpts) ([]JobArrayResponse, error)

	// Delete cancels/removes a job.
	Delete(ctx context.Context, jobID string, opts ...jobs.DeleteOpts) error
}

// =========================================================================
// Nodes Interface
// =========================================================================

// NodesInterface provides operations on Slurm nodes.
type NodesInterface interface {
	// List retrieves all nodes.
	List(ctx context.Context, opts nodes.ListOpts) ([]Node, error)

	// Get retrieves a single node by name.
	Get(ctx context.Context, nodeName string, opts ...nodes.GetOpts) (*Node, error)

	// Update modifies a node's properties.
	Update(ctx context.Context, nodeName string, opts nodes.UpdateOpts) error

	// Delete removes a node from Slurm's configuration.
	Delete(ctx context.Context, nodeName string) error
}

// =========================================================================
// Partitions Interface
// =========================================================================

// PartitionsInterface provides operations on Slurm partitions.
type PartitionsInterface interface {
	// List retrieves all partitions.
	List(ctx context.Context, opts partitions.ListOpts) ([]Partition, error)

	// Get retrieves a single partition by name.
	Get(ctx context.Context, partitionName string, opts ...partitions.GetOpts) (*Partition, error)
}

// =========================================================================
// Reservations Interface
// =========================================================================

// ReservationsInterface provides operations on Slurm reservations.
type ReservationsInterface interface {
	// List retrieves all reservations.
	List(ctx context.Context, opts reservations.ListOpts) ([]Reservation, error)

	// Get retrieves a single reservation by name.
	Get(ctx context.Context, reservationName string, opts ...reservations.GetOpts) (*Reservation, error)
}

// =========================================================================
// Licenses Interface
// =========================================================================

// LicensesInterface provides operations on Slurm licenses.
type LicensesInterface interface {
	// List retrieves all licenses.
	List(ctx context.Context, opts licenses.ListOpts) ([]License, error)
}

// =========================================================================
// Shares Interface
// =========================================================================

// SharesInterface provides operations on Slurm fair share information.
type SharesInterface interface {
	// Get retrieves fair share information.
	Get(ctx context.Context, opts shares.GetOpts) ([]ShareInfo, error)
}

// =========================================================================
// Diag Interface
// =========================================================================

// DiagInterface provides operations on Slurm diagnostics.
type DiagInterface interface {
	// Get retrieves diagnostic information.
	Get(ctx context.Context) (*DiagInfo, error)
}

// =========================================================================
// Ping Interface
// =========================================================================

// PingInterface provides operations on Slurm controller ping.
type PingInterface interface {
	// Ping checks controller status.
	Ping(ctx context.Context) ([]PingInfo, error)
}

// =========================================================================
// Reconfigure Interface
// =========================================================================

// ReconfigureInterface provides operations on Slurm reconfiguration.
type ReconfigureInterface interface {
	// Trigger initiates a reconfiguration.
	Trigger(ctx context.Context) error
}


