package morpheus_test

import (
	"testing"

	"github.com/HewlettPackard/hpe-morpheus-go-sdk/legacy"
)

func TestListVDIGateways(t *testing.T) {
	client := getTestClient(t)
	req := &morpheus.Request{}
	resp, err := client.ListVDIGateways(req)
	assertResponse(t, resp, err)
}

func TestGetVDIGateway(t *testing.T) {
	client := getTestClient(t)
	req := &morpheus.Request{}
	resp, err := client.ListVDIGateways(req)
	assertResponse(t, resp, err)

	// parse JSON and fetch the first one by ID
	result := resp.Result.(*morpheus.ListVDIGatewaysResult)
	recordCount := result.Meta.Total
	t.Logf("Found %d VDI Gateways.", recordCount)
	if recordCount != 0 {
		// Get by ID
		record := (*result.VDIGateways)[0]
		resp, err = client.GetVDIGateway(record.ID, &morpheus.Request{})
		assertResponse(t, resp, err)
	}
}
