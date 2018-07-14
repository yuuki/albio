package command

import (
	"fmt"
	"log"

	"github.com/yuuki/albio/pkg/alb"
	"github.com/yuuki/albio/pkg/awsapi"
	"github.com/yuuki/albio/pkg/ec2"
	"github.com/yuuki/albio/pkg/elb"
)

type AttachParam struct {
	InstanceID       string
	LoadBalancerName string
}

func Attach(param *AttachParam) error {
	sess, err := awsapi.Session()
	if err != nil {
		return err
	}

	ec2Client := ec2.New(sess)

	var instanceID string
	if param.InstanceID == "" {
		var err error
		instanceID, err = ec2Client.GetLocalInstanceID()
		if err != nil {
			return err
		}
	}

	lbNames, err := ec2Client.GetLoadBalancerNamesFromTag(instanceID)
	if err != nil {
		return err
	}

	if param.LoadBalancerName == "" {
		if len(lbNames) < 1 {
			return fmt.Errorf("not found loadbalancers with %s. specify loadbalancer with --loadlalancer option", instanceID)
		}

		elbClient := elb.New(sess)
		elbs, err := elbClient.GetLoadBalancersByNames(lbNames)
		if err != nil {
			return err
		}
		log.Println("-->", "Attaching", instanceID, "to", elbs)
		if err := elbClient.AddInstanceToLoadBalancers(instanceID, elbs); err != nil {
			return err
		}

		albClient := alb.New(sess)
		albs, err := albClient.GetLoadBalancersByNames(lbNames)
		if err != nil {
			return err
		}
		log.Println("-->", "Attaching", instanceID, "to", albs)
		if err := albClient.AddInstanceToLoadBalancers(instanceID, albs); err != nil {
			return err
		}
	} else {
		elbClient := elb.New(sess)
		elbs, err := elbClient.GetLoadBalancersByNames([]string{param.LoadBalancerName})
		if err != nil {
			return err
		}
		log.Println("-->", "Attaching", instanceID, "to", elbs)
		if err := elbClient.AddInstanceToLoadBalancers(instanceID, elbs); err != nil {
			return err
		}

		albClient := alb.New(sess)
		albs, err := albClient.GetLoadBalancersByNames([]string{param.LoadBalancerName})
		if err != nil {
			return err
		}
		log.Println("-->", "Attaching", instanceID, "to", albs)
		if err := albClient.AddInstanceToLoadBalancers(instanceID, albs); err != nil {
			return err
		}

		lbNames = append(lbNames, param.LoadBalancerName)
	}

	if err := ec2Client.SaveLoadBalancersToTag(instanceID, lbNames); err != nil {
		return err
	}

	return nil
}
