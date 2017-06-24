package ec2

import (
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/ec2metadata"
	goec2 "github.com/aws/aws-sdk-go/service/ec2"
	"github.com/kylelemons/godebug/pretty"
	"github.com/yuuki/albio/pkg/awsapi"
)

func TestGetLocalInstanceID(t *testing.T) {
	obj := &_ec2{
		metadataSvc: &awsapi.FakeEC2Metadata{
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

func TestGetLoadBalancerNamesFromTag(t *testing.T) {
	fakeSvc := &awsapi.FakeEC2API{
		FakeDescribeTags: func(input *goec2.DescribeTagsInput) (*goec2.DescribeTagsOutput, error) {
			return &goec2.DescribeTagsOutput{
				Tags: []*goec2.TagDescription{
					{
						Key:   aws.String("albio-loadbalancers"),
						Value: aws.String("albiolb01,albiolb02"),
					},
				},
			}, nil
		},
	}
	obj := &_ec2{svc: fakeSvc}
	got, err := obj.GetLoadBalancerNamesFromTag("instanceid")
	if err != nil {
		t.Errorf("should not raise error: %v", err)
	}
	expected := []string{"albiolb01", "albiolb02"}
	if diff := pretty.Compare(got, expected); diff != "" {
		t.Errorf("GetLoadBalancerNamesFromTag(\"instanceid\") = %v; want %v", got, expected)
	}
}

func TestSaveLoadBalancersToTag(t *testing.T) {
	fakeSvc := &awsapi.FakeEC2API{
		FakeCreateTags: func(input *goec2.CreateTagsInput) (*goec2.CreateTagsOutput, error) {
			expectedResources := []*string{aws.String("instanceid")}
			if diff := pretty.Compare(input.Resources, expectedResources); diff != "" {
				t.Errorf("diff: (-actual +expected)\n%s", diff)
			}
			expectedTags := []*goec2.Tag{
				{
					Key:   aws.String("albio-loadbalancers"),
					Value: aws.String("albiolb01,albiolb02"),
				},
			}
			if diff := pretty.Compare(input.Tags, expectedTags); diff != "" {
				t.Errorf("diff: (-actual +expected)\n%s", diff)
			}
			return &goec2.CreateTagsOutput{}, nil
		},
	}
	obj := &_ec2{svc: fakeSvc}
	err := obj.SaveLoadBalancersToTag("instanceid", []string{"albiolb01", "albiolb02"})
	if err != nil {
		t.Errorf("should not raise error: %v", err)
	}
}
