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

	"github.com/HPE/terraform-provider-hpe/internal/subproviders/morpheus/datasources/serviceplan"
	"github.com/HPE/terraform-provider-hpe/internal/subproviders/morpheus/errors"
)

func CreateServicePlan(t *testing.T) (*sdk.GetServicePlans200ResponseServicePlan, error) {
	t.Helper()

	name := fmt.Sprintf("testacc-%s-%s", t.Name(), rand.Text())
	provisionType := sdk.NewAddClusterLayoutsRequestLayoutProvisionType(44)

	addServicePlan := sdk.NewAddServicePlansRequestServicePlanWithDefaults()

	ctx := context.TODO()
	client := newClient(ctx, t)

	addServicePlan.SetName(name)
	addServicePlan.SetCode(strings.ToLower(name))
	addServicePlan.SetProvisionType(*provisionType)
	addServicePlan.SetMaxMemory(0)
	addServicePlan.SetMaxStorage(0)
	addServicePlanReq := sdk.NewAddServicePlansRequest(*addServicePlan)

	sp, hresp, err := client.ServicePlansAPI.AddServicePlans(ctx).AddServicePlansRequest(
		*addServicePlanReq).Execute()

	if err != nil || hresp.StatusCode != http.StatusOK || sp == nil {
		return nil, fmt.Errorf("POST failed for service_plan: %s ", errors.ErrMsg(err, hresp))
	}

	return serviceplan.GetServicePlanByID(ctx, sp.GetId(), client)
}

func DeleteServicePlan(t *testing.T, id int64) error {
	t.Helper()

	ctx := context.TODO()

	client := newClient(ctx, t)

	_, hresp, err := client.ServicePlansAPI.RemoveServicePlans(ctx, id).Execute()
	if err != nil || hresp.StatusCode != http.StatusOK {
		t.Fatalf("DELETE failed for service_plan %d: %s", id, errors.ErrMsg(err, hresp))
	}
	for range 6 {
		_, hresp, _ := client.ServicePlansAPI.GetServicePlans(ctx, id).Execute()
		if hresp.StatusCode == http.StatusNotFound {
			return nil
		}

		time.Sleep(time.Second * 10)
	}

	return fmt.Errorf("DELETE failed for service_plan %d", id)
}
