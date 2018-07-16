package command

import (
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/yuuki/albio/pkg/awsapi"
	"github.com/yuuki/albio/pkg/ec2"
	"github.com/yuuki/albio/pkg/elbv2"
	"github.com/yuuki/albio/pkg/storage"
)

type DetachParam struct {
	InstanceID string
	Self       bool
}

func Detach(param *DetachParam) error {
	if param.InstanceID == "" && !param.Self {
		return errors.New("require either --instance-id or --self option.")
	}

	sess, err := awsapi.Session()
	if err != nil {
		return err
	}

	ec2Client := ec2.New(sess)

	var instanceID string
	if param.Self {
		var err error
		instanceID, err = ec2Client.GetLocalInstanceID()
		if err != nil {
			return err
		}
	}

	albClient := elbv2.New(sess)
	lbs, err := albClient.GetLoadBalancersFromInstanceID(instanceID)
	if err != nil {
		return err
	}

	if len(lbs) < 1 {
		return fmt.Errorf("%v is not attached any loadbalancers", instanceID)
	}

	if err := storage.SaveLoadBalancers(os.Stdout, instanceID, lbs); err != nil {
		return err
	}

	log.Println("-->", "Detaching", instanceID, "from", lbs)
	if err := albClient.RemoveInstanceFromLoadBalancers(instanceID, lbs); err != nil {
		return err
	}

	return nil
}
