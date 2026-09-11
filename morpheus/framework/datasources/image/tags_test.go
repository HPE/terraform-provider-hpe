// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package image

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-go/tftypes"

	sdk "github.com/HPE/terraform-provider-hpe/internal/sdk/oapigen"

	"github.com/HPE/terraform-provider-hpe/utils/convert"
)

type testTag struct {
	name  string
	value string
}

// tagElements builds the tags set exactly as the read path does (driving the
// production tagValue mapper), then lowers it to the tftypes values Terraform
// actually compares. Building an equivalent value by hand would pass whether or
// not the mapper is correct, which is the mistake this test exists to catch --
// the fault is invisible at the attr.Value layer and only appears once the value
// is lowered to tftypes, in SetType.Validate, after it has left the provider.
func tagElements(t *testing.T, tags []testTag) []tftypes.Value {
	t.Helper()

	inputs := make([]sdk.GetVirtualImage200ResponseVirtualImageTagsInner, 0, len(tags))
	for _, tg := range tags {
		inputs = append(inputs, sdk.GetVirtualImage200ResponseVirtualImageTagsInner{
			Name:  &tg.name,
			Value: &tg.value,
		})
	}

	set, diags := convert.ToSetType(context.Background(), inputs, tagValue)
	if diags.HasError() {
		t.Fatalf("building the tags set failed: %v", diags.Errors())
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

// TestUnitImageTagsAreDistinguishable is the MORPH-16289 regression guard for the
// image tags mapper -- the sibling of tenants, and the mapper the ticket cites as
// having been correct while its neighbour (tenants) was not. Two tags with
// different names/values must not lower to equal tftypes values.
func TestUnitImageTagsAreDistinguishable(t *testing.T) {
	t.Parallel()

	elements := tagElements(t, []testTag{
		{name: "env", value: "prod"},
		{name: "tier", value: "web"},
	})

	if len(elements) != 2 {
		t.Fatalf("expected 2 elements, got %d", len(elements))
	}

	if elements[0].Equal(elements[1]) {
		t.Error(
			"two tags with different names and values compare equal; " +
				"Terraform will reject the set as containing duplicates",
		)
	}
}

// TestUnitImageRepeatedTagsAreEqual is the counterpart: genuinely identical tags
// should compare equal (deduplication is not this mapper's job).
func TestUnitImageRepeatedTagsAreEqual(t *testing.T) {
	t.Parallel()

	elements := tagElements(t, []testTag{
		{name: "env", value: "prod"},
		{name: "env", value: "prod"},
	})

	if len(elements) != 2 {
		t.Fatalf("expected 2 elements, got %d", len(elements))
	}

	if !elements[0].Equal(elements[1]) {
		t.Error("genuinely identical tags should compare equal")
	}
}
