package morpheus_test

import (
	"testing"

	"github.com/HewlettPackard/hpe-morpheus-go-sdk/legacy"
)

func TestListNetworks(t *testing.T) {
	client := getTestClient(t)
	req := &morpheus.Request{}
	resp, err := client.ListNetworks(req)
	assertResponse(t, resp, err)
}

func TestGetNetwork(t *testing.T) {
	client := getTestClient(t)
	req := &morpheus.Request{}
	resp, err := client.ListNetworks(req)
	assertResponse(t, resp, err)

	// parse JSON and fetch the first one by ID

	result := resp.Result.(*morpheus.ListNetworksResult)
	recordCount := result.Meta.Total
	t.Logf("Found %d Networks.", recordCount)
	if recordCount != 0 {
		// Get by ID
		record := (*result.Networks)[0]
		resp, err = client.GetNetwork(record.ID, &morpheus.Request{})
		assertResponse(t, resp, err)
	}
}
