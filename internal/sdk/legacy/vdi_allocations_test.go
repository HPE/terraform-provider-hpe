package morpheus_test

import (
	"testing"

	"github.com/HewlettPackard/hpe-morpheus-go-sdk/legacy"
)

func TestListVDIAllocations(t *testing.T) {
	client := getTestClient(t)
	req := &morpheus.Request{}
	resp, err := client.ListVDIAllocations(req)
	assertResponse(t, resp, err)
}

func TestGetVDIAllocation(t *testing.T) {
	client := getTestClient(t)
	req := &morpheus.Request{}
	resp, err := client.ListVDIAllocations(req)
	assertResponse(t, resp, err)

	// parse JSON and fetch the first one by ID

	result := resp.Result.(*morpheus.ListVDIAllocationsResult)
	recordCount := result.Meta.Total
	t.Logf("Found %d VDI Allocations.", recordCount)
	if recordCount != 0 {
		// Get by ID
		record := (*result.VDIAllocations)[0]
		resp, err = client.GetVDIAllocation(record.ID, &morpheus.Request{})
		assertResponse(t, resp, err)
	}
}
