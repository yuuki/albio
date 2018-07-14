package awsapi

import (
	goelb "github.com/aws/aws-sdk-go/service/elb"
)

type FakeELBv1API struct {
	ELBv1API
	FakeDescribeLoadBalancers               func(input *goelb.DescribeLoadBalancersInput) (*goelb.DescribeLoadBalancersOutput, error)
	FakeRegisterInstancesWithLoadBalancer   func(input *goelb.RegisterInstancesWithLoadBalancerInput) (*goelb.RegisterInstancesWithLoadBalancerOutput, error)
	FakeWaitUntilInstanceInService          func(input *goelb.DescribeInstanceHealthInput) error
	FakeDeregisterInstancesFromLoadBalancer func(input *goelb.DeregisterInstancesFromLoadBalancerInput) (*goelb.DeregisterInstancesFromLoadBalancerOutput, error)
	FakeWaitUntilInstanceDeregistered       func(input *goelb.DescribeInstanceHealthInput) error
}

func (e *FakeELBv1API) DescribeLoadBalancers(input *goelb.DescribeLoadBalancersInput) (*goelb.DescribeLoadBalancersOutput, error) {
	return e.FakeDescribeLoadBalancers(input)
}

func (e *FakeELBv1API) RegisterInstancesWithLoadBalancer(input *goelb.RegisterInstancesWithLoadBalancerInput) (
	*goelb.RegisterInstancesWithLoadBalancerOutput, error) {
	return e.FakeRegisterInstancesWithLoadBalancer(input)
}

func (e *FakeELBv1API) WaitUntilInstanceInService(input *goelb.DescribeInstanceHealthInput) error {
	return e.FakeWaitUntilInstanceInService(input)
}

func (e *FakeELBv1API) DeregisterInstancesFromLoadBalancer(input *goelb.DeregisterInstancesFromLoadBalancerInput) (
	*goelb.DeregisterInstancesFromLoadBalancerOutput, error) {
	return e.FakeDeregisterInstancesFromLoadBalancer(input)
}

func (e *FakeELBv1API) WaitUntilInstanceDeregistered(input *goelb.DescribeInstanceHealthInput) error {
	return e.FakeWaitUntilInstanceDeregistered(input)
}
