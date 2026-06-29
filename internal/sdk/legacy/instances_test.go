package morpheus_test

import (
	"testing"

	"github.com/HewlettPackard/hpe-morpheus-go-sdk/legacy"
)

func TestListInstances(t *testing.T) {
	client := getTestClient(t)
	resp, err := client.ListInstances(&morpheus.Request{})
	assertResponse(t, resp, err)
}

func TestGetInstance(t *testing.T) {
	client := getTestClient(t)
	resp, err := client.ListInstances(&morpheus.Request{})
	assertResponse(t, resp, err)
	// parse JSON and fetch the first one by ID
	listInstancesResult := resp.Result.(*morpheus.ListInstancesResult)
	instancesCount := listInstancesResult.Meta.Total
	t.Logf("Found %d Instances.", instancesCount)
	if instancesCount != 0 {
		firstInstance := (*listInstancesResult.Instances)[0]
		// log.Printf(fmt.Sprintf("First Instance: [%d] %v: ", firstInstance.ID, firstInstance.Name))
		resp, err = client.GetInstance(firstInstance.ID, &morpheus.Request{})
		assertResponse(t, resp, err)
	}

}

// this requires params zoneId&layoutId&siteId, heh
// func TestListInstancePlans(t *testing.T) {
// 	client := getTestClient(t)
// 	resp, err := client.ListInstancePlans(&morpheus.Request{})
// 	assertResponse(t, resp, err)
// }
