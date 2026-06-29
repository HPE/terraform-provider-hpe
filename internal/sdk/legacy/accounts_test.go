package morpheus_test

import (
	"testing"

	"github.com/HewlettPackard/hpe-morpheus-go-sdk/legacy"
)

func TestListTenants(t *testing.T) {
	client := getTestClient(t)
	req := &morpheus.Request{}
	resp, err := client.ListTenants(req)
	assertResponse(t, resp, err)
}

func TestGetTenant(t *testing.T) {
	client := getTestClient(t)
	req := &morpheus.Request{}
	resp, err := client.ListTenants(req)
	assertResponse(t, resp, err)

	// parse JSON and fetch the first one by ID
	result := resp.Result.(*morpheus.ListTenantsResult)
	recordCount := result.Meta.Total
	t.Logf("Found %d Tenants.", recordCount)
	if recordCount != 0 {
		// Get by ID
		record := (*result.Accounts)[0]
		resp, err = client.GetTenant(record.ID, &morpheus.Request{})
		assertResponse(t, resp, err)
	}
}
