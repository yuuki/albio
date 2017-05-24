package command

import (
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/yuuki/albio/pkg/ec2"
	"github.com/yuuki/albio/pkg/elb"
)

func Attach() error {
	sess := session.New()
	ec2Client := ec2.New(sess)
	instanceID, err := ec2Client.GetLocalInstanceID()
	if err != nil {
		return err
	}

	lbNames, err := ec2Client.GetLoadBalancersFromTag(instanceID)
	if err != nil {
		return err
	}

	if len(lbNames) < 1 {
		return fmt.Errorf("not found loadbalancers with %s. specify loadbalancer with --loadlalancer option", instanceID)
	}

	lbClient := elb.New(sess)
	log.Println("-->", "Attaching", instanceID, "to", lbNames)
	if err := lbClient.AddInstanceToLoadBalancers(instanceID, lbNames); err != nil {
		return err
	}

	if err := ec2Client.SaveLoadBalancersToTag(instanceID, lbNames); err != nil {
		return err
	}

	return nil
}
