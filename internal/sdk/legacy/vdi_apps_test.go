package morpheus_test

import (
	"testing"

	"github.com/HewlettPackard/hpe-morpheus-go-sdk/legacy"
)

func TestListVDIApps(t *testing.T) {
	client := getTestClient(t)
	req := &morpheus.Request{}
	resp, err := client.ListVDIApps(req)
	assertResponse(t, resp, err)
}

func TestGetVDIApp(t *testing.T) {
	client := getTestClient(t)
	req := &morpheus.Request{}
	resp, err := client.ListVDIApps(req)
	assertResponse(t, resp, err)

	// parse JSON and fetch the first one by ID

	result := resp.Result.(*morpheus.ListVDIAppsResult)
	recordCount := result.Meta.Total
	t.Logf("Found %d VDI Apps.", recordCount)
	if recordCount != 0 {
		// Get by ID
		record := (*result.VDIApps)[0]
		resp, err = client.GetVDIApp(record.ID, &morpheus.Request{})
		assertResponse(t, resp, err)
	}
}
