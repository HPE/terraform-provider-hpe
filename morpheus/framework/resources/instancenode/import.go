package instancenode

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"
)

// ImportState is implemented solely to refuse import with a reason.
//
// A node could be located from a container id, and most of its state read
// back — the server, the address, and the pool the server currently sits in.
// What cannot be recovered is how the node was created: `pre_provisioned`
// records whether an existing server was attached rather than a new one
// provisioned, and nothing in the API reports that after the fact.
//
// Since `pre_provisioned` forces replacement, an imported node would plan a
// destroy and recreate for any configuration that sets it — returning a
// working server to its pool and building another for no reason. Refusing the
// import is safer than adopting a resource that cannot describe itself.
//
// Without this method the framework reports only that the resource does not
// support import, which does not say why or what to do instead.
func (r *Resource) ImportState(
	_ context.Context,
	_ resource.ImportStateRequest,
	resp *resource.ImportStateResponse,
) {
	resp.Diagnostics.AddError(
		"import hpe_morpheus_instance_node resource",
		"importing an instance node is not supported. A node's origin cannot "+
			"be read back from the API: whether it attached a pre-provisioned "+
			"server or provisioned a new one is not recorded, and that "+
			"attribute forces replacement, so an imported node would plan an "+
			"unnecessary destroy and recreate. Add nodes through Terraform "+
			"rather than adopting existing ones.",
	)
}
