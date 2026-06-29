package morpheus_test

import (
	"testing"

	"github.com/HewlettPackard/hpe-morpheus-go-sdk/legacy"
)

func TestIncidents(t *testing.T) {
	client := getTestClient(t)
	req := &morpheus.Request{}
	resp, err := client.ListIncidents(req)
	assertResponse(t, resp, err)
}

func TestGetIncident(t *testing.T) {
	client := getTestClient(t)
	req := &morpheus.Request{}
	resp, err := client.ListIncidents(req)
	assertResponse(t, resp, err)

	// parse JSON and fetch the first one by ID

	result := resp.Result.(*morpheus.ListIncidentsResult)
	recordCount := result.Meta.Total
	t.Logf("Found %d Incidents.", recordCount)
	if recordCount != 0 {
		// Get by ID
		record := (*result.Incidents)[0]
		resp, err = client.GetIncident(record.ID, &morpheus.Request{})
		assertResponse(t, resp, err)
	}
}
