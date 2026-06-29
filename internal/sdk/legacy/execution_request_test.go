package morpheus_test

import (
	"strconv"
	"testing"

	"github.com/HewlettPackard/hpe-morpheus-go-sdk/legacy"
)

// not clear how to test the "get" function since it requires a uniqueID and there is no list method
// func TestGetExecutionRequest(t *testing.T) {

// }

func TestCreateExecutionRequest(t *testing.T) {
	client := getTestClient(t)
	resp, err := client.ListInstances(&morpheus.Request{})
	assertResponse(t, resp, err)

	// parse JSON and fetch the first one by ID
	listInstancesResult := resp.Result.(*morpheus.ListInstancesResult)
	instancesCount := listInstancesResult.Meta.Total
	t.Logf("Found %d Instances.", instancesCount)

	if instancesCount != 0 {
		firstInstance := (*listInstancesResult.Instances)[0]
		resp, err := client.CreateExecutionRequest(&morpheus.Request{
			QueryParams: map[string]string{
				"instanceId": strconv.Itoa(int(firstInstance.ID)),
			},
			Body: map[string]interface{}{
				"script": "pwd",
			},
		})
		assertResponse(t, resp, err)
	}
}
