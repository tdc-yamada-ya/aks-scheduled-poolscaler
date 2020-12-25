package scaler

import (
	"context"
	"log"
	"time"
)

type Scaler struct {
	Logger        *log.Logger
	PoolUpdater   PoolUpdater
	Configuration *Configuration
}

func (s *Scaler) Scale(t time.Time) error {
	for _, r := range s.Configuration.Resources {
		s.Logger.Printf("start scaling - resourceGroupName: %v, resourceName: %v, agentPoolName: %v", r.ResourceGroupName, r.ResourceName, r.AgentPoolName)

		err := s.scaleWithResource(&r, t)
		if err != nil {
			s.Logger.Printf("scale error: %s", err)
			return err
		}
	}
	return nil
}

func (s *Scaler) scaleWithResource(r *Resource, t time.Time) error {
	p, ok := FindParameters(r.Rules, t, s.Configuration.ParametersDefinitions)
	if !ok {
		return nil
	}

	s.Logger.Printf("start updating agent pool - parameters: %v", p)

	return s.PoolUpdater.UpdatePool(context.Background(), r.ResourceGroupName, r.ResourceName, r.AgentPoolName, p)
}
