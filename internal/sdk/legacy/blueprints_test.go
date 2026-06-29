package morpheus_test

import (
	"testing"

	"github.com/HewlettPackard/hpe-morpheus-go-sdk/legacy"
)

func TestListBlueprints(t *testing.T) {
	client := getTestClient(t)
	req := &morpheus.Request{}
	resp, err := client.ListBlueprints(req)
	assertResponse(t, resp, err)
}

func TestGetBlueprint(t *testing.T) {
	client := getTestClient(t)
	req := &morpheus.Request{}
	resp, err := client.ListBlueprints(req)
	assertResponse(t, resp, err)

	// parse JSON and fetch the first one by ID

	result := resp.Result.(*morpheus.ListBlueprintsResult)
	recordCount := result.Meta.Total
	t.Logf("Found %d Blueprints.", recordCount)
	if recordCount != 0 {
		// Get by ID
		record := (*result.Blueprints)[0]
		resp, err = client.GetBlueprint(record.ID, &morpheus.Request{})
		assertResponse(t, resp, err)
	}
}
