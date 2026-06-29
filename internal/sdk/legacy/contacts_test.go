package morpheus_test

import (
	"testing"

	"github.com/HewlettPackard/hpe-morpheus-go-sdk/legacy"
)

func TestContacts(t *testing.T) {
	client := getTestClient(t)
	req := &morpheus.Request{}
	resp, err := client.ListContacts(req)
	assertResponse(t, resp, err)
}

func TestGetContact(t *testing.T) {
	client := getTestClient(t)
	req := &morpheus.Request{}
	resp, err := client.ListContacts(req)
	assertResponse(t, resp, err)

	// parse JSON and fetch the first one by ID

	result := resp.Result.(*morpheus.ListContactsResult)
	recordCount := result.Meta.Total
	t.Logf("Found %d Contacts.", recordCount)
	if recordCount != 0 {
		// Get by ID
		record := (*result.Contacts)[0]
		resp, err = client.GetContact(record.ID, &morpheus.Request{})
		assertResponse(t, resp, err)
	}
}
