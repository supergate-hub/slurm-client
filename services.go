package slurmclient

import (
	"context"
	"fmt"

	v0040 "github.com/supergate-hub/slurm-client/api/v0040"
)

// --- Slurm Service Interfaces ---

// JobsService provides CRUD operations for Slurm jobs.
type JobsService interface {
	List(ctx context.Context, opts ...JobListOpts) ([]v0040.V0040JobInfo, error)
	Get(ctx context.Context, jobID string) (*v0040.V0040JobInfo, error)
	Submit(ctx context.Context, opts JobSubmitOpts) (*JobSubmitResult, error)
	Update(ctx context.Context, jobID string, opts JobUpdateOpts) error
	Delete(ctx context.Context, jobID string) error
}

// NodesService provides operations for Slurm nodes.
type NodesService interface {
	List(ctx context.Context, opts ...NodeListOpts) ([]v0040.V0040Node, error)
	Get(ctx context.Context, nodeName string) (*v0040.V0040Node, error)
	Delete(ctx context.Context, nodeName string) error
}

// PartitionsService provides operations for Slurm partitions.
type PartitionsService interface {
	List(ctx context.Context) ([]v0040.V0040PartitionInfo, error)
	Get(ctx context.Context, partitionName string) (*v0040.V0040PartitionInfo, error)
}

// ReservationsService provides operations for Slurm reservations.
type ReservationsService interface {
	List(ctx context.Context) ([]v0040.V0040ReservationInfo, error)
	Get(ctx context.Context, reservationName string) (*v0040.V0040ReservationInfo, error)
}

// LicensesService provides operations for Slurm licenses.
type LicensesService interface {
	List(ctx context.Context) ([]v0040.V0040License, error)
}

// SharesService provides operations for Slurm fair-share info.
type SharesService interface {
	Get(ctx context.Context) (*v0040.V0040OpenapiSharesResp, error)
}

// DiagService provides operations for Slurm diagnostics.
type DiagService interface {
	Get(ctx context.Context) (*v0040.V0040StatsMsg, error)
}

// PingService provides operations for Slurm controller ping.
type PingService interface {
	Get(ctx context.Context) ([]v0040.V0040ControllerPing, error)
}

// ReconfigureService provides operations for Slurm reconfigure.
type ReconfigureService interface {
	Reconfigure(ctx context.Context) error
}

// --- SlurmDB Service Interfaces ---

// AccountsService provides operations for SlurmDB accounts.
type AccountsService interface {
	List(ctx context.Context) ([]v0040.V0040Account, error)
	Get(ctx context.Context, accountName string) (*v0040.V0040Account, error)
}

// UsersService provides operations for SlurmDB users.
type UsersService interface {
	List(ctx context.Context) ([]v0040.V0040User, error)
	Get(ctx context.Context, userName string) (*v0040.V0040User, error)
}

// AssociationsService provides operations for SlurmDB associations.
type AssociationsService interface {
	List(ctx context.Context) ([]v0040.V0040Assoc, error)
}

// ClustersService provides operations for SlurmDB clusters.
type ClustersService interface {
	List(ctx context.Context) ([]v0040.V0040ClusterRec, error)
	Get(ctx context.Context, clusterName string) (*v0040.V0040ClusterRec, error)
}

// QOSService provides operations for SlurmDB QOS.
type QOSService interface {
	List(ctx context.Context) ([]v0040.V0040Qos, error)
	Get(ctx context.Context, qosName string) (*v0040.V0040Qos, error)
}

// SlurmdbJobsService provides operations for SlurmDB job accounting.
type SlurmdbJobsService interface {
	List(ctx context.Context) ([]v0040.V0040Job, error)
	Get(ctx context.Context, jobID string) (*v0040.V0040Job, error)
}

// WCKeysService provides operations for SlurmDB WCKeys.
type WCKeysService interface {
	List(ctx context.Context) ([]v0040.V0040Wckey, error)
}

// TRESService provides operations for SlurmDB TRES.
type TRESService interface {
	List(ctx context.Context) ([]v0040.V0040Tres, error)
}

// InstancesService provides operations for SlurmDB instances.
type InstancesService interface {
	List(ctx context.Context) ([]v0040.V0040Instance, error)
}

// ConfigService provides operations for SlurmDB configuration.
type ConfigService interface {
	Get(ctx context.Context) (*v0040.V0040OpenapiSlurmdbdConfigResp, error)
}

// SlurmdbDiagService provides operations for SlurmDB diagnostics.
type SlurmdbDiagService interface {
	Get(ctx context.Context) (*v0040.V0040OpenapiSlurmdbdStatsResp, error)
}

// --- Option & Result Types (in root package for single-import DX) ---

// JobListOpts configures job listing.
type JobListOpts struct {
	UpdateTime *int64
}

// NodeListOpts configures node listing.
type NodeListOpts struct {
	UpdateTime *int64
}

// JobSubmitOpts configures job submission.
type JobSubmitOpts struct {
	// Script is the job batch script. Required.
	Script string
	// Job is the job description. Optional.
	Job *v0040.V0040JobDescMsg
	// Jobs is a list of job descriptions for heterogeneous jobs. Optional.
	Jobs []v0040.V0040JobDescMsg
}

// Validate checks that the submit options are valid.
func (o JobSubmitOpts) Validate() error {
	if o.Script == "" {
		return fmt.Errorf("%w: script is required for job submission", ErrValidation)
	}
	return nil
}

// JobUpdateOpts configures job update.
type JobUpdateOpts struct {
	Job *v0040.V0040JobDescMsg
}

// JobSubmitResult contains the result of a job submission.
type JobSubmitResult struct {
	JobID   int64
	StepID  string
	Warning string
}
