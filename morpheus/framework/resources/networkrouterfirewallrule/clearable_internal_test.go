// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package networkrouterfirewallrule

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

	schemaAttrs := NetworkRouterFirewallRuleResourceSchema(context.Background()).Attributes
	problems := clearable.Check(schemaAttrs, clearableAttributes, notClearable)
	if len(problems) != 0 {
		t.Fatalf("clearable classification:\n  %s", strings.Join(problems, "\n  "))
	}
}
