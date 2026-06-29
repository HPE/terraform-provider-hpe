package morpheus_test

import (
	"testing"

	"github.com/HewlettPackard/hpe-morpheus-go-sdk/legacy"
)

func TestClusterLayouts(t *testing.T) {
	client := getTestClient(t)
	req := &morpheus.Request{}
	resp, err := client.ListClusterLayouts(req)
	assertResponse(t, resp, err)
}

func TestGetClusterLayout(t *testing.T) {
	client := getTestClient(t)
	req := &morpheus.Request{}
	resp, err := client.ListClusterLayouts(req)
	assertResponse(t, resp, err)

	// parse JSON and fetch the first one by ID
	result := resp.Result.(*morpheus.ListClusterLayoutsResult)
	recordCount := result.Meta.Total
	t.Logf("Found %d Cluster Layouts.", recordCount)
	if recordCount != 0 {
		// Get by ID
		record := (*result.ClusterLayouts)[0]
		resp, err = client.GetClusterLayout(record.ID, &morpheus.Request{})
		assertResponse(t, resp, err)
	}
}
