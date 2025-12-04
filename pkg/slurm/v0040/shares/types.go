
package shares

import (
	api "github.com/supergate-hub/slurm-client/api/v0040"
)

// ShareInfo represents fair share information for an association.
type ShareInfo = api.V0040AssocSharesObjWrap

// GetOpts specifies the options for getting share information.
type GetOpts struct {
	// Accounts to filter by.
	Accounts []string
	// Users to filter by.
	Users []string
}

