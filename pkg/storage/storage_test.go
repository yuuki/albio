package storage

import (
	"bytes"
	"strings"
	"testing"

	"github.com/kylelemons/godebug/pretty"
	"github.com/yuuki/albio/pkg/model"
)

func TestSaveLoadBalancers(t *testing.T) {
	out := new(bytes.Buffer)

	err := SaveLoadBalancers(out, "i-0f5ffb9f0a75e6b85", []*model.LoadBalancer{
		{
			Name:    "albiotestalb01",
			DNSName: "albiotestalb01-2a3abff5251f4a60.elb.ap-northeast-1.amazonaws.com",
			Arn:     "arn:aws:elasticloadbalancing:ap-northeast-1:962625566896:targetgroup/albiotestalb01/3057e052ca8ea602",
			Type:    "application",
			Scheme:  "internet-facing",
		},
		{
			Name:    "albiotestnlb01",
			DNSName: "albiotestnlb01-2a3abff5251f4a60.elb.ap-northeast-1.amazonaws.com",
			Arn:     "arn:aws:elasticloadbalancing:ap-northeast-1:962625566896:targetgroup/albiotestnlb01/3057e052ca8ea602",
			Type:    "network",
			Scheme:  "internet-facing",
		},
	})
	if err != nil {
		t.Fatalf("is should not raise error:%s", err)
	}

	want := `{"loadbalancers":[{"name":"albiotestalb01","arn":"arn:aws:elasticloadbalancing:ap-northeast-1:962625566896:targetgroup/albiotestalb01/3057e052ca8ea602","type":"application"},{"name":"albiotestnlb01","arn":"arn:aws:elasticloadbalancing:ap-northeast-1:962625566896:targetgroup/albiotestnlb01/3057e052ca8ea602","type":"network"}],"instance-id":"i-0f5ffb9f0a75e6b85"}
`
	if out.String() != want {
		t.Errorf("got: %v, want: %v", out.String(), want)
	}
}

func TestLoadLoadBalancers(t *testing.T) {
	in := strings.NewReader(`
{"loadbalancers":[{"name":"albiotestalb01","arn":"arn:aws:elasticloadbalancing:ap-northeast-1:962625566896:targetgroup/albiotestalb01/3057e052ca8ea602","type":"application"},{"name":"albiotestnlb01","arn":"arn:aws:elasticloadbalancing:ap-northeast-1:962625566896:targetgroup/albiotestnlb01/3057e052ca8ea602","type":"network"}],"instance-id":"i-0f5ffb9f0a75e6b85"}
`)

	got, err := LoadLoadBalancers(in, "i-0f5ffb9f0a75e6b85")
	if err != nil {
		t.Fatalf("%v", err)
	}

	want := []LoadBalancer{
		{
			Name: "albiotestalb01",
			Arn:  "arn:aws:elasticloadbalancing:ap-northeast-1:962625566896:targetgroup/albiotestalb01/3057e052ca8ea602",
			Type: "application",
		},
		{
			Name: "albiotestnlb01",
			Arn:  "arn:aws:elasticloadbalancing:ap-northeast-1:962625566896:targetgroup/albiotestnlb01/3057e052ca8ea602",
			Type: "network",
		},
	}
	if diff := pretty.Compare(got, want); diff != "" {
		t.Errorf("LoadLoadBalancer diff: (-got +want)\n%v", diff)
	}
}
