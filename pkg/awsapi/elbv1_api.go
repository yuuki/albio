package awsapi

import (
	goelb "github.com/aws/aws-sdk-go/service/elb"
)

// ELBv1API defines the interface for ELB API stubbing.
type ELBv1API interface {
	DescribeLoadBalancers(*goelb.DescribeLoadBalancersInput) (*goelb.DescribeLoadBalancersOutput, error)
	RegisterInstancesWithLoadBalancer(*goelb.RegisterInstancesWithLoadBalancerInput) (*goelb.RegisterInstancesWithLoadBalancerOutput, error)
	WaitUntilInstanceInService(*goelb.DescribeInstanceHealthInput) error
	DeregisterInstancesFromLoadBalancer(*goelb.DeregisterInstancesFromLoadBalancerInput) (*goelb.DeregisterInstancesFromLoadBalancerOutput, error)
	WaitUntilInstanceDeregistered(*goelb.DescribeInstanceHealthInput) error
}
