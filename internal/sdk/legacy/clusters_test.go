package morpheus_test

import (
	"testing"

	"github.com/HewlettPackard/hpe-morpheus-go-sdk/legacy"
)

func TestListClusters(t *testing.T) {
	client := getTestClient(t)
	req := &morpheus.Request{}
	resp, err := client.ListClusters(req)
	assertResponse(t, resp, err)
}

func TestGetCluster(t *testing.T) {
	client := getTestClient(t)
	req := &morpheus.Request{}
	resp, err := client.ListClusters(req)
	assertResponse(t, resp, err)

	// parse JSON and fetch the first one by ID

	result := resp.Result.(*morpheus.ListClustersResult)
	recordCount := result.Meta.Total
	t.Logf("Found %d Clusters.", recordCount)
	if recordCount != 0 {
		// Get by ID
		record := (*result.Clusters)[0]
		resp, err = client.GetCluster(record.ID, &morpheus.Request{})
		assertResponse(t, resp, err)
	}
}
