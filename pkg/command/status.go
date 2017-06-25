package command

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/yuuki/albio/pkg/alb"
	"github.com/yuuki/albio/pkg/ec2"
	"github.com/yuuki/albio/pkg/elb"
)

type StatusParam struct {
	InstanceID string
}

func Status(param *StatusParam) error {
	sess := session.New()

	var instanceID string
	if param.InstanceID == "" {
		var err error
		instanceID, err = ec2.New(sess).GetLocalInstanceID()
		if err != nil {
			return err
		}
	}

	elbs, err := elb.New(sess).GetLoadBalancersFromInstanceID(instanceID)
	if err != nil {
		return err
	}
	albs, err := alb.New(sess).GetLoadBalancersFromInstanceID(instanceID)
	if err != nil {
		return err
	}
	lbs := append(elbs, albs...)

	b, err := json.MarshalIndent(lbs, "", "    ")
	if err != nil {
		return err
	}
	fmt.Fprintf(os.Stdout, "%s\n", b)

	return nil
}
