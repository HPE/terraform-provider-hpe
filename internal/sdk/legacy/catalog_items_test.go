package morpheus_test

import (
	"testing"

	"github.com/HewlettPackard/hpe-morpheus-go-sdk/legacy"
)

func TestListCatalogItems(t *testing.T) {
	client := getTestClient(t)
	req := &morpheus.Request{}
	resp, err := client.ListCatalogItems(req)
	assertResponse(t, resp, err)
}

func TestGetCatalogItem(t *testing.T) {
	client := getTestClient(t)
	req := &morpheus.Request{}
	resp, err := client.ListCatalogItems(req)
	assertResponse(t, resp, err)

	// parse JSON and fetch the first one by ID

	result := resp.Result.(*morpheus.ListCatalogItemsResult)
	recordCount := result.Meta.Total
	t.Logf("Found %d Catalog Items.", recordCount)
	if recordCount != 0 {
		// Get by ID
		record := (*result.CatalogItems)[0]
		resp, err = client.GetCatalogItem(record.ID, &morpheus.Request{})
		assertResponse(t, resp, err)
	}
}
