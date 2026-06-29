package morpheus_test

import (
	"testing"

	"github.com/HewlettPackard/hpe-morpheus-go-sdk/legacy"
)

func TestListSpecTemplates(t *testing.T) {
	client := getTestClient(t)
	req := &morpheus.Request{}
	resp, err := client.ListSpecTemplates(req)
	assertResponse(t, resp, err)
}

func TestGetSpecTemplate(t *testing.T) {
	client := getTestClient(t)
	req := &morpheus.Request{}
	resp, err := client.ListSpecTemplates(req)
	assertResponse(t, resp, err)

	// parse JSON and fetch the first one by ID

	result := resp.Result.(*morpheus.ListSpecTemplatesResult)
	recordCount := result.Meta.Total
	t.Logf("Found %d Spec Templates.", recordCount)
	if recordCount != 0 {
		// Get by ID
		record := (*result.SpecTemplates)[0]
		resp, err = client.GetSpecTemplate(record.ID, &morpheus.Request{})
		assertResponse(t, resp, err)
	}
}
