package morpheus_test

import (
	"testing"

	"github.com/HewlettPackard/hpe-morpheus-go-sdk/legacy"
)

func TestListSecurityScans(t *testing.T) {
	client := getTestClient(t)
	req := &morpheus.Request{}
	resp, err := client.ListSecurityScans(req)
	assertResponse(t, resp, err)
}

func TestGetSecurityScan(t *testing.T) {
	client := getTestClient(t)
	req := &morpheus.Request{}
	resp, err := client.ListSecurityScans(req)
	assertResponse(t, resp, err)

	// parse JSON and fetch the first one by ID
	result := resp.Result.(*morpheus.ListSecurityScansResult)
	recordCount := result.Meta.Total
	t.Logf("Found %d Security scans.", recordCount)
	if recordCount != 0 {
		// Get by ID
		record := (*result.SecurityScans)[0]
		resp, err = client.GetSecurityScan(record.ID, &morpheus.Request{})
		assertResponse(t, resp, err)
	}
}
