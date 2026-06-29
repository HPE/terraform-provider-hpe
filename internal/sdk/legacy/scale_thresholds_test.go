package morpheus_test

import (
	"testing"

	"github.com/HewlettPackard/hpe-morpheus-go-sdk/legacy"
)

func TestListScaleThresholds(t *testing.T) {
	client := getTestClient(t)
	req := &morpheus.Request{}
	resp, err := client.ListScaleThresholds(req)
	assertResponse(t, resp, err)
}

func TestGetScaleThreshold(t *testing.T) {
	client := getTestClient(t)
	req := &morpheus.Request{}
	resp, err := client.ListScaleThresholds(req)
	assertResponse(t, resp, err)

	// parse JSON and fetch the first one by ID
	result := resp.Result.(*morpheus.ListScaleThresholdsResult)
	recordCount := result.Meta.Total
	t.Logf("Found %d Scale Thresholds.", recordCount)
	if recordCount != 0 {
		// Get by ID
		record := (*result.ScaleThresholds)[0]
		resp, err = client.GetScaleThreshold(record.ID, &morpheus.Request{})
		assertResponse(t, resp, err)
	}
}
