package morpheus_test

import (
	"testing"

	"github.com/HewlettPackard/hpe-morpheus-go-sdk/legacy"
)

func TestListInstanceTypes(t *testing.T) {
	client := getTestClient(t)
	req := &morpheus.Request{}
	resp, err := client.ListInstanceTypes(req)
	assertResponse(t, resp, err)
}

func TestGetInstanceType(t *testing.T) {
	client := getTestClient(t)
	req := &morpheus.Request{}
	resp, err := client.ListInstanceTypes(req)
	assertResponse(t, resp, err)

	// parse JSON and fetch the first one by ID
	result := resp.Result.(*morpheus.ListInstanceTypesResult)
	recordCount := result.Meta.Total
	t.Logf("Found %d Instance Types.", recordCount)
	if recordCount != 0 {
		// Get by ID
		record := (*result.InstanceTypes)[0]
		resp, err = client.GetInstanceType(record.ID, &morpheus.Request{})
		assertResponse(t, resp, err)
	}
}
