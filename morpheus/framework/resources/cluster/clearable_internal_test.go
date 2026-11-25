// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package cluster

import (
	"context"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/diag"

	"github.com/HPE/terraform-provider-hpe/utils/convert"

	"github.com/HPE/terraform-provider-hpe/morpheus/utils/clearable"
)

// Every optional attribute is either clearable or excluded with a reason; a
// new one fails this test until it is classified.
func TestUnitClearableClassificationIsComplete(t *testing.T) {
	t.Parallel()

	problems := clearable.Check(ClusterResourceSchema(context.Background()).Attributes, clearableAttributes, notClearable)
	if len(problems) != 0 {
		t.Fatalf("clearable classification:\n  %s", strings.Join(problems, "\n  "))
	}
}

// The resource keeps no record of what was configured and plans no clears
// itself: every clearable attribute is Optional without Computed, so the
// plan carries the null and Update sends it. Declaring an Optional+Computed
// attribute clearable needs clearable.Record in Create and Update and
// clearable.Plan in ModifyPlan first (see the network resource); this test
// fails until that wiring exists.
func TestUnitClearableAttributesNeedNoRecord(t *testing.T) {
	t.Parallel()

	names := clearable.Computed(ClusterResourceSchema(context.Background()).Attributes, clearableAttributes)
	if len(names) != 0 {
		t.Fatalf("computed attributes declared clearable without Record/Plan wiring: %s", strings.Join(names, ", "))
	}
}

// Removing host_profile_id from the configuration detaches the profile: the
// attribute is Optional without Computed, so the plan carries the null and
// Update calls the detach endpoint. Making it Computed (with or without
// UseStateForUnknown) would keep the prior value in the plan and turn a
// removal into a silent no-op, so this test fails first.
func TestUnitHostProfileIDRemovalPlansNull(t *testing.T) {
	t.Parallel()

	attr, ok := ClusterResourceSchema(context.Background()).Attributes["host_profile_id"]
	if !ok {
		t.Fatal("host_profile_id is missing from the cluster schema")
	}

	if !attr.IsOptional() || attr.IsComputed() {
		t.Fatalf("host_profile_id: optional=%t computed=%t; want optional without computed",
			attr.IsOptional(), attr.IsComputed())
	}
}

// TestUnitWarnCredentialsInConfig pins that a config carried in state with a
// credential key produces a warning naming the key, and a clean one none.
func TestUnitWarnCredentialsInConfig(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	id := int64(5)

	dirty, err := convert.MapToDynamic(ctx, map[string]any{"apiUrl": "https://k8s:6443", "apiToken": "eyJ…"})
	if err != nil {
		t.Fatal(err)
	}

	var diags diag.Diagnostics

	warnCredentialsInConfig(ctx, &id, dirty, &diags)

	if diags.WarningsCount() != 1 || diags.ErrorsCount() != 0 {
		t.Fatalf("dirty config: got %d warnings, %d errors", diags.WarningsCount(), diags.ErrorsCount())
	}

	if d := diags.Warnings()[0].Detail(); !strings.Contains(d, "cluster 5") || !strings.Contains(d, "apiToken") {
		t.Errorf("warning detail = %q", d)
	}

	clean, err := convert.MapToDynamic(ctx, map[string]any{"apiUrl": "https://k8s:6443"})
	if err != nil {
		t.Fatal(err)
	}

	var none diag.Diagnostics

	warnCredentialsInConfig(ctx, &id, clean, &none)

	if none.WarningsCount() != 0 {
		t.Errorf("clean config: got %d warnings", none.WarningsCount())
	}
}
