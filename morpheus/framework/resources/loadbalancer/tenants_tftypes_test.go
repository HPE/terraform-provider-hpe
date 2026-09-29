// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package loadbalancer

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-go/tftypes"

	sdk "github.com/HPE/terraform-provider-hpe/internal/sdk/oapigen"

	"github.com/HPE/terraform-provider-hpe/utils/convert"
)

// tenantSetElements builds the load balancer tenants set exactly as the read
// path does (driving the production mapTenant mapper), then lowers it to the
// tftypes values Terraform actually compares in SetType.Validate. Constructing
// the value by hand would pass regardless of whether mapTenant is correct; the
// fault this guards (a *Value built without state) is invisible at the
// attr.Value layer and only surfaces once the value is lowered to tftypes.
func tenantSetElements(
	t *testing.T,
	in []sdk.GetLoadBalancer200ResponseLoadBalancerTenantsInner,
) []tftypes.Value {
	t.Helper()

	set, diags := convert.ToSetType(context.Background(), in, mapTenant)
	if diags.HasError() {
		t.Fatalf("building the tenants set failed: %v", diags.Errors())
	}

	raw, err := set.ToTerraformValue(context.Background())
	if err != nil {
		t.Fatalf("lowering to tftypes failed: %v", err)
	}

	var elements []tftypes.Value
	if err := raw.As(&elements); err != nil {
		t.Fatalf("reading elements failed: %v", err)
	}

	return elements
}

// TestUnitLoadBalancerTenantsAreDistinguishable is the MORPH-16289 regression
// guard for the load balancer tenants mapper. Two tenants with different ids and
// names must not lower to equal tftypes values, or Terraform rejects the set with
// "Duplicate Set Element".
func TestUnitLoadBalancerTenantsAreDistinguishable(t *testing.T) {
	t.Parallel()

	id1, id2 := int64(1), int64(2)
	n1, n2 := "acme", "globex"

	elements := tenantSetElements(t, []sdk.GetLoadBalancer200ResponseLoadBalancerTenantsInner{
		{Id: &id1, Name: &n1},
		{Id: &id2, Name: &n2},
	})

	if len(elements) != 2 {
		t.Fatalf("expected 2 elements, got %d", len(elements))
	}

	if elements[0].Equal(elements[1]) {
		t.Error(
			"two tenants with different ids and names compare equal; " +
				"Terraform will reject the set as containing duplicates",
		)
	}
}

// TestUnitLoadBalancerRepeatedTenantsAreEqual is the counterpart: genuinely
// identical tenants should compare equal.
func TestUnitLoadBalancerRepeatedTenantsAreEqual(t *testing.T) {
	t.Parallel()

	id := int64(1)
	name := "acme"

	elements := tenantSetElements(t, []sdk.GetLoadBalancer200ResponseLoadBalancerTenantsInner{
		{Id: &id, Name: &name},
		{Id: &id, Name: &name},
	})

	if len(elements) != 2 {
		t.Fatalf("expected 2 elements, got %d", len(elements))
	}

	if !elements[0].Equal(elements[1]) {
		t.Error("genuinely identical tenants should compare equal")
	}
}
