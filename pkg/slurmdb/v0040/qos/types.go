// SPDX-FileCopyrightText: Copyright (C) SchedMD LLC.
// SPDX-License-Identifier: Apache-2.0

package qos

import (
	api "github.com/supergate-hub/slurm-client/api/v0040"
)

// QOS represents a Slurm Quality of Service.
type QOS = api.V0040Qos

// ListOpts specifies the options for listing QOS.
type ListOpts struct {
	// Description to filter by.
	Description *string
	// WithDeleted includes deleted QOS.
	WithDeleted *bool
}

// GetOpts specifies the options for getting a single QOS.
type GetOpts struct {
	// WithDeleted includes deleted QOS.
	WithDeleted *bool
}

// CreateOpts specifies the options for creating QOS.
type CreateOpts struct {
	// QOS list to create.
	QOS []QOS
}

