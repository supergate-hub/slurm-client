
package reservations

import (
	api "github.com/supergate-hub/slurm-client/api/v0040"
)

// Reservation represents a Slurm reservation.
type Reservation = api.V0040ReservationInfo

// ListOpts specifies the options for listing reservations.
type ListOpts struct {
	// UpdateTime filters reservations updated after this time.
	UpdateTime *int64
}

// GetOpts specifies the options for getting a single reservation.
type GetOpts struct {
	// UpdateTime filters if the reservation was updated after this time.
	UpdateTime *int64
}

