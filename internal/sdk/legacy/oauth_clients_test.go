package morpheus_test

import (
	"testing"

	"github.com/HewlettPackard/hpe-morpheus-go-sdk/legacy"
)

func TestListOauthClients(t *testing.T) {
	client := getTestClient(t)
	req := &morpheus.Request{}
	resp, err := client.ListOauthClients(req)
	assertResponse(t, resp, err)
}

func TestGetOauthClient(t *testing.T) {
	client := getTestClient(t)
	req := &morpheus.Request{}
	resp, err := client.ListOauthClients(req)
	assertResponse(t, resp, err)

	// parse JSON and fetch the first one by ID
	result := resp.Result.(*morpheus.ListOauthClientsResult)
	recordCount := result.Meta.Total
	t.Logf("Found %d Oauth Clients.", recordCount)
	if recordCount != 0 {
		// Get by ID
		record := (*result.OauthClients)[0]
		resp, err = client.GetOauthClient(record.ID, &morpheus.Request{})
		assertResponse(t, resp, err)
	}
}
