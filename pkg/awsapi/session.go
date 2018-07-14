package awsapi

import (
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/ec2metadata"
	"github.com/aws/aws-sdk-go/aws/session"
)

// Session gets the AWS session.
func Session() (*session.Session, error) {
	region := os.Getenv("AWS_REGION")
	sess, err := session.NewSession(aws.NewConfig())
	if err != nil {
		return nil, err
	}
	if region == "" {
		region, _ = ec2metadata.New(sess).Region()
	}
	return session.NewSession(
		aws.NewConfig().WithRegion(region).WithCredentialsChainVerboseErrors(true),
	)
}
