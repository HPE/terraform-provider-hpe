// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package networkfloatingip

import (
	"context"
	"strings"
	"testing"

	"github.com/HPE/terraform-provider-hpe/morpheus/utils/clearable"
)

// Every optional attribute is either clearable or excluded with a reason; a new
// one fails this test until it is classified. This resource has no optional
// attributes today, so both lists are empty and the check passes.
func TestUnitClearableClassificationIsComplete(t *testing.T) {
	t.Parallel()

	problems := clearable.Check(
		NetworkFloatingIpResourceSchema(context.Background()).Attributes,
		clearableAttributes,
		notClearable,
	)
	if len(problems) != 0 {
		t.Fatalf("clearable classification:\n  %s", strings.Join(problems, "\n  "))
	}
}
