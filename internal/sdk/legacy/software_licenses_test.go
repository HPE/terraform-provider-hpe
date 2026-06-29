package morpheus_test

import (
	"testing"

	"github.com/HewlettPackard/hpe-morpheus-go-sdk/legacy"
)

func TestListSoftwareLicenses(t *testing.T) {
	client := getTestClient(t)
	req := &morpheus.Request{}
	resp, err := client.ListSoftwareLicenses(req)
	assertResponse(t, resp, err)
}

func TestGetSoftwareLicense(t *testing.T) {
	client := getTestClient(t)
	req := &morpheus.Request{}
	resp, err := client.ListSoftwareLicenses(req)
	assertResponse(t, resp, err)

	// parse JSON and fetch the first one by ID
	result := resp.Result.(*morpheus.ListSoftwareLicensesResult)
	recordCount := result.Meta.Total
	t.Logf("Found %d Software Licenses.", recordCount)
	if recordCount != 0 {
		// Get by ID
		record := (*result.SoftwareLicenses)[0]
		resp, err = client.GetSoftwareLicense(record.ID, &morpheus.Request{})
		assertResponse(t, resp, err)
	}
}
