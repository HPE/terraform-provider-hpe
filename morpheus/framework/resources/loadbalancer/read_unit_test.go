// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package loadbalancer

import (
	"testing"

	sdk "github.com/HPE/terraform-provider-hpe/internal/sdk/oapigen"
)

// TestMapTenantNilFields guards the raw-dereference fix for load balancer tenants: the
// API may omit a tenant's id or name, which previously panicked at *in.Id / *in.Name.
// Both must map to a Terraform null via the nil-safe convert.* helpers.
func TestMapTenantNilFields(t *testing.T) {
	t.Parallel()

	got := mapTenant(sdk.GetLoadBalancer200ResponseLoadBalancerTenantsInner{}) // nil Id and Name

	if !got.Id.IsNull() {
		t.Errorf("tenant id = %v, want null", got.Id)
	}

	if !got.Name.IsNull() {
		t.Errorf("tenant name = %v, want null", got.Name)
	}
}

// TestMapTenantPopulated confirms populated id/name still map through.
func TestMapTenantPopulated(t *testing.T) {
	t.Parallel()

	id := int64(5)
	name := "acme"
	got := mapTenant(sdk.GetLoadBalancer200ResponseLoadBalancerTenantsInner{Id: &id, Name: &name})

	if got.Id.IsNull() || got.Id.ValueInt64() != id {
		t.Errorf("tenant id = %v, want %d", got.Id, id)
	}

	if got.Name.IsNull() || got.Name.ValueString() != name {
		t.Errorf("tenant name = %v, want %q", got.Name, name)
	}
}
