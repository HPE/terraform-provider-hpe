package morpheus_test

import (
	"testing"

	"github.com/HewlettPackard/hpe-morpheus-go-sdk/legacy"
)

func TestCheckGroups(t *testing.T) {
	client := getTestClient(t)
	req := &morpheus.Request{}
	resp, err := client.ListCheckGroups(req)
	assertResponse(t, resp, err)
}

func TestGetCheckGroup(t *testing.T) {
	client := getTestClient(t)
	req := &morpheus.Request{}
	resp, err := client.ListCheckGroups(req)
	assertResponse(t, resp, err)

	// parse JSON and fetch the first one by ID

	result := resp.Result.(*morpheus.ListCheckGroupsResult)
	recordCount := result.Meta.Total
	t.Logf("Found %d Check Groups.", recordCount)
	if recordCount != 0 {
		// Get by ID
		record := (*result.CheckGroups)[0]
		resp, err = client.GetCheckGroup(record.ID, &morpheus.Request{})
		assertResponse(t, resp, err)
	}
}
