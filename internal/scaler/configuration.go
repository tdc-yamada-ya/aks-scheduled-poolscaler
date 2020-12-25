package scaler

import (
	"fmt"
	"time"

	"gopkg.in/yaml.v2"
)

type (
	Parameters struct {
		EnableAutoScaling *bool  `yaml:"enableAutoScaling"`
		Count             *int32 `yaml:"count"`
		MinCount          *int32 `yaml:"minCount"`
		MaxCount          *int32 `yaml:"maxCount"`
	}

	ParametersName string

	Rule struct {
		Expression     Expression     `yaml:"expr"`
		ParametersName ParametersName `yaml:"paramsRef"`
	}

	ResourceGroupName string
	ResourceName      string
	AgentPoolName     string
	Rules             []Rule

	Resource struct {
		ResourceGroupName ResourceGroupName `yaml:"resourceGroupName"`
		ResourceName      ResourceName      `yaml:"resourceName"`
		AgentPoolName     AgentPoolName     `yaml:"agentPoolName"`
		Rules             Rules             `yaml:"rules"`
	}

	ParametersMap map[ParametersName]*Parameters
	Resources     []Resource

	Configuration struct {
		ParametersDefinitions ParametersMap `yaml:"paramsDefs"`
		Resources             Resources     `yaml:"resources"`
	}
)

func (rules Rules) FindFirst(t time.Time) (*Rule, bool) {
	for _, rule := range rules {
		if rule.Expression.Match(t) {
			return &rule, true
		}
	}
	return nil, false
}

func FindParameters(rs Rules, t time.Time, m ParametersMap) (*Parameters, bool) {
	r, ok := rs.FindFirst(t)
	if !ok {
		return nil, false
	}

	p, ok := m[r.ParametersName]
	return p, ok
}

func UnmarshalConfigurationWithYaml(b []byte) (*Configuration, error) {
	c := Configuration{}
	err := yaml.Unmarshal(b, &c)
	if err != nil {
		return nil, err
	}

	return &c, nil
}

func (p *Parameters) String() string {
	bpsf := func(p *bool) string {
		if p == nil {
			return "nil"
		}
		return fmt.Sprintf("%t", *p)
	}
	ipsf := func(p *int32) string {
		if p == nil {
			return "nil"
		}
		return fmt.Sprintf("%d", *p)
	}
	return fmt.Sprintf("EnableAutoScaling: %s, Count: %s, MinCount: %s, MaxCount: %s", bpsf(p.EnableAutoScaling), ipsf(p.Count), ipsf(p.MinCount), ipsf(p.MaxCount))
}
