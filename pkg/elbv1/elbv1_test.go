package elbv1

import (
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	goelb "github.com/aws/aws-sdk-go/service/elb"
	"github.com/kylelemons/godebug/pretty"
	"github.com/yuuki/albio/pkg/awsapi"
	"github.com/yuuki/albio/pkg/model"
)

func TestGetLoadBalancersFromInstanceID(t *testing.T) {
	fakeSvc := &awsapi.FakeELBv1API{
		FakeDescribeLoadBalancers: func(input *goelb.DescribeLoadBalancersInput) (*goelb.DescribeLoadBalancersOutput, error) {
			return &goelb.DescribeLoadBalancersOutput{
				LoadBalancerDescriptions: []*goelb.LoadBalancerDescription{
					{
						LoadBalancerName: aws.String("albiolb01"),
						DNSName:          aws.String("albio-gotest01.ap-northeast-1.elb.amazonaws.com"),
						Instances: []*goelb.Instance{
							{InstanceId: aws.String("albioinstance01-id")},
							{InstanceId: aws.String("albioinstance02-id")},
						},
					},
					{
						LoadBalancerName: aws.String("albiolb02"),
						DNSName:          aws.String("albio-gotest02.ap-northeast-1.elb.amazonaws.com"),
						Instances: []*goelb.Instance{
							{InstanceId: aws.String("albioinstance01-id")},
							{InstanceId: aws.String("albioinstance03-id")},
						},
					},
					{
						LoadBalancerName: aws.String("albiolb03"),
						DNSName:          aws.String("albio-gotest03.ap-northeast-1.elb.amazonaws.com"),
						Instances: []*goelb.Instance{
							{InstanceId: aws.String("albioinstance02-id")},
							{InstanceId: aws.String("albioinstance03-id")},
						},
					},
				},
			}, nil
		},
	}
	obj := &_elbv1{svc: fakeSvc}
	got, err := obj.GetLoadBalancersFromInstanceID("albioinstance01-id")
	if err != nil {
		t.Errorf("should not raise error: %v", err)
	}
	expected := model.LoadBalancers{
		{
			Name:    "albiolb01",
			DNSName: "albio-gotest01.ap-northeast-1.elb.amazonaws.com",
			Instances: []*model.Instance{
				{InstanceID: "albioinstance01-id"},
				{InstanceID: "albioinstance02-id"},
			},
		},
		{
			Name:    "albiolb02",
			DNSName: "albio-gotest02.ap-northeast-1.elb.amazonaws.com",
			Instances: []*model.Instance{
				{InstanceID: "albioinstance01-id"},
				{InstanceID: "albioinstance03-id"},
			},
		},
	}
	if diff := pretty.Compare(got, expected); diff != "" {
		t.Errorf("diff: (-actual +expected)\n%s", diff)
	}
}

func TestGetLoadBalancersByNames(t *testing.T) {
	fakeSvc := &awsapi.FakeELBv1API{
		FakeDescribeLoadBalancers: func(input *goelb.DescribeLoadBalancersInput) (*goelb.DescribeLoadBalancersOutput, error) {
			expected := []*string{aws.String("albiolb01"), aws.String("albiolb02")}
			if diff := pretty.Compare(input.LoadBalancerNames, expected); diff != "" {
				t.Errorf("diff: (-actual +expected)\n%s", diff)
			}
			return &goelb.DescribeLoadBalancersOutput{
				LoadBalancerDescriptions: []*goelb.LoadBalancerDescription{
					{
						LoadBalancerName: aws.String("albiolb01"),
						DNSName:          aws.String("albio-gotest01.ap-northeast-1.elb.amazonaws.com"),
						Instances:        []*goelb.Instance{},
					},
					{
						LoadBalancerName: aws.String("albiolb02"),
						DNSName:          aws.String("albio-gotest02.ap-northeast-1.elb.amazonaws.com"),
						Instances:        []*goelb.Instance{},
					},
				},
			}, nil
		},
	}
	obj := &_elbv1{svc: fakeSvc}
	got, err := obj.GetLoadBalancersByNames([]string{"albiolb01", "albiolb02"})
	if err != nil {
		t.Errorf("should not raise error: %v", err)
	}
	expected := model.LoadBalancers{
		{
			Name:      "albiolb01",
			DNSName:   "albio-gotest01.ap-northeast-1.elb.amazonaws.com",
			Instances: []*model.Instance{},
		},
		{
			Name:      "albiolb02",
			DNSName:   "albio-gotest02.ap-northeast-1.elb.amazonaws.com",
			Instances: []*model.Instance{},
		},
	}
	if diff := pretty.Compare(got, expected); diff != "" {
		t.Errorf("diff: (-actual +expected)\n%s", diff)
	}
}

func TestAddInstanceToLoadBalancers(t *testing.T) {
	fakeSvc := &awsapi.FakeELBv1API{
		FakeRegisterInstancesWithLoadBalancer: func(input *goelb.RegisterInstancesWithLoadBalancerInput) (
			*goelb.RegisterInstancesWithLoadBalancerOutput, error) {
			return &goelb.RegisterInstancesWithLoadBalancerOutput{}, nil
		},
		FakeWaitUntilInstanceInService: func(input *goelb.DescribeInstanceHealthInput) error {
			return nil
		},
	}
	obj := &_elbv1{svc: fakeSvc}
	err := obj.AddInstanceToLoadBalancers("albioinstance01-id", model.LoadBalancers{
		{
			Name:      "albiolb01",
			DNSName:   "albio-gotest01.ap-northeast-1.elb.amazonaws.com",
			Instances: []*model.Instance{},
		},
	})
	if err != nil {
		t.Errorf("should not raise error: %v", err)
	}
}

func TestRemoveInstanceToLoadBalancers(t *testing.T) {
	fakeSvc := &awsapi.FakeELBv1API{
		FakeDeregisterInstancesFromLoadBalancer: func(input *goelb.DeregisterInstancesFromLoadBalancerInput) (
			*goelb.DeregisterInstancesFromLoadBalancerOutput, error) {
			return &goelb.DeregisterInstancesFromLoadBalancerOutput{}, nil
		},
		FakeWaitUntilInstanceDeregistered: func(input *goelb.DescribeInstanceHealthInput) error {
			return nil
		},
	}
	obj := &_elbv1{svc: fakeSvc}
	err := obj.RemoveInstanceFromLoadBalancers("albioinstance01-id", model.LoadBalancers{
		{
			Name:      "albiolb01",
			DNSName:   "albio-gotest01.ap-northeast-1.elb.amazonaws.com",
			Instances: []*model.Instance{},
		},
	})
	if err != nil {
		t.Errorf("should not raise error: %v", err)
	}
}
