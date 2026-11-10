// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package networkserver

import (
	"testing"

	"github.com/HPE/terraform-provider-hpe/morpheus/utils/networkserverlookup"
)

// TestUnitErrorStringsMatchLookup pins the exported error constants to the
// lookup package's sentinel errors, which are what the data source reports.
func TestUnitErrorStringsMatchLookup(t *testing.T) {
	t.Parallel()

	if got := networkserverlookup.ErrNotFound.Error(); got != ErrorNoNetworkServerFound {
		t.Errorf("ErrNotFound = %q, const = %q", got, ErrorNoNetworkServerFound)
	}

	if got := networkserverlookup.ErrMultiple.Error(); got != ErrorMultipleNetworkServers {
		t.Errorf("ErrMultiple = %q, const = %q", got, ErrorMultipleNetworkServers)
	}
}
