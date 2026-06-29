package morpheus_test

import (
	"testing"

	"github.com/HewlettPackard/hpe-morpheus-go-sdk/legacy"
)

func TestListRoles(t *testing.T) {
	client := getTestClient(t)
	req := &morpheus.Request{}
	resp, err := client.ListRoles(req)
	assertResponse(t, resp, err)
}

func TestGetRole(t *testing.T) {
	client := getTestClient(t)
	req := &morpheus.Request{}
	resp, err := client.ListRoles(req)
	assertResponse(t, resp, err)

	// parse JSON and fetch the first one by ID
	result := resp.Result.(*morpheus.ListRolesResult)
	recordCount := result.Meta.Total
	t.Logf("Found %d Roles.", recordCount)
	if recordCount != 0 {
		// Get by ID
		record := (*result.Roles)[0]
		resp, err = client.GetRole(record.ID, &morpheus.Request{})
		assertResponse(t, resp, err)
	}
}
