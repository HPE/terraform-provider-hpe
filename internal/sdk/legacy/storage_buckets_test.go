package morpheus_test

import (
	"testing"

	"github.com/HewlettPackard/hpe-morpheus-go-sdk/legacy"
)

func TestListStorageBuckets(t *testing.T) {
	client := getTestClient(t)
	req := &morpheus.Request{}
	resp, err := client.ListStorageBuckets(req)
	assertResponse(t, resp, err)
}

func TestGetStorageBucket(t *testing.T) {
	client := getTestClient(t)
	req := &morpheus.Request{}
	resp, err := client.ListStorageBuckets(req)
	assertResponse(t, resp, err)

	// parse JSON and fetch the first one by ID
	result := resp.Result.(*morpheus.ListStorageBucketsResult)
	recordCount := result.Meta.Total
	t.Logf("Found %d Storage Buckets.", recordCount)
	if recordCount != 0 {
		// Get by ID
		record := (*result.StorageBuckets)[0]
		resp, err = client.GetStorageBucket(record.ID, &morpheus.Request{})
		assertResponse(t, resp, err)
	}
}
