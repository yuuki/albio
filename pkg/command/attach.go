package command

import (
	"fmt"
	"log"
	"os"

	"github.com/yuuki/albio/pkg/awsapi"
	"github.com/yuuki/albio/pkg/ec2"
	"github.com/yuuki/albio/pkg/elbv2"
	"github.com/yuuki/albio/pkg/storage"
)

type AttachParam struct {
	InstanceID       string
	Self             bool
	LoadBalancerName string
}

func Attach(param *AttachParam) error {
	if param.InstanceID == "" && !param.Self {
		log.Println("require either --instance-id or --self option.")
	}

	sess, err := awsapi.Session()
	if err != nil {
		return err
	}

	ec2Client := ec2.New(sess)
	elbv2Client := elbv2.New(sess)

	var instanceID string
	if param.Self {
		var err error
		instanceID, err = ec2Client.GetLocalInstanceID()
		if err != nil {
			return err
		}

		lbs, err := storage.LoadLoadBalancers(os.Stdin, instanceID)
		if err != nil {
			return err
		}

		var lbNames []string
		for _, lb := range lbs {
			lbNames = append(lbNames, lb.Name)
		}
		if len(lbNames) < 1 {
			return fmt.Errorf("not found loadbalancers with %s. specify loadbalancer with --loadlalancer option", instanceID)
		}

		albs, err := elbv2Client.GetLoadBalancersByNames(lbNames)
		if err != nil {
			return err
		}
		log.Println("-->", "Attaching", instanceID, "to", albs)
		if err := elbv2Client.AddInstanceToLoadBalancers(instanceID, albs); err != nil {
			return err
		}
	}

	if param.LoadBalancerName != "" {
		albs, err := elbv2Client.GetLoadBalancersByNames([]string{param.LoadBalancerName})
		if err != nil {
			return err
		}
		log.Println("-->", "Attaching", instanceID, "to", albs)
		if err := elbv2Client.AddInstanceToLoadBalancers(instanceID, albs); err != nil {
			return err
		}
	}

	newLBs, err := elbv2Client.GetLoadBalancersFromInstanceID(instanceID)
	if err := storage.SaveLoadBalancers(os.Stdout, instanceID, newLBs); err != nil {
		return err
	}

	return nil
}
