package model

import (
	goelb "github.com/aws/aws-sdk-go/service/elb"
)

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
