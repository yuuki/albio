package elb

import (
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/defaults"
	"github.com/aws/aws-sdk-go/aws/ec2metadata"
	"github.com/aws/aws-sdk-go/aws/session"
	goelb "github.com/aws/aws-sdk-go/service/elb"
	"github.com/yuuki/albio/pkg/awsapi"
	"github.com/yuuki/albio/pkg/model"
)

type ELB interface {
	GetLoadBalancersFromInstanceID(string) (model.LoadBalancers, error)
	GetLoadBalancersByNames([]string) (model.LoadBalancers, error)
	AddInstanceToLoadBalancers(string, model.LoadBalancers) error
	RemoveInstanceFromLoadBalancers(string, model.LoadBalancers) error
}

type _elb struct {
	svc awsapi.ELBAPI
}

func New(sess *session.Session) ELB {
	region := os.Getenv("AWS_REGION")
	if region == "" {
		region, _ = ec2metadata.New(sess).Region()
	}
	return &_elb{
		svc: goelb.New(sess,
			defaults.Config().WithRegion(region).WithCredentialsChainVerboseErrors(true),
		),
	}
}

func (e *_elb) GetLoadBalancersFromInstanceID(instanceID string) (model.LoadBalancers, error) {
	resp, err := e.svc.DescribeLoadBalancers(&goelb.DescribeLoadBalancersInput{})
	if err != nil {
		return nil, err
	}
	return model.NewLoadBalancersByInstanceID(resp.LoadBalancerDescriptions, instanceID), nil
}

func (e *_elb) GetLoadBalancersByNames(lbNames []string) (model.LoadBalancers, error) {
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
