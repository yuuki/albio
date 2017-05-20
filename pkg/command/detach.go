package command

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws/session"

	"github.com/yuuki/albio/pkg/ec2"
	"github.com/yuuki/albio/pkg/elb"
)

func Detach() error {
	sess := session.New()
	instanceID, err := ec2.New(sess).GetLocalInstanceID()
	if err != nil {
		return err
	}
	lbNames, err := elb.New(sess).GetLoadBalancersFromInstanceID(instanceID)
	if err != nil {
		return err
	}
	fmt.Printf("%+v", lbNames)
	return nil
}
