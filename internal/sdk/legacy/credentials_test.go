package morpheus_test

import (
	"testing"

	"github.com/HewlettPackard/hpe-morpheus-go-sdk/legacy"
)

func TestListCredentials(t *testing.T) {
	client := getTestClient(t)
	req := &morpheus.Request{}
	resp, err := client.ListCredentials(req)
	assertResponse(t, resp, err)
}

func TestGetCredential(t *testing.T) {
	client := getTestClient(t)
	req := &morpheus.Request{}
	resp, err := client.ListCredentials(req)
	assertResponse(t, resp, err)

	// parse JSON and fetch the first one by ID
	result := resp.Result.(*morpheus.ListCredentialsResult)
	recordCount := result.Meta.Total
	t.Logf("Found %d Credentials.", recordCount)
	if recordCount != 0 {
		// Get by ID
		record := (*result.Credentials)[0]
		resp, err = client.GetCredential(record.ID, &morpheus.Request{})
		assertResponse(t, resp, err)
	}
}
