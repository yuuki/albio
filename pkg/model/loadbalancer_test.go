package model

import (
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	goelb "github.com/aws/aws-sdk-go/service/elb"
	"github.com/kylelemons/godebug/pretty"
)

func TestNewLoadBalancers(t *testing.T) {
	got := NewLoadBalancers([]*goelb.LoadBalancerDescription{
		{
			LoadBalancerName: aws.String("albio-gotest01"),
			DNSName:          aws.String("albio-gotest01.ap-northeast-1.elb.amazonaws.com"),
			Instances:        []*goelb.Instance{},
		},
		{
			LoadBalancerName: aws.String("albio-gotest02"),
			DNSName:          aws.String("albio-gotest02.ap-northeast-1.elb.amazonaws.com"),
			Instances:        []*goelb.Instance{},
		},
	})
	expected := LoadBalancers{
		{
			Name:      "albio-gotest01",
			DNSName:   "albio-gotest01.ap-northeast-1.elb.amazonaws.com",
			Instances: []*Instance{},
		},
		{
			Name:      "albio-gotest02",
			DNSName:   "albio-gotest02.ap-northeast-1.elb.amazonaws.com",
			Instances: []*Instance{},
		},
	}
	if diff := pretty.Compare(got, expected); diff != "" {
		t.Errorf("diff: (-actual +expected)\n%s", diff)
	}
}

func TestNewLoadBalancersByInstanceID(t *testing.T) {
	got := NewLoadBalancersByInstanceID([]*goelb.LoadBalancerDescription{
		{
			LoadBalancerName: aws.String("albio-gotest01"),
			DNSName:          aws.String("albio-gotest01.ap-northeast-1.elb.amazonaws.com"),
			Instances:        []*goelb.Instance{{InstanceId: aws.String("instance01")}},
		},
		{
			LoadBalancerName: aws.String("albio-gotest02"),
			DNSName:          aws.String("albio-gotest02.ap-northeast-1.elb.amazonaws.com"),
			Instances:        []*goelb.Instance{{InstanceId: aws.String("instance02")}},
		},
	}, "instance02")
	expected := LoadBalancers{
		{
			Name:      "albio-gotest02",
			DNSName:   "albio-gotest02.ap-northeast-1.elb.amazonaws.com",
			Instances: []*Instance{{InstanceID: "instance02"}},
		},
	}
	if diff := pretty.Compare(got, expected); diff != "" {
		t.Errorf("diff: (-actual +expected)\n%s", diff)
	}
}

func TestNewLoadBalancersByInstanceID_NotExistsInstance(t *testing.T) {
	got := NewLoadBalancersByInstanceID([]*goelb.LoadBalancerDescription{
		{
			LoadBalancerName: aws.String("albio-gotest01"),
			DNSName:          aws.String("albio-gotest01.ap-northeast-1.elb.amazonaws.com"),
			Instances:        []*goelb.Instance{{InstanceId: aws.String("instance01")}},
		},
	}, "instance02")
	if len(got) != 0 {
		t.Errorf("len(loadbalancers) = %v; want 0", len(got))
	}
}

func TestNewLoadBalancer(t *testing.T) {
	got := NewLoadBalancer(&goelb.LoadBalancerDescription{
		LoadBalancerName: aws.String("albio-gotest"),
		DNSName:          aws.String("albio-gotest.ap-northeast-1.elb.amazonaws.com"),
		Instances:        []*goelb.Instance{},
	})
	expected := LoadBalancer{
		Name:      "albio-gotest",
		DNSName:   "albio-gotest.ap-northeast-1.elb.amazonaws.com",
		Instances: []*Instance{},
	}
	if diff := pretty.Compare(got, expected); diff != "" {
		t.Errorf("diff: (-actual +expected)\n%s", diff)
	}
}

func TestLoadBalancerString(t *testing.T) {
	l := NewLoadBalancer(&goelb.LoadBalancerDescription{
		LoadBalancerName: aws.String("albio-gotest"),
		DNSName:          aws.String("albio-gotest.ap-northeast-1.elb.amazonaws.com"),
		Instances:        []*goelb.Instance{},
	})
	if got := l.String(); got != "albio-gotest" {
		t.Errorf("model.LoadBalancer.String() = '%v'; want 'albio-gotest'\n", got)
	}
}
