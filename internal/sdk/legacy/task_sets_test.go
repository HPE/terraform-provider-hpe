package morpheus_test

import (
	"testing"

	"github.com/HewlettPackard/hpe-morpheus-go-sdk/legacy"
)

func TestListTaskSets(t *testing.T) {
	client := getTestClient(t)
	req := &morpheus.Request{}
	resp, err := client.ListTaskSets(req)
	assertResponse(t, resp, err)
}

func TestGetTaskSet(t *testing.T) {
	client := getTestClient(t)
	req := &morpheus.Request{}
	resp, err := client.ListTaskSets(req)
	assertResponse(t, resp, err)

	// parse JSON and fetch the first one by ID

	result := resp.Result.(*morpheus.ListTaskSetsResult)
	recordCount := result.Meta.Total
	t.Logf("Found %d Task Sets.", recordCount)
	if recordCount != 0 {
		// Get by ID
		record := (*result.TaskSets)[0]
		resp, err = client.GetTaskSet(record.ID, &morpheus.Request{})
		assertResponse(t, resp, err)

		// List by name

	}
}
