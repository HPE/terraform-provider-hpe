package morpheus_test

import (
	"testing"

	"github.com/HewlettPackard/hpe-morpheus-go-sdk/legacy"
)

func TestPrices(t *testing.T) {
	client := getTestClient(t)
	req := &morpheus.Request{}
	resp, err := client.ListPrices(req)
	assertResponse(t, resp, err)
}

func TestGetPrice(t *testing.T) {
	client := getTestClient(t)
	req := &morpheus.Request{}
	resp, err := client.ListPrices(req)
	assertResponse(t, resp, err)

	// parse JSON and fetch the first one by ID

	result := resp.Result.(*morpheus.ListPricesResult)
	recordCount := result.Meta.Total
	t.Logf("Found %d Prices.", recordCount)
	if recordCount != 0 {
		// Get by ID
		record := (*result.Prices)[0]
		resp, err = client.GetPrice(record.ID, &morpheus.Request{})
		assertResponse(t, resp, err)
	}
}

func TestPricesCRUD(t *testing.T) {
	client := getTestClient(t)
	//create
	req := &morpheus.Request{
		Body: map[string]interface{}{
			"price": map[string]interface{}{
				"name": "sdk-test",
				"code": "sdk-test",
				"account": map[string]interface{}{
					"id": 1,
				},
				"priceType":    "fixed",
				"priceUnit":    "month",
				"incurCharges": "running",
				"currency":     "USD",
				"cost":         10,
			},
		},
	}
	resp, err := client.CreatePrice(req)
	result := resp.Result.(*morpheus.CreatePriceResult)
	assertResponse(t, resp, err)
	assertNotNil(t, result)
	assertEqual(t, result.Success, true)
	getRequest := &morpheus.Request{}
	getResponse, err := client.GetPrice(result.ID, getRequest)
	if err != nil {
		t.Error(err)
	}
	getResult := getResponse.Result.(*morpheus.GetPriceResult)
	assertNotEqual(t, getResult.Price.Cost, 85.00)
	assertEqual(t, getResult.Price.Cost, 10.00)

	// update
	updateReq := &morpheus.Request{
		Body: map[string]interface{}{
			"price": map[string]interface{}{
				"markupType": "fixed",
				"markup":     1.25,
			},
		},
	}
	updateResp, updateErr := client.UpdatePrice(result.ID, updateReq)
	updateResult := updateResp.Result.(*morpheus.UpdatePriceResult)
	getRequest = &morpheus.Request{}
	getResponse, err = client.GetPrice(result.ID, getRequest)
	if err != nil {
		t.Error(err)
	}
	getResult = getResponse.Result.(*morpheus.GetPriceResult)
	assertNotEqual(t, getResult.Price.Markup, 2.00)
	assertEqual(t, getResult.Price.Markup, 1.25)
	assertResponse(t, updateResp, updateErr)
	assertNotNil(t, updateResult)
	assertEqual(t, updateResult.Success, true)

	// delete
	deleteReq := &morpheus.Request{}
	deleteResp, deleteErr := client.DeletePrice(result.ID, deleteReq)
	deleteResult := deleteResp.Result.(*morpheus.DeletePriceResult)
	assertResponse(t, deleteResp, deleteErr)
	assertNotNil(t, deleteResult)
	assertEqual(t, deleteResult.Success, true)
}
