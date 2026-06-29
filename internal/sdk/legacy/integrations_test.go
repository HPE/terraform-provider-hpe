package morpheus_test

import (
	"testing"

	"github.com/HewlettPackard/hpe-morpheus-go-sdk/legacy"
)

func TestIntegrations(t *testing.T) {
	client := getTestClient(t)
	req := &morpheus.Request{}
	resp, err := client.ListIntegrations(req)
	assertResponse(t, resp, err)
}

func TestGetIntegration(t *testing.T) {
	client := getTestClient(t)
	req := &morpheus.Request{}
	resp, err := client.ListIntegrations(req)
	assertResponse(t, resp, err)

	// parse JSON and fetch the first one by ID

	result := resp.Result.(*morpheus.ListIntegrationsResult)
	recordCount := result.Meta.Total
	t.Logf("Found %d integrations.", recordCount)
	if recordCount != 0 {
		// Get by ID
		record := (*result.Integrations)[0]
		resp, err = client.GetIntegration(record.ID, &morpheus.Request{})
		assertResponse(t, resp, err)
	}
}
