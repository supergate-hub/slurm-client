package nodes

import (
	"context"
	"fmt"
	"net/url"

	api "github.com/supergate-hub/slurm-client/api/v0040"
	"github.com/supergate-hub/slurm-client/pkg/slurmclient"
)

// List retrieves all nodes from the Slurm cluster.
func List(ctx context.Context, client *slurmclient.ServiceClient, opts ListOpts) ([]Node, error) {
	query := url.Values{}
	if opts.UpdateTime != nil {
		query.Set("update_time", fmt.Sprintf("%d", *opts.UpdateTime))
	}

	resp, err := client.GetWithQuery(ctx, listURL(client), query)
	if err != nil {
		return nil, err
	}

	var result api.V0040OpenapiNodesResp
	if err := slurmclient.DecodeResponse(resp, &result); err != nil {
		return nil, err
	}

	if len(result.Nodes) == 0 {
		return []Node{}, nil
	}
	return result.Nodes, nil
}

// Get retrieves a single node by name.
func Get(ctx context.Context, client *slurmclient.ServiceClient, nodeName string, opts ...GetOpts) (*Node, error) {
	query := url.Values{}
	if len(opts) > 0 && opts[0].UpdateTime != nil {
		query.Set("update_time", fmt.Sprintf("%d", *opts[0].UpdateTime))
	}

	resp, err := client.GetWithQuery(ctx, getURL(client, nodeName), query)
	if err != nil {
		return nil, err
	}

	var result api.V0040OpenapiNodesResp
	if err := slurmclient.DecodeResponse(resp, &result); err != nil {
		return nil, err
	}

	if len(result.Nodes) == 0 {
		return nil, fmt.Errorf("node %s not found", nodeName)
	}

	return &result.Nodes[0], nil
}

// Update modifies a node's properties.
func Update(ctx context.Context, client *slurmclient.ServiceClient, nodeName string, opts UpdateOpts) error {
	reqBody := api.V0040UpdateNodeMsg{
		Comment:   opts.Comment,
		Extra:     opts.Extra,
		Reason:    opts.Reason,
		ReasonUid: opts.ReasonUID,
	}

	// Convert State if provided
	if opts.State != nil {
		states := make([]api.V0040UpdateNodeMsgState, len(*opts.State))
		for i, s := range *opts.State {
			states[i] = api.V0040UpdateNodeMsgState(s)
		}
		reqBody.State = &states
	}

	resp, err := client.Post(ctx, updateURL(client, nodeName), reqBody)
	if err != nil {
		return err
	}

	return slurmclient.CheckResponse(resp)
}

// Delete removes a node from Slurm's configuration.
func Delete(ctx context.Context, client *slurmclient.ServiceClient, nodeName string) error {
	resp, err := client.Delete(ctx, deleteURL(client, nodeName))
	if err != nil {
		return err
	}

	return slurmclient.CheckResponse(resp)
}
