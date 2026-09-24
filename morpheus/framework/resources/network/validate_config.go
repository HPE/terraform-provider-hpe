// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package network

import (
	"context"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

var _ resource.ResourceWithValidateConfig = &Resource{}

// ValidateConfig checks the typed configuration blocks against each other and
// against type_code without touching the appliance, so `terraform validate`
// catches them: at most one block may be set, and a block must be the one for
// the configured network type. When the type is given as type_id, or is not
// in the catalogue, the check waits for ModifyPlan and Create, which know it.
func (r *Resource) ValidateConfig(
	ctx context.Context,
	req resource.ValidateConfigRequest,
	resp *resource.ValidateConfigResponse,
) {
	var config NetworkModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	blocks, diags := configuredBlocks(ctx, config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if len(blocks) > 1 {
		names := make([]string, 0, len(blocks))
		for _, b := range blocks {
			names = append(names, b.block.Attr)
		}

		for _, b := range blocks {
			resp.Diagnostics.AddAttributeError(path.Root(b.block.Attr),
				"only one configuration block may be set",
				"A network has one type, so it takes one typed configuration block; got "+strings.Join(names, ", ")+".")
		}

		return
	}

	if len(blocks) == 0 {
		return
	}

	code := knownString(config.TypeCode)
	if code == "" {
		return
	}

	typ := typeSpecFor(code)
	if typ == nil {
		// A type the catalogue does not know (a plugin's, say) has no block, so
		// a block set with it is wrong — but that is Create's message once it
		// has confirmed the type exists.
		return
	}

	if d := checkBlockMatchesType(&blocks[0].block, typ); d != nil {
		resp.Diagnostics.Append(d)
	}
}
