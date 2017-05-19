package ec2

import (
	"github.com/aws/aws-sdk-go/aws/ec2metadata"
	"github.com/aws/aws-sdk-go/aws/session"
	goec2 "github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/ec2/ec2iface"
)

type EC2 interface {
	GetLocalInstanceID() (string, error)
}

type _ec2 struct {
	svc         ec2iface.EC2API
	metadataSvc *ec2metadata.EC2Metadata // TODO interface
}

func New(sess *session.Session) EC2 {
	return &_ec2{
		svc:         goec2.New(sess),
		metadataSvc: ec2metadata.New(sess),
	}
}

func (e *_ec2) GetLocalInstanceID() (string, error) {
	doc, err := e.metadataSvc.GetInstanceIdentityDocument()
	if err != nil {
		return "", err
	}
	return doc.InstanceID, nil
}
