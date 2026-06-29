package morpheus_test

import (
	"testing"

	"github.com/HewlettPackard/hpe-morpheus-go-sdk/legacy"
)

func TestListTasks(t *testing.T) {
	client := getTestClient(t)
	req := &morpheus.Request{}
	resp, err := client.ListTasks(req)
	assertResponse(t, resp, err)
}

func TestGetTask(t *testing.T) {
	client := getTestClient(t)
	req := &morpheus.Request{}
	resp, err := client.ListTasks(req)
	assertResponse(t, resp, err)

	// parse JSON and fetch the first one by ID

	result := resp.Result.(*morpheus.ListTasksResult)
	recordCount := result.Meta.Total
	t.Logf("Found %d Tasks.", recordCount)
	if recordCount != 0 {
		// Get by ID
		record := (*result.Tasks)[0]
		resp, err = client.GetTask(record.ID, &morpheus.Request{})
		assertResponse(t, resp, err)
	}
}
