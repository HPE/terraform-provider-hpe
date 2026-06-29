package morpheus_test

import (
	"testing"

	"github.com/HewlettPackard/hpe-morpheus-go-sdk/legacy"
)

func TestCheckApps(t *testing.T) {
	client := getTestClient(t)
	req := &morpheus.Request{}
	resp, err := client.ListCheckApps(req)
	assertResponse(t, resp, err)
}

func TestGetCheckApp(t *testing.T) {
	client := getTestClient(t)
	req := &morpheus.Request{}
	resp, err := client.ListCheckApps(req)
	assertResponse(t, resp, err)

	// parse JSON and fetch the first one by ID
	result := resp.Result.(*morpheus.ListCheckAppsResult)
	recordCount := result.Meta.Total
	t.Logf("Found %d Check Apps.", recordCount)
	if recordCount != 0 {
		// Get by ID
		record := (*result.CheckApps)[0]
		resp, err = client.GetCheckApp(record.ID, &morpheus.Request{})
		assertResponse(t, resp, err)
	}
}
