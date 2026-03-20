package slurmclient

import (
	"context"

	v0040 "github.com/supergate-hub/slurm-client/api/v0040"
)

// --- Accounts ---

type accountsServiceImpl struct{ t *Transport }

func (s *accountsServiceImpl) List(ctx context.Context) ([]v0040.V0040Account, error) {
	var resp v0040.V0040OpenapiAccountsResp
	if err := s.t.Get(ctx, "/accounts", &resp); err != nil {
		return nil, err
	}
	return resp.Accounts, nil
}

func (s *accountsServiceImpl) Get(ctx context.Context, name string) (*v0040.V0040Account, error) {
	var resp v0040.V0040OpenapiAccountsResp
	if err := s.t.Get(ctx, "/account/"+name, &resp); err != nil {
		return nil, err
	}
	if len(resp.Accounts) == 0 {
		return nil, ErrNotFound
	}
	return &resp.Accounts[0], nil
}

// --- Users ---

type usersServiceImpl struct{ t *Transport }

func (s *usersServiceImpl) List(ctx context.Context) ([]v0040.V0040User, error) {
	var resp v0040.V0040OpenapiUsersResp
	if err := s.t.Get(ctx, "/users", &resp); err != nil {
		return nil, err
	}
	return resp.Users, nil
}

func (s *usersServiceImpl) Get(ctx context.Context, name string) (*v0040.V0040User, error) {
	var resp v0040.V0040OpenapiUsersResp
	if err := s.t.Get(ctx, "/user/"+name, &resp); err != nil {
		return nil, err
	}
	if len(resp.Users) == 0 {
		return nil, ErrNotFound
	}
	return &resp.Users[0], nil
}

// --- Associations ---

type associationsServiceImpl struct{ t *Transport }

func (s *associationsServiceImpl) List(ctx context.Context) ([]v0040.V0040Assoc, error) {
	var resp v0040.V0040OpenapiAssocsResp
	if err := s.t.Get(ctx, "/associations", &resp); err != nil {
		return nil, err
	}
	return resp.Associations, nil
}

// --- Clusters ---

type clustersServiceImpl struct{ t *Transport }

func (s *clustersServiceImpl) List(ctx context.Context) ([]v0040.V0040ClusterRec, error) {
	var resp v0040.V0040OpenapiClustersResp
	if err := s.t.Get(ctx, "/clusters", &resp); err != nil {
		return nil, err
	}
	return resp.Clusters, nil
}

func (s *clustersServiceImpl) Get(ctx context.Context, name string) (*v0040.V0040ClusterRec, error) {
	var resp v0040.V0040OpenapiClustersResp
	if err := s.t.Get(ctx, "/cluster/"+name, &resp); err != nil {
		return nil, err
	}
	if len(resp.Clusters) == 0 {
		return nil, ErrNotFound
	}
	return &resp.Clusters[0], nil
}

// --- QOS ---

type qosServiceImpl struct{ t *Transport }

func (s *qosServiceImpl) List(ctx context.Context) ([]v0040.V0040Qos, error) {
	var resp v0040.V0040OpenapiSlurmdbdQosResp
	if err := s.t.Get(ctx, "/qos", &resp); err != nil {
		return nil, err
	}
	return resp.Qos, nil
}

func (s *qosServiceImpl) Get(ctx context.Context, name string) (*v0040.V0040Qos, error) {
	var resp v0040.V0040OpenapiSlurmdbdQosResp
	if err := s.t.Get(ctx, "/qos/"+name, &resp); err != nil {
		return nil, err
	}
	if len(resp.Qos) == 0 {
		return nil, ErrNotFound
	}
	return &resp.Qos[0], nil
}

// --- SlurmDB Jobs ---

type slurmdbJobsServiceImpl struct{ t *Transport }

func (s *slurmdbJobsServiceImpl) List(ctx context.Context) ([]v0040.V0040Job, error) {
	var resp v0040.V0040OpenapiSlurmdbdJobsResp
	if err := s.t.Get(ctx, "/jobs", &resp); err != nil {
		return nil, err
	}
	return resp.Jobs, nil
}

func (s *slurmdbJobsServiceImpl) Get(ctx context.Context, jobID string) (*v0040.V0040Job, error) {
	var resp v0040.V0040OpenapiSlurmdbdJobsResp
	if err := s.t.Get(ctx, "/job/"+jobID, &resp); err != nil {
		return nil, err
	}
	if len(resp.Jobs) == 0 {
		return nil, ErrNotFound
	}
	return &resp.Jobs[0], nil
}

// --- WCKeys ---

type wckeysServiceImpl struct{ t *Transport }

func (s *wckeysServiceImpl) List(ctx context.Context) ([]v0040.V0040Wckey, error) {
	var resp v0040.V0040OpenapiWckeyResp
	if err := s.t.Get(ctx, "/wckeys", &resp); err != nil {
		return nil, err
	}
	return resp.Wckeys, nil
}

// --- TRES ---

type tresServiceImpl struct{ t *Transport }

func (s *tresServiceImpl) List(ctx context.Context) ([]v0040.V0040Tres, error) {
	var resp v0040.V0040OpenapiTresResp
	if err := s.t.Get(ctx, "/tres", &resp); err != nil {
		return nil, err
	}
	return resp.TRES, nil
}

// --- Instances ---

type instancesServiceImpl struct{ t *Transport }

func (s *instancesServiceImpl) List(ctx context.Context) ([]v0040.V0040Instance, error) {
	var resp v0040.V0040OpenapiInstancesResp
	if err := s.t.Get(ctx, "/instances", &resp); err != nil {
		return nil, err
	}
	return resp.Instances, nil
}

// --- Config ---

type configServiceImpl struct{ t *Transport }

func (s *configServiceImpl) Get(ctx context.Context) (*v0040.V0040OpenapiSlurmdbdConfigResp, error) {
	var resp v0040.V0040OpenapiSlurmdbdConfigResp
	if err := s.t.Get(ctx, "/config", &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// --- SlurmDB Diag ---

type slurmdbDiagServiceImpl struct{ t *Transport }

func (s *slurmdbDiagServiceImpl) Get(ctx context.Context) (*v0040.V0040OpenapiSlurmdbdStatsResp, error) {
	var resp v0040.V0040OpenapiSlurmdbdStatsResp
	if err := s.t.Get(ctx, "/diag", &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}
