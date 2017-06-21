package model

import (
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	goelb "github.com/aws/aws-sdk-go/service/elb"
	"github.com/kylelemons/godebug/pretty"
)

func TestNewInstance(t *testing.T) {
	got := NewInstance(&goelb.Instance{InstanceId: aws.String("instanceid")})
	expected := Instance{InstanceID: "instanceid"}
	if diff := pretty.Compare(got, expected); diff != "" {
		t.Errorf("diff: (-actual +expected)\n%s", diff)
	}
}

func TestInstanceString(t *testing.T) {
	i := NewInstance(&goelb.Instance{InstanceId: aws.String("instanceid")})
	if got := i.String(); got != "instanceid" {
		t.Errorf("model.Instance.String() = '%v'; want 'instanceid'\n", got)
	}
}
