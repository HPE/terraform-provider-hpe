package morpheus_test

import (
	"testing"

	"github.com/HewlettPackard/hpe-morpheus-go-sdk/legacy"
)

func TestNodeTypes(t *testing.T) {
	client := getTestClient(t)
	req := &morpheus.Request{}
	resp, err := client.ListNodeTypes(req)
	assertResponse(t, resp, err)
}

func TestGetNodeType(t *testing.T) {
	client := getTestClient(t)
	req := &morpheus.Request{}
	resp, err := client.ListNodeTypes(req)
	assertResponse(t, resp, err)

	// parse JSON and fetch the first one by ID

	result := resp.Result.(*morpheus.ListNodeTypesResult)
	recordCount := result.Meta.Total
	t.Logf("Found %d Nodes Types.", recordCount)
	if recordCount != 0 {
		// Get by ID
		record := (*result.NodeTypes)[0]
		resp, err = client.GetNodeType(record.ID, &morpheus.Request{})
		assertResponse(t, resp, err)
	}
}
