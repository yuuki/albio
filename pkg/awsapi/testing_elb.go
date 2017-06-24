package awsapi

import (
	goelb "github.com/aws/aws-sdk-go/service/elb"
)

type FakeELBAPI struct {
	ELBAPI
	FakeDescribeLoadBalancers               func(input *goelb.DescribeLoadBalancersInput) (*goelb.DescribeLoadBalancersOutput, error)
	FakeRegisterInstancesWithLoadBalancer   func(input *goelb.RegisterInstancesWithLoadBalancerInput) (*goelb.RegisterInstancesWithLoadBalancerOutput, error)
	FakeWaitUntilInstanceInService          func(input *goelb.DescribeInstanceHealthInput) error
	FakeDeregisterInstancesFromLoadBalancer func(input *goelb.DeregisterInstancesFromLoadBalancerInput) (*goelb.DeregisterInstancesFromLoadBalancerOutput, error)
	FakeWaitUntilInstanceDeregistered       func(input *goelb.DescribeInstanceHealthInput) error
}

func (e *FakeELBAPI) DescribeLoadBalancers(input *goelb.DescribeLoadBalancersInput) (*goelb.DescribeLoadBalancersOutput, error) {
	return e.FakeDescribeLoadBalancers(input)
}

func (e *FakeELBAPI) RegisterInstancesWithLoadBalancer(input *goelb.RegisterInstancesWithLoadBalancerInput) (
	*goelb.RegisterInstancesWithLoadBalancerOutput, error) {
	return e.FakeRegisterInstancesWithLoadBalancer(input)
}

func (e *FakeELBAPI) WaitUntilInstanceInService(input *goelb.DescribeInstanceHealthInput) error {
	return e.FakeWaitUntilInstanceInService(input)
}

func (e *FakeELBAPI) DeregisterInstancesFromLoadBalancer(input *goelb.DeregisterInstancesFromLoadBalancerInput) (
	*goelb.DeregisterInstancesFromLoadBalancerOutput, error) {
	return e.FakeDeregisterInstancesFromLoadBalancer(input)
}

func (e *FakeELBAPI) WaitUntilInstanceDeregistered(input *goelb.DescribeInstanceHealthInput) error {
	return e.FakeWaitUntilInstanceDeregistered(input)
}
