package morpheus_test

import (
	"testing"

	"github.com/HewlettPackard/hpe-morpheus-go-sdk/legacy"
)

func TestListUserGroups(t *testing.T) {
	client := getTestClient(t)
	req := &morpheus.Request{}
	resp, err := client.ListUserGroups(req)
	assertResponse(t, resp, err)
}

func TestGetUserGroup(t *testing.T) {
	client := getTestClient(t)
	req := &morpheus.Request{}
	resp, err := client.ListUserGroups(req)
	assertResponse(t, resp, err)

	// parse JSON and fetch the first one by ID
	result := resp.Result.(*morpheus.ListUserGroupsResult)
	recordCount := result.Meta.Total
	t.Logf("Found %d User Groups.", recordCount)
	if recordCount != 0 {
		// Get by ID
		record := (*result.UserGroups)[0]
		resp, err = client.GetUserGroup(record.ID, &morpheus.Request{})
		assertResponse(t, resp, err)
	}
}
