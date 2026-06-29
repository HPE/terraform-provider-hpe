package morpheus_test

import (
	"testing"

	"github.com/HewlettPackard/hpe-morpheus-go-sdk/legacy"
)

func TestSecurityPackage(t *testing.T) {
	client := getTestClient(t)
	req := &morpheus.Request{}
	resp, err := client.ListSecurityPackages(req)
	assertResponse(t, resp, err)
}

func TestGetSecurityPackage(t *testing.T) {
	client := getTestClient(t)
	req := &morpheus.Request{}
	resp, err := client.ListSecurityPackages(req)
	assertResponse(t, resp, err)

	// parse JSON and fetch the first one by ID

	result := resp.Result.(*morpheus.ListSecurityPackagesResult)
	recordCount := result.Meta.Total
	t.Logf("Found %d Security Packages.", recordCount)
	if recordCount != 0 {
		// Get by ID
		record := (*result.SecurityPackages)[0]
		resp, err = client.GetSecurityPackage(record.ID, &morpheus.Request{})
		assertResponse(t, resp, err)
	}
}
