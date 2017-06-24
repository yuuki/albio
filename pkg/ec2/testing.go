package ec2

import (
	"github.com/aws/aws-sdk-go/aws/ec2metadata"
	goec2 "github.com/aws/aws-sdk-go/service/ec2"
)

type FakeEC2API struct {
	EC2API
	FakeDescribeTags func(input *goec2.DescribeTagsInput) (*goec2.DescribeTagsOutput, error)
	FakeCreateTags   func(input *goec2.CreateTagsInput) (*goec2.CreateTagsOutput, error)
}

func (e *FakeEC2API) DescribeTags(input *goec2.DescribeTagsInput) (*goec2.DescribeTagsOutput, error) {
	return e.FakeDescribeTags(input)
}

func (e *FakeEC2API) CreateTags(input *goec2.CreateTagsInput) (*goec2.CreateTagsOutput, error) {
	return e.FakeCreateTags(input)
}

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
