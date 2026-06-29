package morpheus_test

import (
	"testing"

	"github.com/HewlettPackard/hpe-morpheus-go-sdk/legacy"
)

var (
	testEnvironmentName = "tfmorph-test"
)

func TestListEnvironments(t *testing.T) {
	client := getTestClient(t)
	req := &morpheus.Request{}
	resp, err := client.ListEnvironments(req)
	assertResponse(t, resp, err)
}

func TestGetEnvironment(t *testing.T) {
	client := getTestClient(t)
	req := &morpheus.Request{}
	resp, err := client.ListEnvironments(req)
	assertResponse(t, resp, err)

	// parse JSON and fetch the first one by ID

	result := resp.Result.(*morpheus.ListEnvironmentsResult)
	recordCount := result.Meta.Total
	t.Logf("Found %d Environments.", recordCount)
	if recordCount != 0 {
		// Get by ID
		record := (*result.Environments)[0]
		resp, err = client.GetEnvironment(record.ID, &morpheus.Request{})
		assertResponse(t, resp, err)
	}
}

func TestEnvironmentCRUD(t *testing.T) {
	client := getTestClient(t)
	//create
	req := &morpheus.Request{
		Body: map[string]interface{}{
			"environment": map[string]interface{}{
				"active":      true,
				"name":        testEnvironmentName,
				"code":        testEnvironmentName,
				"description": "Test Bunker",
				"visibility":  "private",
			},
		},
	}
	resp, err := client.CreateEnvironment(req)
	result := resp.Result.(*morpheus.CreateEnvironmentResult)
	assertResponse(t, resp, err)
	assertNotNil(t, result)
	assertEqual(t, result.Success, true)
	assertNotNil(t, result.Environment)
	assertNotEqual(t, result.Environment.ID, 0)
	assertEqual(t, result.Environment.Active, true)
	assertEqual(t, result.Environment.Name, testEnvironmentName)
	assertEqual(t, result.Environment.Code, testEnvironmentName)
	assertEqual(t, result.Environment.Description, "Test Bunker")
	assertEqual(t, result.Environment.Visibility, "private")

	// update
	updateReq := &morpheus.Request{
		Body: map[string]interface{}{
			"environment": map[string]interface{}{
				"description": "Test Lab",
			},
		},
	}
	updateResp, updateErr := client.UpdateEnvironment(result.Environment.ID, updateReq)
	updateResult := updateResp.Result.(*morpheus.UpdateEnvironmentResult)
	assertResponse(t, updateResp, updateErr)
	assertNotNil(t, updateResult)
	assertEqual(t, updateResult.Success, true)
	assertNotNil(t, updateResult.Environment)
	assertNotEqual(t, updateResult.Environment.ID, 0)
	assertEqual(t, updateResult.Environment.Description, "Test Lab")

	// delete
	deleteReq := &morpheus.Request{}
	deleteResp, deleteErr := client.DeleteEnvironment(result.Environment.ID, deleteReq)
	deleteResult := deleteResp.Result.(*morpheus.DeleteEnvironmentResult)
	assertResponse(t, deleteResp, deleteErr)
	assertNotNil(t, deleteResult)
	assertEqual(t, deleteResult.Success, true)
}
