package nodes

import (
	api "github.com/supergate-hub/slurm-client/api/v0040"
)

// Node represents a Slurm node.
type Node = api.V0040Node

// NodeState represents a Slurm node state.
type NodeState = api.V0040NodeState

// ListOpts specifies the options for listing nodes.
type ListOpts struct {
	// UpdateTime filters nodes updated after this time.
	UpdateTime *int64
}

// GetOpts specifies the options for getting a single node.
type GetOpts struct {
	// UpdateTime filters if the node was updated after this time.
	UpdateTime *int64
}

// UpdateOpts specifies the options for updating a node.
type UpdateOpts struct {
	// Comment to set on the node.
	Comment *string
	// Extra arbitrary string to set on the node.
	Extra *string
	// State to set the node to.
	State *[]NodeState
	// Reason for the state change.
	Reason *string
	// User who is making the change.
	ReasonUID *string
}

