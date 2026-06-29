package morpheus_test

import (
	"testing"

	"github.com/HewlettPackard/hpe-morpheus-go-sdk/legacy"
)

func TestDeployments(t *testing.T) {
	client := getTestClient(t)
	req := &morpheus.Request{}
	resp, err := client.ListDeployments(req)
	assertResponse(t, resp, err)
}

func TestGetDeployment(t *testing.T) {
	client := getTestClient(t)
	req := &morpheus.Request{}
	resp, err := client.ListDeployments(req)
	assertResponse(t, resp, err)

	// parse JSON and fetch the first one by ID
	result := resp.Result.(*morpheus.ListDeploymentsResult)
	recordCount := result.Meta.Total
	t.Logf("Found %d Deployments.", recordCount)
	if recordCount != 0 {
		// Get by ID
		record := (*result.Deployments)[0]
		resp, err = client.GetDeployment(record.ID, &morpheus.Request{})
		assertResponse(t, resp, err)
	}
}
