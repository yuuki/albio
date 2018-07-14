package command

import (
	"fmt"
	"log"

	"github.com/yuuki/albio/pkg/awsapi"
	"github.com/yuuki/albio/pkg/ec2"
	"github.com/yuuki/albio/pkg/elbv1"
	"github.com/yuuki/albio/pkg/elbv2"
)

type DetachParam struct {
	InstanceID string
}

func Detach(param *DetachParam) error {
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

	elbClient := elbv1.New(sess)
	elbs, err := elbClient.GetLoadBalancersFromInstanceID(instanceID)
	if err != nil {
		return err
	}

	albClient := elbv2.New(sess)
	albs, err := albClient.GetLoadBalancersFromInstanceID(instanceID)
	if err != nil {
		return err
	}

	if len(elbs) < 1 && len(albs) < 1 {
		return fmt.Errorf("%v is not attached any loadbalancers", instanceID)
	}

	lbs := append(elbs, albs...)
	if err := ec2Client.SaveLoadBalancersToTag(instanceID, lbs.NameSlice()); err != nil {
		return err
	}

	if len(elbs) > 0 {
		log.Println("-->", "Detaching", instanceID, "from", elbs)
		if err := elbClient.RemoveInstanceFromLoadBalancers(instanceID, elbs); err != nil {
			return err
		}
	}
	if len(albs) > 0 {
		log.Println("-->", "Detaching", instanceID, "from", albs)
		if err := albClient.RemoveInstanceFromLoadBalancers(instanceID, albs); err != nil {
			return err
		}
	}

	return nil
}
