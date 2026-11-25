// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package settingbackup

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Delete forgets the resource. The settings always exist and the API has no
// way to remove or reset them, so the appliance keeps what it has, as it did
// with the SDKv2 resource.
func (r *Resource) Delete(
	ctx context.Context,
	_ resource.DeleteRequest,
	_ *resource.DeleteResponse,
) {
	tflog.Info(ctx, "backup settings cannot be deleted; removing the resource from state only")
}
