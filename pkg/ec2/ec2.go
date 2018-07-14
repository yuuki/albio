package ec2

import (
	"os"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/defaults"
	"github.com/aws/aws-sdk-go/aws/ec2metadata"
	"github.com/aws/aws-sdk-go/aws/session"
	goec2 "github.com/aws/aws-sdk-go/service/ec2"
	"github.com/yuuki/albio/pkg/awsapi"
)

const LOAD_BALANCERS_TAG_KEY = "albio-loadbalancers"

type EC2 interface {
	GetLocalInstanceID() (string, error)
	GetLoadBalancerNamesFromTag(string) ([]string, error)
	SaveLoadBalancersToTag(string, []string) error
}

type _ec2 struct {
	svc         awsapi.EC2API
	metadataSvc awsapi.EC2MetadataAPI
}

func New(sess *session.Session) EC2 {
	mdSvc := ec2metadata.New(sess)
	region := os.Getenv("AWS_REGION")
	if region == "" {
		region, _ = mdSvc.Region()
	}
	return &_ec2{
		svc: goec2.New(sess,
			defaults.Config().WithRegion(region).WithCredentialsChainVerboseErrors(true),
		),
		metadataSvc: mdSvc,
	}
}

// GetLocalInstanceID gets an ID of local EC2 instance from EC2 metadata.
func (e *_ec2) GetLocalInstanceID() (string, error) {
	doc, err := e.metadataSvc.GetInstanceIdentityDocument()
	if err != nil {
		return "", err
	}
	return doc.InstanceID, nil
}

// GetLoadBalancersFromTag gets the information of loadbalancers that
// saves by SaveLoadBalancersToTag.
func (e *_ec2) GetLoadBalancerNamesFromTag(instanceID string) ([]string, error) {
	resp, err := e.svc.DescribeTags(&goec2.DescribeTagsInput{
		Filters: []*goec2.Filter{
			&goec2.Filter{Name: aws.String("resource-id"), Values: []*string{aws.String(instanceID)}},
			&goec2.Filter{Name: aws.String("key"), Values: []*string{aws.String(LOAD_BALANCERS_TAG_KEY)}},
		},
		MaxResults: aws.Int64(5), // minimu value is defined as 5.
	})
	if err != nil {
		return nil, err
	}
	if len(resp.Tags) < 1 {
		return []string{}, nil
	}

	var lbNames []string
	for _, name := range strings.Split(*resp.Tags[0].Value, ",") {
		lbNames = append(lbNames, name)
	}
	return lbNames, nil
}

// SaveLoadBalancersToUserdata saves the information of loadbalancers
// to rejoin into EC2 instance tag.
// ex. Key: albio-loadbalancers, Value: elb01, elb02
func (e *_ec2) SaveLoadBalancersToTag(instanceID string, lbNames []string) error {
	_, err := e.svc.CreateTags(&goec2.CreateTagsInput{
		Resources: []*string{aws.String(instanceID)},
		Tags: []*goec2.Tag{&goec2.Tag{
			Key:   aws.String(LOAD_BALANCERS_TAG_KEY),
			Value: aws.String(strings.Join(lbNames, ",")),
		}},
	})
	return err
}
