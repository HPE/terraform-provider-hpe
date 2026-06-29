package morpheus_test

import (
	"testing"

	"github.com/HewlettPackard/hpe-morpheus-go-sdk/legacy"
)

func TestPowerSchedules(t *testing.T) {
	client := getTestClient(t)
	req := &morpheus.Request{}
	resp, err := client.ListPowerSchedules(req)
	assertResponse(t, resp, err)
}

func TestGetPowerSchedule(t *testing.T) {
	client := getTestClient(t)
	req := &morpheus.Request{}
	resp, err := client.ListPowerSchedules(req)
	assertResponse(t, resp, err)

	// parse JSON and fetch the first one by ID

	result := resp.Result.(*morpheus.ListPowerSchedulesResult)
	recordCount := result.Meta.Total
	t.Logf("Found %d Power Schedules.", recordCount)
	if recordCount != 0 {
		// Get by ID
		record := (*result.PowerSchedules)[0]
		resp, err = client.GetPowerSchedule(record.ID, &morpheus.Request{})
		assertResponse(t, resp, err)
	}
}

func TestPowerSchedulesCRUD(t *testing.T) {
	client := getTestClient(t)
	//create
	req := &morpheus.Request{
		Body: map[string]interface{}{
			"schedule": map[string]interface{}{
				"name":             "business hours only",
				"description":      "testing",
				"enabled":          true,
				"scheduleType":     "power",
				"scheduleTimezone": "UTC",
				"sundayOn":         0,
				"sundayOff":        0,
				"mondayOn":         7,
				"mondayOff":        15,
				"tuesdayOn":        7,
				"tuesdayOff":       15,
				"wednesdayOn":      7,
				"wednesdayOff":     15,
				"thursdayOn":       7,
				"thursdayOff":      15,
				"fridayOn":         7,
				"fridayOff":        15,
				"saturdayOn":       0,
				"saturdayOff":      0,
			},
		},
	}
	resp, err := client.CreatePowerSchedule(req)
	result := resp.Result.(*morpheus.CreatePowerScheduleResult)
	assertResponse(t, resp, err)
	assertNotNil(t, result)
	assertEqual(t, result.Success, true)
	assertNotNil(t, result.PowerSchedule)
	assertNotEqual(t, result.PowerSchedule.ID, 0)
	assertEqual(t, result.PowerSchedule.Name, "business hours only")
	assertEqual(t, result.PowerSchedule.Description, "testing")

	// update
	updateReq := &morpheus.Request{
		Body: map[string]interface{}{
			"schedule": map[string]interface{}{
				"fridayOff": 15,
			},
		},
	}
	updateResp, updateErr := client.UpdatePowerSchedule(result.PowerSchedule.ID, updateReq)
	updateResult := updateResp.Result.(*morpheus.UpdatePowerScheduleResult)
	assertResponse(t, updateResp, updateErr)
	assertNotNil(t, updateResult)
	assertEqual(t, updateResult.Success, true)
	assertNotNil(t, updateResult.PowerSchedule)
	assertNotEqual(t, updateResult.PowerSchedule.ID, 0)
	assertEqual(t, updateResult.PowerSchedule.FridayOff, 15.0)

	// delete
	deleteReq := &morpheus.Request{}
	deleteResp, deleteErr := client.DeletePowerSchedule(result.PowerSchedule.ID, deleteReq)
	deleteResult := deleteResp.Result.(*morpheus.DeletePowerScheduleResult)
	assertResponse(t, deleteResp, deleteErr)
	assertNotNil(t, deleteResult)
	assertEqual(t, deleteResult.Success, true)
}
