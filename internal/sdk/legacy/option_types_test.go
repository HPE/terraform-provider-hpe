package morpheus_test

import (
	"testing"

	"github.com/HewlettPackard/hpe-morpheus-go-sdk/legacy"
)

func TestListOptionTypes(t *testing.T) {
	client := getTestClient(t)
	req := &morpheus.Request{}
	resp, err := client.ListOptionTypes(req)
	assertResponse(t, resp, err)
}

func TestGetOptionType(t *testing.T) {
	client := getTestClient(t)
	req := &morpheus.Request{}
	resp, err := client.ListOptionTypes(req)
	assertResponse(t, resp, err)

	// parse JSON and fetch the first one by ID

	result := resp.Result.(*morpheus.ListOptionTypesResult)
	recordCount := result.Meta.Total
	t.Logf("Found %d Option Types.", recordCount)
	if recordCount != 0 {
		// Get by ID
		record := (*result.OptionTypes)[0]
		resp, err = client.GetOptionType(record.ID, &morpheus.Request{})
		assertResponse(t, resp, err)
	}
}
