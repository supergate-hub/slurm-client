package jobs

import (
	api "github.com/supergate-hub/slurm-client/api/v0040"
)

// Job represents a Slurm job. This is an alias to the generated API type.
type Job = api.V0040JobInfo

// JobDescMsg represents a Slurm job description.
type JobDescMsg = api.V0040JobDescMsg

// JobSubmitResponse represents the response from a job submission.
type JobSubmitResponse = api.V0040JobSubmitResponseMsg

// JobArrayResponse represents a response for array job operations.
type JobArrayResponse = api.V0040JobArrayResponseMsgEntry

// ListOpts specifies the options for listing jobs.
type ListOpts struct {
	// UpdateTime filters jobs updated after this time (Unix timestamp).
	UpdateTime *int64
}

// GetOpts specifies the options for getting a single job.
type GetOpts struct {
	// UpdateTime filters if the job was updated after this time.
	UpdateTime *int64
}

// SubmitOpts specifies the options for submitting a new job.
type SubmitOpts struct {
	// Script is the batch script content.
	Script string

	// JobDescMsg contains the full job description.
	// If provided, this takes precedence over individual fields.
	Job *JobDescMsg

	// Jobs contains multiple job descriptions for het jobs or array jobs.
	Jobs []JobDescMsg
}

// UpdateOpts specifies the options for updating a job.
type UpdateOpts struct {
	// Job contains the job properties to update.
	Job *JobDescMsg
}

// DeleteOpts specifies the options for deleting a job.
type DeleteOpts struct {
	// Signal to send to the job (optional).
	Signal *string
	// Flags for the delete operation.
	Flags *[]string
}


