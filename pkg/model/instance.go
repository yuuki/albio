package model

import (
	goelb "github.com/aws/aws-sdk-go/service/elb"
	"github.com/aws/aws-sdk-go/service/elbv2"
)

// Instance represents a EC2 instance.
type Instance struct {
	InstanceID string `json:"instance-id"`
}

// NewInstance creates a Instance object.
func NewInstance(instance *goelb.Instance) *Instance {
	return &Instance{InstanceID: *instance.InstanceId}
}

// NewInstanceFromTarget creates a Instance object from elbv2.TargetDescription.
func NewInstanceFromTarget(target *elbv2.TargetDescription) *Instance {
	return &Instance{InstanceID: *target.Id}
}

// String returns a string reprentation of Instance.
func (i *Instance) String() string {
	return i.InstanceID
}
