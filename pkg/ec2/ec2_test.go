package ec2

import (
	"testing"

	"github.com/aws/aws-sdk-go/aws/ec2metadata"
)

func TestGetLocalInstanceID(t *testing.T) {
	obj := &_ec2{
		metadataSvc: &FakeEC2Metadata{
			FakeRegion: func() (string, error) {
				return "ap-northeast-1", nil
			},
			FakeGetInstanceIdentityDocument: func() (ec2metadata.EC2InstanceIdentityDocument, error) {
				return ec2metadata.EC2InstanceIdentityDocument{InstanceID: "instanceid"}, nil
			},
		},
	}
	got, err := obj.GetLocalInstanceID()
	if err != nil {
		t.Errorf("should not raise error: %v", err)
	}
	if got != "instanceid" {
		t.Errorf("GetLocalInstanceID() = %v; want \"instanceid\"", got)
	}
}
