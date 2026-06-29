package morpheus_test

import (
	"testing"

	"github.com/HewlettPackard/hpe-morpheus-go-sdk/legacy"
)

func TestExecuteSchedules(t *testing.T) {
	client := getTestClient(t)
	req := &morpheus.Request{}
	resp, err := client.ListExecuteSchedules(req)
	assertResponse(t, resp, err)
}

func TestGetExecutionSchedule(t *testing.T) {
	client := getTestClient(t)
	req := &morpheus.Request{}
	resp, err := client.ListExecuteSchedules(req)
	assertResponse(t, resp, err)

	// parse JSON and fetch the first one by ID

	result := resp.Result.(*morpheus.ListExecuteSchedulesResult)
	recordCount := result.Meta.Total
	t.Logf("Found %d Execute Schedules.", recordCount)
	if recordCount != 0 {
		// Get by ID
		record := (*result.ExecuteSchedules)[0]
		resp, err = client.GetExecuteSchedule(record.ID, &morpheus.Request{})
		assertResponse(t, resp, err)
	}
}
