package morpheus_test

import (
	"testing"

	"github.com/HewlettPackard/hpe-morpheus-go-sdk/legacy"
)

var (
	testWikiName     = "tfmorph-test-wiki"
	testWikiCategory = "sdktest"
)

func TestWikis(t *testing.T) {
	client := getTestClient(t)
	req := &morpheus.Request{}
	resp, err := client.ListWikis(req)
	assertResponse(t, resp, err)
}

func TestGetWiki(t *testing.T) {
	client := getTestClient(t)
	req := &morpheus.Request{}
	resp, err := client.ListWikis(req)
	assertResponse(t, resp, err)

	// parse JSON and fetch the first one by ID

	result := resp.Result.(*morpheus.ListWikisResult)
	recordCount := result.Meta.Total
	t.Logf("Found %d Wikis.", recordCount)
	if recordCount != 0 {
		// Get by ID
		record := (*result.Wikis)[0]
		resp, err = client.GetWiki(record.ID, &morpheus.Request{})
		assertResponse(t, resp, err)
	}
}

func TestWikisCRUD(t *testing.T) {
	client := getTestClient(t)
	//create
	req := &morpheus.Request{
		Body: map[string]interface{}{
			"page": map[string]interface{}{
				"name":     testWikiName,
				"category": testWikiCategory,
				"content":  "Test Bunker",
			},
		},
	}
	resp, err := client.CreateWiki(req)
	result := resp.Result.(*morpheus.CreateWikiResult)
	assertResponse(t, resp, err)
	assertNotNil(t, result)
	assertEqual(t, result.Success, true)
	assertNotNil(t, result.Wiki)
	assertNotEqual(t, result.Wiki.ID, 0)
	assertEqual(t, result.Wiki.Name, testWikiName)
	assertEqual(t, result.Wiki.Category, testWikiCategory)
	assertEqual(t, result.Wiki.Content, "Test Bunker")

	// update
	updateReq := &morpheus.Request{
		Body: map[string]interface{}{
			"page": map[string]interface{}{
				"content": "Test Lab",
			},
		},
	}
	updateResp, updateErr := client.UpdateWiki(result.Wiki.ID, updateReq)
	updateResult := updateResp.Result.(*morpheus.UpdateWikiResult)
	assertResponse(t, updateResp, updateErr)
	assertNotNil(t, updateResult)
	assertEqual(t, updateResult.Success, true)
	assertNotNil(t, updateResult.Wiki)
	assertNotEqual(t, updateResult.Wiki.ID, 0)
	assertEqual(t, updateResult.Wiki.Content, "Test Lab")

	// delete
	deleteReq := &morpheus.Request{}
	deleteResp, deleteErr := client.DeleteWiki(result.Wiki.ID, deleteReq)
	deleteResult := deleteResp.Result.(*morpheus.DeleteWikiResult)
	assertResponse(t, deleteResp, deleteErr)
	assertNotNil(t, deleteResult)
	assertEqual(t, deleteResult.Success, true)
}
