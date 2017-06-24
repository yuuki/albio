package awsapi

import (
	"github.com/aws/aws-sdk-go/aws/ec2metadata"
	goec2 "github.com/aws/aws-sdk-go/service/ec2"
)

// EC2API defines the interface for EC2 API stubbing.
type EC2API interface {
	DescribeTags(input *goec2.DescribeTagsInput) (*goec2.DescribeTagsOutput, error)
	CreateTags(input *goec2.CreateTagsInput) (*goec2.CreateTagsOutput, error)
}

// EC2MetadataAPI defines the interface for EC2Metadata API stubbing.
type EC2MetadataAPI interface {
	Region() (string, error)
	GetInstanceIdentityDocument() (ec2metadata.EC2InstanceIdentityDocument, error)
}
