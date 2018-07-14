package command

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/yuuki/albio/pkg/awsapi"
	"github.com/yuuki/albio/pkg/ec2"
	"github.com/yuuki/albio/pkg/elbv1"
	"github.com/yuuki/albio/pkg/elbv2"
)

type StatusParam struct {
	InstanceID string
}

func Status(param *StatusParam) error {
	sess, err := awsapi.Session()
	if err != nil {
		return err
	}

	var instanceID string
	if param.InstanceID == "" {
		var err error
		instanceID, err = ec2.New(sess).GetLocalInstanceID()
		if err != nil {
			return err
		}
	}

	elbs, err := elbv1.New(sess).GetLoadBalancersFromInstanceID(instanceID)
	if err != nil {
		return err
	}

	albs, err := elbv2.New(sess).GetLoadBalancersFromInstanceID(instanceID)
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
