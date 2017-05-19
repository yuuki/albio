package command

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws/session"

	"github.com/yuuki/albio/pkg/ec2"
)

func Detach() error {
	sess := session.New()
	instanceID, err := ec2.New(sess).GetLocalInstanceID()
	if err != nil {
		return err
	}
	fmt.Println(instanceID)
	return nil
}
