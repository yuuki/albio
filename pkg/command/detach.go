package command

import (
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws/session"

	"github.com/yuuki/albio/pkg/ec2"
	"github.com/yuuki/albio/pkg/elb"
)

type DetachParam struct {
	InstanceID string
}

func Detach(param *DetachParam) error {
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

	lbClient := elb.New(sess)
	lbs, err := lbClient.GetLoadBalancersFromInstanceID(instanceID)
	if err != nil {
		return err
	}
	if len(lbs) < 1 {
		return fmt.Errorf("%s is not attached any loadbalancers")
	}

	if err := ec2Client.SaveLoadBalancersToTag(instanceID, lbs.NameSlice()); err != nil {
		return err
	}

	log.Println("-->", "Detaching", instanceID, "from", lbs)
	if err := lbClient.RemoveInstanceFromLoadBalancers(instanceID, lbs); err != nil {
		return err
	}

	return nil
}
