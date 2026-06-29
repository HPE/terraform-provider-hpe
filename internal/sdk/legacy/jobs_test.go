package morpheus_test

import (
	"testing"

	"github.com/HewlettPackard/hpe-morpheus-go-sdk/legacy"
)

func TestJobs(t *testing.T) {
	client := getTestClient(t)
	req := &morpheus.Request{}
	resp, err := client.ListJobs(req)
	assertResponse(t, resp, err)
}

func TestGetJob(t *testing.T) {
	client := getTestClient(t)
	req := &morpheus.Request{}
	resp, err := client.ListJobs(req)
	assertResponse(t, resp, err)

	// parse JSON and fetch the first one by ID

	result := resp.Result.(*morpheus.ListJobsResult)
	recordCount := result.Meta.Total
	t.Logf("Found %d jobs.", recordCount)
	if recordCount != 0 {
		// Get by ID
		record := (*result.Jobs)[0]
		resp, err = client.GetJob(record.ID, &morpheus.Request{})
		assertResponse(t, resp, err)
	}
}
