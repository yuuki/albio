package ec2

import (
	"github.com/aws/aws-sdk-go/aws/ec2metadata"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/yuuki/albio/pkg/awsapi"
)

type EC2 interface {
	GetLocalInstanceID() (string, error)
}

type _ec2 struct {
	metadataSvc awsapi.EC2MetadataAPI
}

func New(sess *session.Session) EC2 {
	return &_ec2{
		metadataSvc: ec2metadata.New(sess),
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
