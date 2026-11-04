// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package settingappliance

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

// ImportState adopts the appliance settings. They are a singleton, so the
// import ID identifies nothing and any value works ("1", or "settings"); the
// id is "1" whatever it is, as it was in the SDKv2 resource. The ID may carry
// the sha256 digests of the configured credentials --
// "1#smtp_password=<hex>&proxy_password=<hex>", see wodigest -- so that a later
// change to one is detected from the start; without them the configured
// credentials are taken as already applied. Read fills in the rest.
func (r *Resource) ImportState(
	ctx context.Context,
	req resource.ImportStateRequest,
	resp *resource.ImportStateResponse,
) {
	if _, ok := wo.Import(ctx, req, resp); !ok {
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), singletonID)...)
}
