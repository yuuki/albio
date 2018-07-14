package awsapi

import (
	"github.com/aws/aws-sdk-go/service/elbv2"
)

// ELBv2API defines the interface for ELB v2 API stubbing.
type ELBv2API interface {
	DescribeTargetGroups(*elbv2.DescribeTargetGroupsInput) (*elbv2.DescribeTargetGroupsOutput, error)
	DescribeTargetHealth(*elbv2.DescribeTargetHealthInput) (*elbv2.DescribeTargetHealthOutput, error)
	DescribeLoadBalancers(*elbv2.DescribeLoadBalancersInput) (*elbv2.DescribeLoadBalancersOutput, error)
	RegisterTargets(*elbv2.RegisterTargetsInput) (*elbv2.RegisterTargetsOutput, error)
	DeregisterTargets(*elbv2.DeregisterTargetsInput) (*elbv2.DeregisterTargetsOutput, error)
}
