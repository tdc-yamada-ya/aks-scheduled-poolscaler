package scaler

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/profiles/latest/containerservice/mgmt/containerservice"
)

type PoolUpdater interface {
	UpdatePool(ctx context.Context, resourceGroupName ResourceGroupName, resourceName ResourceName, agentPoolName AgentPoolName, parameters *Parameters) error
}

type AzPoolUpdater struct {
	Client containerservice.AgentPoolsClient
}

func (u *AzPoolUpdater) UpdatePool(ctx context.Context, resourceGroupName ResourceGroupName, resourceName ResourceName, agentPoolName AgentPoolName, parameters *Parameters) error {
	ap, err := u.Client.Get(ctx, string(resourceGroupName), string(resourceName), string(agentPoolName))
	if err != nil {
		return err
	}

	ap.EnableAutoScaling = parameters.EnableAutoScaling
	ap.Count = parameters.Count
	ap.MinCount = parameters.MinCount
	ap.MaxCount = parameters.MaxCount

	_, err = u.Client.CreateOrUpdate(ctx, string(resourceGroupName), string(resourceName), string(agentPoolName), ap)
	if err != nil {
		return err
	}

	return nil
}
