package morpheus_test

import (
	"testing"

	"github.com/HewlettPackard/hpe-morpheus-go-sdk/legacy"
)

func TestServicePlan(t *testing.T) {
	client := getTestClient(t)
	req := &morpheus.Request{}
	resp, err := client.ListServicePlans(req)
	assertResponse(t, resp, err)
}

func TestGetServicePlan(t *testing.T) {
	client := getTestClient(t)
	req := &morpheus.Request{}
	resp, err := client.ListServicePlans(req)
	assertResponse(t, resp, err)

	// parse JSON and fetch the first one by ID

	result := resp.Result.(*morpheus.ListServicePlansResult)
	recordCount := result.Meta.Total
	t.Logf("Found %d Service Plans.", recordCount)
	if recordCount != 0 {
		// Get by ID
		record := (*result.ServicePlans)[0]
		resp, err = client.GetServicePlan(record.ID, &morpheus.Request{})
		assertResponse(t, resp, err)
	}
}

func TestServicePlansCRUD(t *testing.T) {
	client := getTestClient(t)
	//create
	req := &morpheus.Request{
		Body: map[string]interface{}{
			"servicePlan": map[string]interface{}{
				"name":       "sdk-test",
				"code":       "sdk-test",
				"maxStorage": 1073741824,
				"maxMemory":  1073741824,
			},
		},
	}
	resp, err := client.CreateServicePlan(req)
	result := resp.Result.(*morpheus.CreateServicePlanResult)
	assertResponse(t, resp, err)
	assertNotNil(t, result)
	assertEqual(t, result.Success, true)
	getRequest := &morpheus.Request{}
	getResponse, err := client.GetServicePlan(result.ID, getRequest)
	if err != nil {
		t.Error(err)
	}
	getResult := getResponse.Result.(*morpheus.GetServicePlanResult)
	assertEqual(t, getResult.ServicePlan.Name, "sdk-test")

	// update
	updateReq := &morpheus.Request{
		Body: map[string]interface{}{
			"servicePlan": map[string]interface{}{
				"name": "sdk-test-1",
			},
		},
	}
	updateResp, updateErr := client.UpdateServicePlan(result.ID, updateReq)
	updateResult := updateResp.Result.(*morpheus.UpdateServicePlanResult)
	getRequest = &morpheus.Request{}
	getResponse, err = client.GetServicePlan(result.ID, getRequest)
	if err != nil {
		t.Error(err)
	}
	getResult = getResponse.Result.(*morpheus.GetServicePlanResult)
	assertEqual(t, getResult.ServicePlan.Name, "sdk-test-1")
	assertResponse(t, updateResp, updateErr)
	assertNotNil(t, updateResult)
	assertEqual(t, updateResult.Success, true)

	// delete
	deleteReq := &morpheus.Request{}
	deleteResp, deleteErr := client.DeleteServicePlan(result.ID, deleteReq)
	deleteResult := deleteResp.Result.(*morpheus.DeleteServicePlanResult)
	assertResponse(t, deleteResp, deleteErr)
	assertNotNil(t, deleteResult)
	assertEqual(t, deleteResult.Success, true)

}
