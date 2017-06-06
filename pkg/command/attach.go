package command

import (
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/yuuki/albio/pkg/ec2"
	"github.com/yuuki/albio/pkg/elb"
)

type AttachParam struct {
	InstanceID string
}

func Attach(param *AttachParam) error {
	sess := session.New()
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

	if len(lbNames) < 1 {
		return fmt.Errorf("not found loadbalancers with %s. specify loadbalancer with --loadlalancer option", instanceID)
	}

	lbClient := elb.New(sess)

	lbs, err := lbClient.GetLoadBalancersByNames(lbNames)
	if err != nil {
		return err
	}

	log.Println("-->", "Attaching", instanceID, "to", lbs)
	if err := lbClient.AddInstanceToLoadBalancers(instanceID, lbs); err != nil {
		return err
	}

	if err := ec2Client.SaveLoadBalancersToTag(instanceID, lbNames); err != nil {
		return err
	}

	return nil
}
