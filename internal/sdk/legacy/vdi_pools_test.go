package morpheus_test

import (
	"testing"

	"github.com/HewlettPackard/hpe-morpheus-go-sdk/legacy"
)

func TestListVDIPools(t *testing.T) {
	client := getTestClient(t)
	req := &morpheus.Request{}
	resp, err := client.ListVDIPools(req)
	assertResponse(t, resp, err)
}

func TestGetVDIPool(t *testing.T) {
	client := getTestClient(t)
	req := &morpheus.Request{}
	resp, err := client.ListVDIPools(req)
	assertResponse(t, resp, err)

	// parse JSON and fetch the first one by ID
	result := resp.Result.(*morpheus.ListVDIPoolsResult)
	recordCount := result.Meta.Total
	t.Logf("Found %d VDI Pools.", recordCount)
	if recordCount != 0 {
		// Get by ID
		record := (*result.VDIPools)[0]
		resp, err = client.GetVDIPool(record.ID, &morpheus.Request{})
		assertResponse(t, resp, err)
	}
}
