package morpheus_test

import (
	"testing"

	"github.com/HewlettPackard/hpe-morpheus-go-sdk/legacy"
)

func TestListUsers(t *testing.T) {
	client := getTestClient(t)
	req := &morpheus.Request{}
	resp, err := client.ListUsers(req)
	assertResponse(t, resp, err)
}

func TestGetUser(t *testing.T) {
	client := getTestClient(t)
	req := &morpheus.Request{}
	resp, err := client.ListUsers(req)
	assertResponse(t, resp, err)

	// parse JSON and fetch the first one by ID
	result := resp.Result.(*morpheus.ListUsersResult)
	recordCount := result.Meta.Total
	t.Logf("Found %d Users.", recordCount)
	if recordCount != 0 {
		// Get by ID
		record := (*result.Users)[0]
		resp, err = client.GetUser(record.ID, &morpheus.Request{})
		assertResponse(t, resp, err)
	}
}
