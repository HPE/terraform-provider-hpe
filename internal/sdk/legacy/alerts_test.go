package morpheus_test

import (
	"testing"

	"github.com/HewlettPackard/hpe-morpheus-go-sdk/legacy"
)

func TestListAlerts(t *testing.T) {
	client := getTestClient(t)
	req := &morpheus.Request{}
	resp, err := client.ListAlerts(req)
	assertResponse(t, resp, err)
}

func TestGetAlert(t *testing.T) {
	client := getTestClient(t)
	req := &morpheus.Request{}
	resp, err := client.ListAlerts(req)
	assertResponse(t, resp, err)

	// parse JSON and fetch the first one by ID
	result := resp.Result.(*morpheus.ListAlertsResult)
	recordCount := result.Meta.Total
	t.Logf("Found %d Alerts.", recordCount)
	if recordCount != 0 {
		// Get by ID
		record := (*result.Alerts)[0]
		resp, err = client.GetAlert(record.ID, &morpheus.Request{})
		assertResponse(t, resp, err)
	}
}
