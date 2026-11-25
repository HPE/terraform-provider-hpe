// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package settingbackup

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

// ImportState adopts the backup settings. They are a singleton, so the import
// ID identifies nothing and any value works ("1", or "settings"); the id is
// "1" whatever it is, as it was in the SDKv2 resource. Read fills in the rest.
func (r *Resource) ImportState(
	ctx context.Context,
	_ resource.ImportStateRequest,
	resp *resource.ImportStateResponse,
) {
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), singletonID)...)
}
