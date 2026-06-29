package morpheus_test

import (
	"testing"

	"github.com/HewlettPackard/hpe-morpheus-go-sdk/legacy"
)

func TestListNetworkProxies(t *testing.T) {
	client := getTestClient(t)
	req := &morpheus.Request{}
	resp, err := client.ListNetworkProxies(req)
	assertResponse(t, resp, err)
}

func TestGetNetworkProxy(t *testing.T) {
	client := getTestClient(t)
	req := &morpheus.Request{}
	resp, err := client.ListNetworkProxies(req)
	assertResponse(t, resp, err)

	// parse JSON and fetch the first one by ID

	result := resp.Result.(*morpheus.ListNetworkProxiesResult)
	recordCount := result.Meta.Total
	t.Logf("Found %d Network Proxies.", recordCount)
	if recordCount != 0 {
		// Get by ID
		record := (*result.NetworkProxies)[0]
		resp, err = client.GetNetworkProxy(record.ID, &morpheus.Request{})
		assertResponse(t, resp, err)
	}
}
