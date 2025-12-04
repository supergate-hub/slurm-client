// SPDX-FileCopyrightText: Copyright (C) SchedMD LLC.
// SPDX-License-Identifier: Apache-2.0

package wckeys

import (
	api "github.com/supergate-hub/slurm-client/api/v0040"
)

// WCKey represents a Slurm WCKey.
type WCKey = api.V0040Wckey

// ListOpts specifies the options for listing WCKeys.
type ListOpts struct {
	// WithDeleted includes deleted WCKeys.
	WithDeleted *bool
}

// CreateOpts specifies the options for creating WCKeys.
type CreateOpts struct {
	// WCKeys to create.
	WCKeys []WCKey
}

