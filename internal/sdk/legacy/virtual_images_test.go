package morpheus_test

import (
	"testing"

	"github.com/HewlettPackard/hpe-morpheus-go-sdk/legacy"
)

func TestVirtualImages(t *testing.T) {
	client := getTestClient(t)
	req := &morpheus.Request{}
	resp, err := client.ListVirtualImages(req)
	assertResponse(t, resp, err)
}

func TestGetVirtualImage(t *testing.T) {
	client := getTestClient(t)
	req := &morpheus.Request{
		QueryParams: map[string]string{"filterType": "All"},
	}
	resp, err := client.ListVirtualImages(req)
	assertResponse(t, resp, err)

	// parse JSON and fetch the first one by ID
	result := resp.Result.(*morpheus.ListVirtualImagesResult)
	recordCount := result.Meta.Total
	t.Logf("Found %d Virtual Images.", recordCount)
	if recordCount != 0 {
		// Get by ID
		record := (*result.VirtualImages)[0]
		resp, err = client.GetVirtualImage(record.ID, &morpheus.Request{})
		assertResponse(t, resp, err)
	}
}
