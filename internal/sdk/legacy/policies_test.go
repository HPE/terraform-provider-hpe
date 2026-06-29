package morpheus_test

import (
	"testing"

	"github.com/HewlettPackard/hpe-morpheus-go-sdk/legacy"
)

func TestPolicies(t *testing.T) {
	client := getTestClient(t)
	req := &morpheus.Request{}
	resp, err := client.ListPolicies(req)
	assertResponse(t, resp, err)
}

func TestGetPolicy(t *testing.T) {
	client := getTestClient(t)
	req := &morpheus.Request{}
	resp, err := client.ListPolicies(req)
	assertResponse(t, resp, err)

	// parse JSON and fetch the first one by ID

	result := resp.Result.(*morpheus.ListPoliciesResult)
	recordCount := result.Meta.Total
	t.Logf("Found %d Policies.", recordCount)
	if recordCount != 0 {
		// Get by ID
		record := (*result.Policies)[0]
		resp, err = client.GetPolicy(record.ID, &morpheus.Request{})
		assertResponse(t, resp, err)
	}
}
