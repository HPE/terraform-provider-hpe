package morpheus_test

import (
	"testing"

	"github.com/HewlettPackard/hpe-morpheus-go-sdk/legacy"
)

func TestListApps(t *testing.T) {
	client := getTestClient(t)
	req := &morpheus.Request{}
	resp, err := client.ListApps(req)
	assertResponse(t, resp, err)
}

func TestGetApp(t *testing.T) {
	client := getTestClient(t)
	req := &morpheus.Request{}
	resp, err := client.ListApps(req)
	assertResponse(t, resp, err)

	// parse JSON and fetch the first one by ID
	result := resp.Result.(*morpheus.ListAppsResult)
	recordCount := result.Meta.Total
	t.Logf("Found %d Apps.", recordCount)
	if recordCount != 0 {
		// Get by ID
		record := (*result.Apps)[0]
		resp, err = client.GetApp(record.ID, &morpheus.Request{})
		assertResponse(t, resp, err)
	}
}
