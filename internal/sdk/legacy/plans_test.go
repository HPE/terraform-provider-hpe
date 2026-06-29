package morpheus_test

import (
	"testing"

	"github.com/HewlettPackard/hpe-morpheus-go-sdk/legacy"
)

func TestPlans(t *testing.T) {
	client := getTestClient(t)
	req := &morpheus.Request{}
	resp, err := client.ListPlans(req)
	assertResponse(t, resp, err)
}

func TestGetPlan(t *testing.T) {
	client := getTestClient(t)
	req := &morpheus.Request{}
	resp, err := client.ListPlans(req)
	assertResponse(t, resp, err)

	// parse JSON and fetch the first one by ID

	result := resp.Result.(*morpheus.ListPlansResult)
	recordCount := result.Meta.Total
	t.Logf("Found %d Plans.", recordCount)
	if recordCount != 0 {
		// Get by ID
		record := (*result.Plans)[0]
		resp, err = client.GetPlan(record.ID, &morpheus.Request{})
		assertResponse(t, resp, err)
	}
}
