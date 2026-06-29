package morpheus_test

import (
	"github.com/HewlettPackard/hpe-morpheus-go-sdk/legacy"
	"testing"
)

func TestListClouds(t *testing.T) {
	client := getTestClient(t)
	resp, err := client.ListClouds(&morpheus.Request{})
	assertResponse(t, resp, err)
}

func TestGetCloud(t *testing.T) {
	client := getTestClient(t)
	resp, err := client.ListClouds(&morpheus.Request{})
	assertResponse(t, resp, err)
	// parse JSON and fetch the first one by ID
	listCloudsResult := resp.Result.(*morpheus.ListCloudsResult)
	cloudsCount := listCloudsResult.Meta.Total
	t.Logf("Found %d Clouds.", cloudsCount)
	if cloudsCount != 0 {
		firstCloud := (*listCloudsResult.Clouds)[0]
		// log.Printf(fmt.Sprintf("First Cloud: [%d] %v: ", firstCloud.ID, firstCloud.Name))
		// log.Printf("resp.Result: ", resp.Result)
		resp, err = client.GetCloud(firstCloud.ID, &morpheus.Request{})
		assertResponse(t, resp, err)
	}

}
