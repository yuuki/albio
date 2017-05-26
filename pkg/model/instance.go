package model

import (
	goelb "github.com/aws/aws-sdk-go/service/elb"
)

// Instance represents a EC2 instance.
type Instance struct {
	InstanceID string `json:"instance-id"`
}

// NewInstance creates a Instance object.
func NewInstance(instance *goelb.Instance) *Instance {
	return &Instance{InstanceID: *instance.InstanceId}
}

// String returns a string reprentation of Instance.
func (i *Instance) String() string {
	return i.InstanceID
}
