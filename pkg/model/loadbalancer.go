package model

import (
	"strings"

	goelb "github.com/aws/aws-sdk-go/service/elb"
)

// LoadBalancers represents an slice of loadbalancer.
type LoadBalancers []*LoadBalancer

// NewLoadBalancers creates the object of LoadBalancers from the slice of elb.LoadBalancerDescription.
func NewLoadBalancers(descs []*goelb.LoadBalancerDescription) LoadBalancers {
	var lbs LoadBalancers
	for _, desc := range descs {
		lbs = append(lbs, NewLoadBalancer(desc))
	}
	return lbs
}

// NewLoadBalancersByInstanceID creates the object of LoadBalancers from elb.LoadBalancerDescription.
func NewLoadBalancersByInstanceID(descs []*goelb.LoadBalancerDescription, instanceID string) LoadBalancers {
	var lbs LoadBalancers
	for _, desc := range descs {
		for _, instance := range desc.Instances {
			if *instance.InstanceId == instanceID {
				lbs = append(lbs, NewLoadBalancer(desc))
			}
		}
	}
	return lbs
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

// LoadBalancer represents a loadbalancer.
type LoadBalancer struct {
	Name      string      `json:"name"`
	DNSName   string      `json:"dnsname"`
	Instances []*Instance `json:"instances"`
}

// NewLoadBalancer creates a LoadBalancer object from elb.LoadBalancerDescription.
func NewLoadBalancer(desc *goelb.LoadBalancerDescription) *LoadBalancer {
	instances := make([]*Instance, 0, len(desc.Instances))
	for _, instance := range desc.Instances {
		instances = append(instances, NewInstance(instance))
	}
	return &LoadBalancer{
		Name:      *desc.LoadBalancerName,
		DNSName:   *desc.DNSName,
		Instances: instances,
	}
}

// String returns a string reprentation of LoadBalancer.
func (l *LoadBalancer) String() string {
	return l.Name
}