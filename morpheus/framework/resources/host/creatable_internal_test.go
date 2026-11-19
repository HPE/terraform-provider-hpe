// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package host

import (
	"strings"
	"testing"
)

// The types a practitioner is most likely to try and cannot create here each
// name the resource that owns them.
func TestUnitNotCreatableHintsNameOwner(t *testing.T) {
	t.Parallel()

	want := map[string]string{
		"mvm":        "hpe_morpheus_cluster",
		"morpheusVm": "hpe_morpheus_instance",
		"vmwareVm":   "hpe_morpheus_instance",
	}
	for code, owner := range want {
		hint, ok := notCreatableHints[code]
		if !ok {
			t.Errorf("%s: no hint", code)

			continue
		}
		if !strings.Contains(hint, owner) {
			t.Errorf("%s: hint %q does not name %s", code, hint, owner)
		}
	}
}
