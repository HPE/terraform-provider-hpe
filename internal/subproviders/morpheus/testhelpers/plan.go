// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package testhelpers

import (
	"context"
	"crypto/rand"
	"fmt"
	"net/http"
	"strings"
	"testing"
	"time"
	"github.com/HewlettPackard/hpe-morpheus-go-sdk/sdk"

	"github.com/HPE/terraform-provider-hpe/internal/subproviders/morpheus/datasources/plan"
	"github.com/HPE/terraform-provider-hpe/internal/subproviders/morpheus/errors"
)

func CreatePlan(t *testing.T) (*sdk.GetServicePlans200ResponseServicePlan, error) {
	t.Helper()

	name := fmt.Sprintf("testacc-%s-%s", t.Name(), rand.Text())
	provisionType := sdk.NewAddClusterLayoutsRequestLayoutProvisionType(44)

	addPlan := sdk.NewAddServicePlansRequestServicePlanWithDefaults()

	ctx := context.TODO()
	client := newClient(ctx, t)

	addPlan.SetName(name)
	addPlan.SetCode(strings.ToLower(name))
	addPlan.SetProvisionType(*provisionType)
	addPlan.SetMaxMemory(0)
	addPlan.SetMaxStorage(0)
	addPlanReq := sdk.NewAddServicePlansRequest(*addPlan)

	p, hresp, err := client.ServicePlansAPI.AddServicePlans(ctx).AddServicePlansRequest(
		*addPlanReq).Execute()

	if err != nil || hresp.StatusCode != http.StatusOK || p == nil {
		return nil, fmt.Errorf("POST failed for plan: %s ", errors.ErrMsg(err, hresp))
	}

	return plan.GetPlanByID(ctx, p.GetId(), client)
}

func DeletePlan(t *testing.T, id int64) error {
	t.Helper()

	ctx := context.TODO()

	client := newClient(ctx, t)

	_, hresp, err := client.ServicePlansAPI.RemoveServicePlans(ctx, id).Execute()
	if err != nil || hresp.StatusCode != http.StatusOK {
		t.Fatalf("DELETE failed for plan %d: %s", id, errors.ErrMsg(err, hresp))
	}
	for range 6 {
		_, hresp, _ := client.ServicePlansAPI.GetServicePlans(ctx, id).Execute()
		if hresp.StatusCode == http.StatusNotFound {
			return nil
		}

		time.Sleep(time.Second * 10)
	}

	return fmt.Errorf("DELETE failed for plan %d", id)
}
