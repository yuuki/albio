package model

import (
	"strings"

	"github.com/aws/aws-sdk-go/service/elbv2"
)

// LoadBalancers represents an slice of loadbalancer.
type LoadBalancers []*LoadBalancer

// NewLoadBalancersFromELBv2 creates the object of LoadBalancers from ELBv2.
func NewLoadBalancersFromELBv2(loadBalancers []*elbv2.LoadBalancer,
	loadBalancerArnToTargets map[string][]*elbv2.TargetDescription) LoadBalancers {
	models := make(LoadBalancers, 0, len(loadBalancers))
	for _, lb := range loadBalancers {
		targets := loadBalancerArnToTargets[*lb.LoadBalancerArn]
		if len(targets) == 0 {
			continue
		}
		models = append(models, NewLoadBalancerFromELBv2(lb, targets))
	}
	return models
}

// String returns a string reprentation of LoadBalancers.
func (lbs LoadBalancers) String() string {
	return strings.Join(lbs.NameSlice(), ",")
}

// NameSlice returns a slice of loadbalancer's name.
func (lbs LoadBalancers) NameSlice() []string {
	names := make([]string, 0, len(lbs))
	for _, lb := range lbs {
		names = append(names, lb.Name)
	}
	return names
}

// NamePointerSlice returns a slice of loadbalancer's name.
func (lbs LoadBalancers) NamePointerSlice() []*string {
	names := make([]*string, 0, len(lbs))
	for _, lb := range lbs {
		name := lb.Name
		names = append(names, &name)
	}
	return names
}

// LoadBalancer represents a loadbalancer.
type LoadBalancer struct {
	Name      string      `json:"name"`
	DNSName   string      `json:"dnsname"`
	Arn       string      `json:"arn"`
	Type      string      `json:"type"`
	Scheme    string      `json:"scheme"`
	Instances []*Instance `json:"instances"`
}

func NewLoadBalancerFromELBv2(desc *elbv2.LoadBalancer, targets []*elbv2.TargetDescription) *LoadBalancer {
	instances := make([]*Instance, 0, len(targets))
	for _, target := range targets {
		instances = append(instances, NewInstanceFromTarget(target))
	}
	return &LoadBalancer{
		Name:      *desc.LoadBalancerName,
		DNSName:   *desc.DNSName,
		Arn:       *desc.LoadBalancerArn,
		Type:      *desc.Type,
		Scheme:    *desc.Scheme,
		Instances: instances,
	}
}

// String returns a string reprentation of LoadBalancer.
func (l *LoadBalancer) String() string {
	return l.Name
}
