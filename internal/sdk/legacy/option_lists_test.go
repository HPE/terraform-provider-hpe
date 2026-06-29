package morpheus_test

import (
	"testing"

	"github.com/HewlettPackard/hpe-morpheus-go-sdk/legacy"
)

func TestListOptionLists(t *testing.T) {
	client := getTestClient(t)
	req := &morpheus.Request{}
	resp, err := client.ListOptionLists(req)
	assertResponse(t, resp, err)
}

func TestGetOptionList(t *testing.T) {
	client := getTestClient(t)
	req := &morpheus.Request{}
	resp, err := client.ListOptionLists(req)
	assertResponse(t, resp, err)

	// parse JSON and fetch the first one by ID

	result := resp.Result.(*morpheus.ListOptionListsResult)
	recordCount := result.Meta.Total
	t.Logf("Found %d Option Lists.", recordCount)
	if recordCount != 0 {
		// Get by ID
		record := (*result.OptionLists)[0]
		resp, err = client.GetOptionList(record.ID, &morpheus.Request{})
		assertResponse(t, resp, err)
	}
}
