package command

import (
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws/session"

	"github.com/yuuki/albio/pkg/ec2"
	"github.com/yuuki/albio/pkg/elb"
)

func Detach() error {
	sess := session.New()
	ec2Client := ec2.New(sess)
	instanceID, err := ec2Client.GetLocalInstanceID()
	if err != nil {
		return err
	}

	lbClient := elb.New(sess)
	lbNames, err := lbClient.GetLoadBalancersFromInstanceID(instanceID)
	if err != nil {
		return err
	}
	if len(lbNames) < 1 {
		return fmt.Errorf("%s is not attached any loadbalancers")
	}

	if err := ec2Client.SaveLoadBalancersToTag(instanceID, lbNames); err != nil {
		return err
	}

	log.Println("-->", "Detaching", instanceID, "from", lbNames)
	if err := lbClient.RemoveInstanceFromLoadBalancers(instanceID, lbNames); err != nil {
		return err
	}

	return nil
}
