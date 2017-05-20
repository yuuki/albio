package elb

import (
	"github.com/aws/aws-sdk-go/aws/session"
	goelb "github.com/aws/aws-sdk-go/service/elb"
	"github.com/aws/aws-sdk-go/service/elb/elbiface"
)

type ELB interface {
	GetLoadBalancersFromInstanceID(string) ([]string, error)
}

type _elb struct {
	svc elbiface.ELBAPI
}

func New(sess *session.Session) ELB {
	return &_elb{
		svc: goelb.New(sess),
	}
}

func (e *_elb) GetLoadBalancersFromInstanceID(instanceID string) ([]string, error) {
	resp, err := e.svc.DescribeLoadBalancers(&goelb.DescribeLoadBalancersInput{})
	if err != nil {
		return nil, err
	}
	var lbNames []string
	for _, lb := range resp.LoadBalancerDescriptions {
		for _, instance := range lb.Instances {
			if *instance.InstanceId == instanceID {
				lbNames = append(lbNames, *lb.LoadBalancerName)
			}
		}
	}
	return lbNames, nil
}
