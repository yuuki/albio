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
	lbs, err := lbClient.GetLoadBalancersFromInstanceID(instanceID)
	if err != nil {
		return err
	}
	if len(lbs) < 1 {
		return fmt.Errorf("%s is not attached any loadbalancers")
	}

	lbNames := make([]string, 0, len(lbs))
	for _, n := range lbs {
		lbNames = append(lbNames, n.Name)
	}
	if err := ec2Client.SaveLoadBalancersToTag(instanceID, lbNames); err != nil {
		return err
	}

	log.Println("-->", "Detaching", instanceID, "from", lbs)
	if err := lbClient.RemoveInstanceFromLoadBalancers(instanceID, lbs); err != nil {
		return err
	}

	return nil
}
