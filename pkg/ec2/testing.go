package ec2

import (
	"github.com/aws/aws-sdk-go/aws/ec2metadata"
)

type FakeEC2Metadata struct {
	EC2MetadataAPI
	FakeRegion                      func() (string, error)
	FakeGetInstanceIdentityDocument func() (ec2metadata.EC2InstanceIdentityDocument, error)
}

func (m *FakeEC2Metadata) Region() (string, error) {
	return m.FakeRegion()
}

func (m *FakeEC2Metadata) GetInstanceIdentityDocument() (ec2metadata.EC2InstanceIdentityDocument, error) {
	return m.FakeGetInstanceIdentityDocument()
}
