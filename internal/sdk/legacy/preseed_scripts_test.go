package morpheus_test

import (
	"testing"

	"github.com/HewlettPackard/hpe-morpheus-go-sdk/legacy"
)

func TestListPreseedScripts(t *testing.T) {
	client := getTestClient(t)
	req := &morpheus.Request{}
	resp, err := client.ListPreseedScripts(req)
	assertResponse(t, resp, err)
}

func TestGetPreseedScript(t *testing.T) {
	client := getTestClient(t)
	req := &morpheus.Request{}
	resp, err := client.ListPreseedScripts(req)
	assertResponse(t, resp, err)

	// parse JSON and fetch the first one by ID
	result := resp.Result.(*morpheus.ListPreseedScriptsResult)
	recordCount := result.Meta.Total
	t.Logf("Found %d Preseed Scripts.", recordCount)
	if recordCount != 0 {
		// Get by ID
		record := (*result.PreseedScripts)[0]
		resp, err = client.GetPreseedScript(record.ID, &morpheus.Request{})
		assertResponse(t, resp, err)
	}
}
