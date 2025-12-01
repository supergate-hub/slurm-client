// SPDX-FileCopyrightText: Copyright (C) SchedMD LLC.
// SPDX-License-Identifier: Apache-2.0

package v0039

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	api "github.com/supergate-hub/slurm-client/api/v0039"
	"github.com/supergate-hub/slurm-client/pkg/utils"
)

type NodeInterface interface {
	GetNode(ctx context.Context, nodeName string) (*api.V0039Node, error)
	ListNodes(ctx context.Context) ([]api.V0039Node, error)
}

// GetNode implements NodeInterface.
func (c *SlurmClient) GetNode(ctx context.Context, nodeName string) (*api.V0039Node, error) {
	resp, err := c.SlurmV0039GetNodeWithResponse(ctx, nodeName)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode() != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode())
	}
	if resp.JSON200 == nil {
		return nil, fmt.Errorf("empty response body")
	}
	if resp.JSON200.Errors != nil && len(*resp.JSON200.Errors) > 0 {
		errs := []error{}
		for _, e := range *resp.JSON200.Errors {
			if e.Error != nil {
				errs = append(errs, errors.New(*e.Error))
			}
		}
		if len(errs) > 0 {
			return nil, errors.Join(errs...)
		}
	}
	if resp.JSON200.Nodes == nil || len(*resp.JSON200.Nodes) == 0 {
		return nil, nil
	}
	
	out := &api.V0039Node{}
	utils.RemarshalOrDie((*resp.JSON200.Nodes)[0], out)
	return out, nil
}

// ListNodes implements NodeInterface.
func (c *SlurmClient) ListNodes(ctx context.Context) ([]api.V0039Node, error) {
	resp, err := c.SlurmV0039GetNodesWithResponse(ctx, &api.SlurmV0039GetNodesParams{})
	if err != nil {
		return nil, err
	}
	if resp.StatusCode() != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode())
	}
	if resp.JSON200 == nil {
		return nil, fmt.Errorf("empty response body")
	}
	if resp.JSON200.Errors != nil && len(*resp.JSON200.Errors) > 0 {
		errs := []error{}
		for _, e := range *resp.JSON200.Errors {
			if e.Error != nil {
				errs = append(errs, errors.New(*e.Error))
			}
		}
		if len(errs) > 0 {
			return nil, errors.Join(errs...)
		}
	}
	
	nodes := []api.V0039Node{}
	if resp.JSON200.Nodes != nil {
		for _, item := range *resp.JSON200.Nodes {
			node := api.V0039Node{}
			utils.RemarshalOrDie(item, &node)
			nodes = append(nodes, node)
		}
	}
	return nodes, nil
}
