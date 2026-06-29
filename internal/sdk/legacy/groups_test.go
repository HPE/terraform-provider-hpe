package morpheus_test

import (
	"testing"

	"github.com/HewlettPackard/hpe-morpheus-go-sdk/legacy"
)

var (
	testGroupName = "tfmorph-test"
)

func TestListGroups(t *testing.T) {
	client := getTestClient(t)
	req := &morpheus.Request{}
	resp, err := client.ListGroups(req)
	assertResponse(t, resp, err)
}

func TestGetGroup(t *testing.T) {
	client := getTestClient(t)
	req := &morpheus.Request{}
	resp, err := client.ListGroups(req)
	assertResponse(t, resp, err)

	// parse JSON and fetch the first one by ID

	result := resp.Result.(*morpheus.ListGroupsResult)
	recordCount := result.Meta.Total
	t.Logf("Found %d Groups.", recordCount)
	if recordCount != 0 {
		// Get by ID
		record := (*result.Groups)[0]
		resp, err = client.GetGroup(record.ID, &morpheus.Request{})
		assertResponse(t, resp, err)
	}
}

func TestGroupsCRUD(t *testing.T) {
	client := getTestClient(t)
	//create
	req := &morpheus.Request{
		Body: map[string]interface{}{
			"group": map[string]interface{}{
				"name":     testGroupName,
				"code":     testGroupName,
				"location": "Test Bunker",
			},
		},
	}
	resp, err := client.CreateGroup(req)
	result := resp.Result.(*morpheus.CreateGroupResult)
	assertResponse(t, resp, err)
	assertNotNil(t, result)
	assertEqual(t, result.Success, true)
	assertNotNil(t, result.Group)
	assertNotEqual(t, result.Group.ID, 0)
	assertEqual(t, result.Group.Name, testGroupName)
	assertEqual(t, result.Group.Code, testGroupName)
	assertEqual(t, result.Group.Location, "Test Bunker")

	// update
	updateReq := &morpheus.Request{
		Body: map[string]interface{}{
			"group": map[string]interface{}{
				"location": "Test Lab",
			},
		},
	}
	updateResp, updateErr := client.UpdateGroup(result.Group.ID, updateReq)
	updateResult := updateResp.Result.(*morpheus.UpdateGroupResult)
	assertResponse(t, updateResp, updateErr)
	assertNotNil(t, updateResult)
	assertEqual(t, updateResult.Success, true)
	assertNotNil(t, updateResult.Group)
	assertNotEqual(t, updateResult.Group.ID, 0)
	assertEqual(t, updateResult.Group.Location, "Test Lab")

	// delete
	deleteReq := &morpheus.Request{}
	deleteResp, deleteErr := client.DeleteGroup(result.Group.ID, deleteReq)
	deleteResult := deleteResp.Result.(*morpheus.DeleteGroupResult)
	assertResponse(t, deleteResp, deleteErr)
	assertNotNil(t, deleteResult)
	assertEqual(t, deleteResult.Success, true)
}
