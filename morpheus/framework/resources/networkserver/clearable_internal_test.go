// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package networkserver

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

	problems := clearable.Check(
		NetworkServerResourceSchema(context.Background()).Attributes,
		clearableAttributes,
		notClearable,
	)
	if len(problems) != 0 {
		t.Fatalf("clearable classification:\n  %s", strings.Join(problems, "\n  "))
	}
}

// Every clearable attribute is Optional-only, so the plan carries the null and
// Update sends it: no clearable.Record/Plan wiring is needed. This test fails
// if a clearable attribute becomes Optional+Computed without that wiring.
func TestUnitClearableAttributesNeedNoRecord(t *testing.T) {
	t.Parallel()

	names := clearable.Computed(
		NetworkServerResourceSchema(context.Background()).Attributes,
		clearableAttributes,
	)
	if len(names) != 0 {
		t.Fatalf("computed attributes declared clearable without Record/Plan wiring: %s", strings.Join(names, ", "))
	}
}
