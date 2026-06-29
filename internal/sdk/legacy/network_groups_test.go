package morpheus_test

import (
	"testing"

	"github.com/HewlettPackard/hpe-morpheus-go-sdk/legacy"
)

func TestListNetworkGroups(t *testing.T) {
	client := getTestClient(t)
	req := &morpheus.Request{}
	resp, err := client.ListNetworkGroups(req)
	assertResponse(t, resp, err)
}

func TestGetNetworkGroup(t *testing.T) {
	client := getTestClient(t)
	req := &morpheus.Request{}
	resp, err := client.ListNetworkGroups(req)
	assertResponse(t, resp, err)

	// parse JSON and fetch the first one by ID
	result := resp.Result.(*morpheus.ListNetworkGroupsResult)
	recordCount := result.Meta.Total
	t.Logf("Found %d Network Groups.", recordCount)
	if recordCount != 0 {
		// Get by ID
		record := (*result.NetworkGroups)[0]
		resp, err = client.GetNetworkGroup(record.ID, &morpheus.Request{})
		assertResponse(t, resp, err)
	}
}
