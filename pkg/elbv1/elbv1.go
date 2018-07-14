package elbv1

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	goelb "github.com/aws/aws-sdk-go/service/elb"
	"github.com/yuuki/albio/pkg/awsapi"
	"github.com/yuuki/albio/pkg/model"
)

type ELBv1 interface {
	GetLoadBalancersFromInstanceID(string) (model.LoadBalancers, error)
	GetLoadBalancersByNames([]string) (model.LoadBalancers, error)
	AddInstanceToLoadBalancers(string, model.LoadBalancers) error
	RemoveInstanceFromLoadBalancers(string, model.LoadBalancers) error
}

type _elbv1 struct {
	svc awsapi.ELBv1API
}

func New(sess *session.Session) ELBv1 {
	return &_elbv1{
		svc: goelb.New(sess),
	}
}

func (e *_elbv1) GetLoadBalancersFromInstanceID(instanceID string) (model.LoadBalancers, error) {
	resp, err := e.svc.DescribeLoadBalancers(&goelb.DescribeLoadBalancersInput{})
	if err != nil {
		return nil, err
	}
	return model.NewLoadBalancersByInstanceIDFromELBv1(resp.LoadBalancerDescriptions, instanceID), nil
}

func (e *_elbv1) GetLoadBalancersByNames(lbNames []string) (model.LoadBalancers, error) {
	names := make([]*string, 0, len(lbNames))
	for _, n := range lbNames {
		names = append(names, aws.String(n))
	}
	resp, err := e.svc.DescribeLoadBalancers(&goelb.DescribeLoadBalancersInput{
		LoadBalancerNames: names,
	})
	if err != nil {
		return nil, err
	}
	return model.NewLoadBalancers(resp.LoadBalancerDescriptions), nil
}

func (e *_elbv1) AddInstanceToLoadBalancers(instanceID string, lbs model.LoadBalancers) error {
	for _, lb := range lbs {
		_, err := e.svc.RegisterInstancesWithLoadBalancer(
			&goelb.RegisterInstancesWithLoadBalancerInput{
				Instances:        []*goelb.Instance{{InstanceId: aws.String(instanceID)}},
				LoadBalancerName: aws.String(lb.Name),
			},
		)
		if err != nil {
			return err
		}
		err = e.svc.WaitUntilInstanceInService(&goelb.DescribeInstanceHealthInput{
			Instances:        []*goelb.Instance{&goelb.Instance{InstanceId: aws.String(instanceID)}},
			LoadBalancerName: aws.String(lb.Name),
		})
		if err != nil {
			return err
		}
	}
	return nil
}

func (e *_elbv1) RemoveInstanceFromLoadBalancers(instanceID string, lbs model.LoadBalancers) error {
	for _, lb := range lbs {
		_, err := e.svc.DeregisterInstancesFromLoadBalancer(
			&goelb.DeregisterInstancesFromLoadBalancerInput{
				Instances:        []*goelb.Instance{{InstanceId: aws.String(instanceID)}},
				LoadBalancerName: aws.String(lb.Name),
			},
		)
		if err != nil {
			return err
		}
		err = e.svc.WaitUntilInstanceDeregistered(&goelb.DescribeInstanceHealthInput{
			Instances:        []*goelb.Instance{&goelb.Instance{InstanceId: aws.String(instanceID)}},
			LoadBalancerName: aws.String(lb.Name),
		})
		if err != nil {
			return err
		}
	}
	return nil
}
