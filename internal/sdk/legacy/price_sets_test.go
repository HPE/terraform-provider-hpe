package morpheus_test

import (
	"testing"

	"github.com/HewlettPackard/hpe-morpheus-go-sdk/legacy"
)

func TestPriceSet(t *testing.T) {
	client := getTestClient(t)
	req := &morpheus.Request{}
	resp, err := client.ListPriceSets(req)
	assertResponse(t, resp, err)
}

func TestGetPriceSet(t *testing.T) {
	client := getTestClient(t)
	req := &morpheus.Request{}
	resp, err := client.ListPriceSets(req)
	assertResponse(t, resp, err)

	// parse JSON and fetch the first one by ID

	result := resp.Result.(*morpheus.ListPriceSetsResult)
	recordCount := result.Meta.Total
	t.Logf("Found %d Price Sets.", recordCount)
	if recordCount != 0 {
		// Get by ID
		record := (*result.PriceSets)[0]
		resp, err = client.GetPriceSet(record.ID, &morpheus.Request{})
		assertResponse(t, resp, err)
	}
}

func TestPriceSetsCRUD(t *testing.T) {
	client := getTestClient(t)
	//create
	req := &morpheus.Request{
		Body: map[string]interface{}{
			"priceSet": map[string]interface{}{
				"name": "sdk-test",
				"code": "sdk-test",
				"account": map[string]interface{}{
					"id": 1,
				},
				"priceUnit": "month",
				"type":      "fixed",
			},
		},
	}
	resp, err := client.CreatePriceSet(req)
	result := resp.Result.(*morpheus.CreatePriceSetResult)
	assertResponse(t, resp, err)
	t.Log("Price Create Complete")
	assertNotNil(t, result)
	assertEqual(t, result.Success, true)
	getRequest := &morpheus.Request{}
	getResponse, err := client.GetPriceSet(result.ID, getRequest)
	if err != nil {
		t.Error(err)
	}
	getResult := getResponse.Result.(*morpheus.GetPriceSetResult)
	assertEqual(t, getResult.PriceSet.Type, "fixed")
	t.Log("Initial Get Complete")

	// update
	updateReq := &morpheus.Request{
		Body: map[string]interface{}{
			"priceSet": map[string]interface{}{
				"restartUsage": false,
			},
		},
	}
	updateResp, updateErr := client.UpdatePriceSet(result.ID, updateReq)
	updateResult := updateResp.Result.(*morpheus.UpdatePriceSetResult)
	getRequest = &morpheus.Request{}
	getResponse, err = client.GetPriceSet(result.ID, getRequest)
	if err != nil {
		t.Error(err)
	}
	getResult = getResponse.Result.(*morpheus.GetPriceSetResult)
	assertEqual(t, getResult.PriceSet.RestartUsage, false)
	assertResponse(t, updateResp, updateErr)
	assertNotNil(t, updateResult)
	assertEqual(t, updateResult.Success, true)
	t.Log("Price Update Complete")

	// delete
	deleteReq := &morpheus.Request{}
	deleteResp, deleteErr := client.DeletePriceSet(result.ID, deleteReq)
	deleteResult := deleteResp.Result.(*morpheus.DeletePriceSetResult)
	assertResponse(t, deleteResp, deleteErr)
	assertNotNil(t, deleteResult)
	assertEqual(t, deleteResult.Success, true)

}
