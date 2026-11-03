// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package ostypeimage

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

	problems := clearable.Check(OsTypeImageResourceSchema(context.Background()).Attributes, clearableAttributes, notClearable)
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

	names := clearable.Computed(OsTypeImageResourceSchema(context.Background()).Attributes, clearableAttributes)
	if len(names) != 0 {
		t.Fatalf("computed attributes declared clearable without Record/Plan wiring: %s", strings.Join(names, ", "))
	}
}
