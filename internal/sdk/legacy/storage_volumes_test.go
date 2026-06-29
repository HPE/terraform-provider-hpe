package morpheus_test

import (
	"testing"

	"github.com/HewlettPackard/hpe-morpheus-go-sdk/legacy"
)

func TestListStorageVolumes(t *testing.T) {
	client := getTestClient(t)
	req := &morpheus.Request{}
	resp, err := client.ListStorageVolumes(req)
	assertResponse(t, resp, err)
}

func TestGetStorageVolume(t *testing.T) {
	client := getTestClient(t)
	req := &morpheus.Request{}
	resp, err := client.ListStorageVolumes(req)
	assertResponse(t, resp, err)

	// parse JSON and fetch the first one by ID
	result := resp.Result.(*morpheus.ListStorageVolumesResult)
	recordCount := result.Meta.Total
	t.Logf("Found %d storage volumes.", recordCount)
	if recordCount != 0 {
		// Get by ID
		record := (*result.StorageVolumes)[0]
		resp, err = client.GetStorageVolume(record.ID.(int64), &morpheus.Request{})
		assertResponse(t, resp, err)
	}
}
