// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package adapter

import (
	opsramp "github.com/HPE/terraform-provider-hpe/opsramp/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider"
)

// Returns an OpsRamp provider via the Adapter Layer
// This saves having to manually wrap an OpsRamp provider where an
// adapted OpsRamp provider is needed, e.g. Acceptance Tests.
func NewOpsRamp() provider.Provider {
	opsrampProvider := opsramp.New(opsramp.Version)()
	return NewAdaptedProvider(opsrampProvider)
}
