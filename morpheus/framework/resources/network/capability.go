// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package network

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"

	sdk "github.com/HPE/terraform-provider-hpe/internal/sdk/oapigen"
	"github.com/HPE/terraform-provider-hpe/morpheus/utils/networkconfig"
	"github.com/HPE/terraform-provider-hpe/morpheus/utils/versioncheck"
)

// checkDeclared compares the configured type-specific values with what the
// appliance's own definition of the network type declares. A typed block
// carries every field the type has on the release the block was written
// against; an older release, or an older version of the plugin that provides
// the type, may lack some of them — the virtual-switch fields of the HVM data
// and private networks arrived with Morpheus 9.1, say — and Morpheus applies
// only the option types the type declares, so a value for a field it does
// not have is accepted and ignored. That is the silent failure the typed
// blocks exist to prevent, so it is an error for a block attribute and a
// warning for a key of the free-form config map, which is the escape hatch.
//
// Only a type read from the appliance is authoritative: the generated table
// is another release's view, so nothing is checked when it is the fallback.
// The appliance version is looked up only when there is something to report.
func checkDeclared(ctx context.Context, client *sdk.APIClient, inputs configInputs) diag.Diagnostics {
	var diags diag.Diagnostics

	typ := inputs.typ
	if typ == nil || !typ.Live {
		return diags
	}

	var missing []networkconfig.BlockField
	if inputs.block != nil {
		missing = inputs.block.Undeclared(*typ, inputs.blockValues)
	}

	undeclared := networkconfig.UndeclaredKeys(*typ, inputs.generic)

	if len(missing) == 0 && len(undeclared) == 0 {
		return diags
	}

	where := "this appliance" + applianceVersion(ctx, client)

	// The block's path, not the attribute's: Terraform maps a write-only
	// nested attribute's path back to the block, and the message names the
	// attribute anyway.
	for _, f := range missing {
		diags.AddAttributeError(path.Root(inputs.block.Attr),
			"attribute not available on this appliance",
			fmt.Sprintf("Network type %s (%s) on %s does not declare %s, which %s.%s sets; "+
				"Morpheus would accept the value and ignore it. The attribute needs a newer Morpheus release, "+
				"or a newer version of the plugin that provides the network type.",
				typ.Code, typ.Name, where, f.Key, inputs.block.Attr, f.Attr))
	}

	if len(undeclared) > 0 {
		diags.AddAttributeWarning(path.Root("config"),
			"config keys the network type does not declare",
			fmt.Sprintf("Network type %s (%s) on %s does not declare %s; Morpheus may ignore them.",
				typ.Code, typ.Name, where, strings.Join(undeclared, ", ")))
	}

	return diags
}

// applianceVersion returns " (Morpheus x.y.z.n)" for the message, or "" when
// the version cannot be read.
func applianceVersion(ctx context.Context, client *sdk.APIClient) string {
	if client == nil {
		return ""
	}

	v, err := versioncheck.Appliance(ctx, client)
	if err != nil || v == nil {
		return ""
	}

	return " (Morpheus " + v.Original() + ")"
}
