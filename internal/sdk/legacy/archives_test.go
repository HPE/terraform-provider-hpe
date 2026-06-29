package morpheus_test

import (
	"testing"

	"github.com/HewlettPackard/hpe-morpheus-go-sdk/legacy"
)

func TestListArchives(t *testing.T) {
	client := getTestClient(t)
	req := &morpheus.Request{}
	resp, err := client.ListArchives(req)
	assertResponse(t, resp, err)
}

func TestGetArchive(t *testing.T) {
	client := getTestClient(t)
	req := &morpheus.Request{}
	resp, err := client.ListArchives(req)
	assertResponse(t, resp, err)

	// parse JSON and fetch the first one by ID
	result := resp.Result.(*morpheus.ListArchivesResult)
	recordCount := result.Meta.Total
	t.Logf("Found %d Archives.", recordCount)
	if recordCount != 0 {
		// Get by ID
		record := (*result.Archives)[0]
		resp, err = client.GetArchive(record.ID, &morpheus.Request{})
		assertResponse(t, resp, err)
	}
}
