package slurmclient

import (
	"context"

	v0040 "github.com/supergate-hub/slurm-client/api/v0040"
)

// --- Reservations ---

type reservationsServiceImpl struct{ t *Transport }

func (s *reservationsServiceImpl) List(ctx context.Context) ([]v0040.V0040ReservationInfo, error) {
	var resp v0040.V0040OpenapiReservationResp
	if err := s.t.Get(ctx, "/reservations", &resp); err != nil {
		return nil, err
	}
	return resp.Reservations, nil
}

func (s *reservationsServiceImpl) Get(ctx context.Context, name string) (*v0040.V0040ReservationInfo, error) {
	var resp v0040.V0040OpenapiReservationResp
	if err := s.t.Get(ctx, "/reservation/"+name, &resp); err != nil {
		return nil, err
	}
	if len(resp.Reservations) == 0 {
		return nil, ErrNotFound
	}
	return &resp.Reservations[0], nil
}

// --- Licenses ---

type licensesServiceImpl struct{ t *Transport }

func (s *licensesServiceImpl) List(ctx context.Context) ([]v0040.V0040License, error) {
	var resp v0040.V0040OpenapiLicensesResp
	if err := s.t.Get(ctx, "/licenses", &resp); err != nil {
		return nil, err
	}
	return resp.Licenses, nil
}

// --- Shares ---

type sharesServiceImpl struct{ t *Transport }

func (s *sharesServiceImpl) Get(ctx context.Context) (*v0040.V0040OpenapiSharesResp, error) {
	var resp v0040.V0040OpenapiSharesResp
	if err := s.t.Get(ctx, "/shares", &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// --- Diag ---

type diagServiceImpl struct{ t *Transport }

func (s *diagServiceImpl) Get(ctx context.Context) (*v0040.V0040StatsMsg, error) {
	var resp v0040.V0040OpenapiDiagResp
	if err := s.t.Get(ctx, "/diag", &resp); err != nil {
		return nil, err
	}
	return &resp.Statistics, nil
}

// --- Ping ---

type pingServiceImpl struct{ t *Transport }

func (s *pingServiceImpl) Get(ctx context.Context) ([]v0040.V0040ControllerPing, error) {
	var resp v0040.V0040OpenapiPingArrayResp
	if err := s.t.Get(ctx, "/ping", &resp); err != nil {
		return nil, err
	}
	return resp.Pings, nil
}

// --- Reconfigure ---

type reconfigureServiceImpl struct{ t *Transport }

func (s *reconfigureServiceImpl) Reconfigure(ctx context.Context) error {
	return s.t.Get(ctx, "/reconfigure", nil)
}
