package morpheus_test

import (
	"testing"

	"github.com/HewlettPackard/hpe-morpheus-go-sdk/legacy"
)

func TestListKeyPairs(t *testing.T) {
	client := getTestClient(t)
	resp, err := client.ListKeyPairs()
	assertResponse(t, resp, err)
}

func TestGetKeyPair(t *testing.T) {
	client := getTestClient(t)
	resp, err := client.ListKeyPairs()
	assertResponse(t, resp, err)

	// parse JSON and fetch the first one by ID

	result := resp.Result.(*morpheus.ListKeyPairsResult)
	recordCount := result.Meta.Total
	t.Logf("Found %d Key Pairs.", recordCount)
	if recordCount != 0 {
		// Get by ID
		record := (*result.KeyPairs)[0]
		resp, err = client.GetKeyPair(record.ID)
		assertResponse(t, resp, err)
	}
}
