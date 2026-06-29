package morpheus_test

import (
	"testing"

	"github.com/HewlettPackard/hpe-morpheus-go-sdk/legacy"
)

func TestListLoadBalancers(t *testing.T) {
	client := getTestClient(t)
	req := &morpheus.Request{}
	resp, err := client.ListLoadBalancers(req)
	assertResponse(t, resp, err)
}

func TestGetLoadBalancer(t *testing.T) {
	client := getTestClient(t)
	req := &morpheus.Request{}
	resp, err := client.ListLoadBalancers(req)
	assertResponse(t, resp, err)

	// parse JSON and fetch the first one by ID

	result := resp.Result.(*morpheus.ListLoadBalancersResult)
	recordCount := result.Meta.Total
	t.Logf("Found %d Load Balancers.", recordCount)
	if recordCount != 0 {
		// Get by ID
		record := (*result.LoadBalancers)[0]
		resp, err = client.GetLoadBalancer(record.ID, &morpheus.Request{})
		assertResponse(t, resp, err)
	}
}
