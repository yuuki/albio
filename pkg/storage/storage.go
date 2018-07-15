package storage

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"

	"github.com/yuuki/albio/pkg/model"
)

type LoadBalancer struct {
	Arn  string `json:"arn"`
	Type string `json:"type"` // "alb" or "nlb"
}

// SaveLoadBalancers saves LoadBalancers as JSON.
func SaveLoadBalancers(w io.Writer, instanceID string, loadbalancers model.LoadBalancers) error {
	var lbs []LoadBalancer
	for _, lb := range loadbalancers {
		lbs = append(lbs, LoadBalancer{
			Arn:  lb.Arn,
			Type: lb.Type,
		})
	}
	result := struct {
		LoadBalancers []LoadBalancer `json:"loadbalancers"`
		InstanceID    string         `json:"instance-id"`
	}{
		LoadBalancers: lbs,
		InstanceID:    instanceID,
	}
	b, err := json.Marshal(result)
	if err != nil {
		return err
	}
	if _, err := fmt.Fprintf(w, "%s", b); err != nil {
		return err
	}
	return nil
}

// LoadLoadBalancers load LoadBalancers as JSON.
func LoadLoadBalancers(r io.Reader) ([]LoadBalancer, error) {
	b, err := ioutil.ReadAll(r)
	if err != nil {
		return []LoadBalancer{}, err
	}
	var result struct {
		LoadBalancers []LoadBalancer `json:"loadbalancers"`
		InstanceID    string         `json:"instance-id"`
	}
	if err := json.Unmarshal(b, &result); err != nil {
		return []LoadBalancer{}, err
	}
	return result.LoadBalancers, nil
}
