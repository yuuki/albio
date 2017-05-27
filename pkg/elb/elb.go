package elb

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	goelb "github.com/aws/aws-sdk-go/service/elb"
	"github.com/aws/aws-sdk-go/service/elb/elbiface"
	"github.com/yuuki/albio/pkg/model"
)

type ELB interface {
	GetLoadBalancersFromInstanceID(string) (model.LoadBalancers, error)
	GetLoadBalancersByNames([]string) (model.LoadBalancers, error)
	AddInstanceToLoadBalancers(string, model.LoadBalancers) error
	RemoveInstanceFromLoadBalancers(string, model.LoadBalancers) error
}

type _elb struct {
	svc elbiface.ELBAPI
}

func New(sess *session.Session) ELB {
	return &_elb{
		svc: goelb.New(sess),
	}
}

func (e *_elb) GetLoadBalancersFromInstanceID(instanceID string) (model.LoadBalancers, error) {
	resp, err := e.svc.DescribeLoadBalancers(&goelb.DescribeLoadBalancersInput{})
	if err != nil {
		return nil, err
	}
	var lbs model.LoadBalancers
	for _, lbdesc := range resp.LoadBalancerDescriptions {
		for _, instance := range lbdesc.Instances {
			if *instance.InstanceId == instanceID {
				lbs = append(lbs, model.NewLoadBalancer(lbdesc))
			}
		}
	}
	return lbs, nil
}

func (e *_elb) GetLoadBalancersByNames(lbNames []string) (model.LoadBalancers, error) {
	names := make([]*string, 0, len(lbNames))
	for _, n := range lbNames {
		names = append(names, &n)
	}
	resp, err := e.svc.DescribeLoadBalancers(&goelb.DescribeLoadBalancersInput{
		LoadBalancerNames: names,
	})
	if err != nil {
		return nil, err
	}
	var lbs model.LoadBalancers
	for _, lbdesc := range resp.LoadBalancerDescriptions {
		lbs = append(lbs, model.NewLoadBalancer(lbdesc))
	}
	return lbs, nil
}

func (e *_elb) AddInstanceToLoadBalancers(instanceID string, lbs model.LoadBalancers) error {
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

func (e *_elb) RemoveInstanceFromLoadBalancers(instanceID string, lbs model.LoadBalancers) error {
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
