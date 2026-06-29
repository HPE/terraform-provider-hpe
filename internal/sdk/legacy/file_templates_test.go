package morpheus_test

import (
	"testing"

	"github.com/HewlettPackard/hpe-morpheus-go-sdk/legacy"
)

func TestFileTemplates(t *testing.T) {
	client := getTestClient(t)
	req := &morpheus.Request{}
	resp, err := client.ListFileTemplates(req)
	assertResponse(t, resp, err)
}

func TestGetFileTemplate(t *testing.T) {
	client := getTestClient(t)
	req := &morpheus.Request{}
	resp, err := client.ListFileTemplates(req)
	assertResponse(t, resp, err)

	// parse JSON and fetch the first one by ID

	result := resp.Result.(*morpheus.ListFileTemplatesResult)
	recordCount := result.Meta.Total
	t.Logf("Found %d File Templates.", recordCount)
	if recordCount != 0 {
		// Get by ID
		record := (*result.FileTemplates)[0]
		resp, err = client.GetFileTemplate(record.ID, &morpheus.Request{})
		assertResponse(t, resp, err)
	}
}
