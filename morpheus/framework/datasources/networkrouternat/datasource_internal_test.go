// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package networkrouternat

import (
	"context"
	"strings"
	"testing"
)

// TestUnitDeprecatedAttributes_MORPH14596 pins the deprecation of action and
// firewall in the generated schema. The NAT read API never returns either, so
// they are deprecated ahead of removal in 2.6: both must stay computed-only
// (never settable), carry a deprecation message, and say so in their
// description. The message is produced by the codegen pipeline's
// deprecation_message override, which older generator builds silently drop on
// data sources; this test turns that regression into a build failure.
func TestUnitDeprecatedAttributes_MORPH14596(t *testing.T) {
	t.Parallel()

	s := NetworkRouterNatDataSourceSchema(context.Background())

	for _, name := range []string{"action", "firewall"} {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			a, ok := s.Attributes[name]
			if !ok {
				t.Fatalf("schema has no %s attribute", name)
			}
			if !a.IsComputed() || a.IsOptional() || a.IsRequired() {
				t.Errorf("%s must be computed-only; computed=%v optional=%v required=%v",
					name, a.IsComputed(), a.IsOptional(), a.IsRequired())
			}
			if a.GetDeprecationMessage() == "" {
				t.Errorf("%s has no deprecation message", name)
			}
			if !strings.HasPrefix(a.GetDescription(), "Deprecated.") {
				t.Errorf("%s description does not start with %q: %q", name, "Deprecated.", a.GetDescription())
			}
			if !strings.Contains(a.GetDeprecationMessage(), "2.6") {
				t.Errorf("%s deprecation message does not name the removal release: %q", name, a.GetDeprecationMessage())
			}
		})
	}

	// Nothing else on the data source is deprecated; a stray message here would
	// mean an override landed on the wrong attribute.
	for name, a := range s.Attributes {
		if name == "action" || name == "firewall" {
			continue
		}
		if a.GetDeprecationMessage() != "" {
			t.Errorf("%s is unexpectedly deprecated: %q", name, a.GetDeprecationMessage())
		}
	}
}
