package licenses

import (
	api "github.com/supergate-hub/slurm-client/api/v0040"
)

// License represents a Slurm license.
type License = api.V0040License

// ListOpts specifies the options for listing licenses.
type ListOpts struct {
	// Currently no options are supported for listing licenses.
}

