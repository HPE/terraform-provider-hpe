package morpheus_test

import (
	"testing"

	"github.com/HewlettPackard/hpe-morpheus-go-sdk/legacy"
)

func TestScriptTemplates(t *testing.T) {
	client := getTestClient(t)
	req := &morpheus.Request{}
	resp, err := client.ListScriptTemplates(req)
	assertResponse(t, resp, err)
}

func TestGetScriptTemplate(t *testing.T) {
	client := getTestClient(t)
	req := &morpheus.Request{}
	resp, err := client.ListScriptTemplates(req)
	assertResponse(t, resp, err)

	// parse JSON and fetch the first one by ID
	result := resp.Result.(*morpheus.ListScriptTemplatesResult)
	recordCount := result.Meta.Total
	t.Logf("Found %d Script Templates.", recordCount)
	if recordCount != 0 {
		// Get by ID
		record := (*result.ScriptTemplates)[0]
		resp, err = client.GetScriptTemplate(record.ID, &morpheus.Request{})
		assertResponse(t, resp, err)
	}
}
