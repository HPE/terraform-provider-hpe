package morpheus_test

import (
	"testing"

	"github.com/HewlettPackard/hpe-morpheus-go-sdk/legacy"
)

func TestListBudgets(t *testing.T) {
	client := getTestClient(t)
	req := &morpheus.Request{}
	resp, err := client.ListBudgets(req)
	assertResponse(t, resp, err)
}

func TestGetBudget(t *testing.T) {
	client := getTestClient(t)
	req := &morpheus.Request{}
	resp, err := client.ListBudgets(req)
	assertResponse(t, resp, err)

	// parse JSON and fetch the first one by ID
	result := resp.Result.(*morpheus.ListBudgetsResult)
	recordCount := result.Meta.Total
	t.Logf("Found %d Budgets.", recordCount)
	if recordCount != 0 {
		// Get by ID
		record := (*result.Budgets)[0]
		resp, err = client.GetBudget(record.ID, &morpheus.Request{})
		assertResponse(t, resp, err)
	}
}
