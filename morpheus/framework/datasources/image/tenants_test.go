package image

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-go/tftypes"

	sdk "github.com/HPE/terraform-provider-hpe/internal/sdk/oapigen"

	"github.com/HPE/terraform-provider-hpe/utils/convert"
)

type testAccount struct {
	id   int64
	name string
}

// tenantElements builds the tenants set exactly as the read path does, then
// lowers it to the tftypes values Terraform actually compares.
//
// The comparison that matters happens in SetType.Validate, after the value has
// left the provider, so a test at the attr.Value layer cannot see the failure
// this guards against.
func tenantElements(t *testing.T, accounts []testAccount) []tftypes.Value {
	t.Helper()

	// Drives the production mapper. Building an equivalent value here would
	// pass whether or not the read path is correct, which is exactly the
	// mistake this test exists to catch.
	inputs := make([]sdk.GetVirtualImage200ResponseVirtualImageAccountsInner, 0, len(accounts))
	for _, a := range accounts {
		inputs = append(inputs, sdk.GetVirtualImage200ResponseVirtualImageAccountsInner{
			Id:   &a.id,
			Name: &a.name,
		})
	}

	set, diags := convert.ToSetType(context.Background(), inputs, tenantValue)
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

// TestTenantsAreDistinguishable is the regression guard for MORPH-16245.
//
// TenantsValue is built as a struct literal, and a zero ValueState is
// ValueStateNull. A null object lowers to a null tftypes value whatever
// attributes were set on it, so every tenant became identical to every other.
// SetType.Validate compares elements pairwise and rejected the set with
// "Duplicate Set Element" -- meaning any image with two or more tenants failed
// to read, whether or not anything was duplicated.
func TestUnitImageTenantsAreDistinguishable(t *testing.T) {
	t.Parallel()

	elements := tenantElements(t, []testAccount{
		{id: 1, name: "acme"},
		{id: 2, name: "globex"},
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

// TestRepeatedTenantsAreEqual is the counterpart. Deduplication is not this
// data source's job -- if the API genuinely repeats a tenant, the elements
// should compare equal, and that is a question for the API rather than
// something to paper over here.
func TestUnitImageRepeatedTenantsAreEqual(t *testing.T) {
	t.Parallel()

	elements := tenantElements(t, []testAccount{
		{id: 1, name: "acme"},
		{id: 1, name: "acme"},
	})

	if !elements[0].Equal(elements[1]) {
		t.Error("genuinely identical tenants should compare equal")
	}
}
