package morpheus_test

import (
	"testing"

	"github.com/HewlettPackard/hpe-morpheus-go-sdk/legacy"
)

func TestListBootScripts(t *testing.T) {
	client := getTestClient(t)
	req := &morpheus.Request{}
	resp, err := client.ListBootScripts(req)
	assertResponse(t, resp, err)
}

func TestGetBootScript(t *testing.T) {
	client := getTestClient(t)
	req := &morpheus.Request{}
	resp, err := client.ListBootScripts(req)
	assertResponse(t, resp, err)

	// parse JSON and fetch the first one by ID
	result := resp.Result.(*morpheus.ListBootScriptsResult)
	recordCount := result.Meta.Total
	t.Logf("Found %d Boot Scripts.", recordCount)
	if recordCount != 0 {
		// Get by ID
		record := (*result.BootScripts)[0]
		resp, err = client.GetBootScript(record.ID, &morpheus.Request{})
		assertResponse(t, resp, err)
	}
}
