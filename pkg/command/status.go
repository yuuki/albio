package command

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/yuuki/albio/pkg/ec2"
	"github.com/yuuki/albio/pkg/elb"
)

func Status() error {
	sess := session.New()
	instanceID, err := ec2.New(sess).GetLocalInstanceID()
	if err != nil {
		return err
	}

	lbs, err := elb.New(sess).GetLoadBalancersFromInstanceID(instanceID)
	if err != nil {
		return err
	}

	b, err := json.MarshalIndent(lbs, "", "    ")
	if err != nil {
		return err
	}
	fmt.Fprintf(os.Stdout, "%s\n", b)

	return nil
}
