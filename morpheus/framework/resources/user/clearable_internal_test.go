// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package user

import (
	"context"
	"strings"
	"testing"

	"github.com/HPE/terraform-provider-hpe/morpheus/utils/clearable"
)

// Every optional attribute is either clearable or excluded with a reason; a
// new one fails this test until it is classified.
func TestUnitClearableClassificationIsComplete(t *testing.T) {
	t.Parallel()

	problems := clearable.Check(UserResourceSchema(context.Background()).Attributes, clearableAttributes, notClearable)
	if len(problems) != 0 {
		t.Fatalf("clearable classification:\n  %s", strings.Join(problems, "\n  "))
	}
}

// The user keeps no record of what was configured and plans no clears
// itself: every clearable attribute carries a Nullable*UpdateModifier that
// plans the null when the configuration drops it. A clearable attribute
// without that modifier needs clearable.Record in Create and Update and
// clearable.Plan in ModifyPlan first (see the network resource).
func TestUnitClearableAttributesNeedNoRecord(t *testing.T) {
	t.Parallel()

	names := clearable.Computed(UserResourceSchema(context.Background()).Attributes, clearableAttributes)
	if len(names) != 0 {
		t.Fatalf("computed attributes declared clearable without Record/Plan wiring: %s", strings.Join(names, ", "))
	}
}
