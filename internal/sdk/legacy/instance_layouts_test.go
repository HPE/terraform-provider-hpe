package morpheus_test

import (
	"testing"

	"github.com/HewlettPackard/hpe-morpheus-go-sdk/legacy"
)

func TestListInstanceLayouts(t *testing.T) {
	client := getTestClient(t)
	req := &morpheus.Request{}
	resp, err := client.ListInstanceLayouts(req)
	assertResponse(t, resp, err)
}

func TestGetInstanceLayout(t *testing.T) {
	client := getTestClient(t)
	req := &morpheus.Request{}
	resp, err := client.ListInstanceLayouts(req)
	assertResponse(t, resp, err)

	// parse JSON and fetch the first one by ID
	result := resp.Result.(*morpheus.ListInstanceLayoutsResult)
	recordCount := result.Meta.Total
	t.Logf("Found %d Instance Layouts.", recordCount)
	if recordCount != 0 {
		// Get by ID
		record := (*result.InstanceLayouts)[0]
		resp, err = client.GetInstanceLayout(record.ID, &morpheus.Request{})
		assertResponse(t, resp, err)
	}
}
