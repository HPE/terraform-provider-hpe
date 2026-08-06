package instancenode

import (
	"context"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/resource"
)

// Import is refused deliberately, so assert both that it fails and that it
// explains why. A bare failure would be indistinguishable from the
// framework's own default and would tell the practitioner nothing.
func TestUnitImportState_RefusesWithReason(t *testing.T) {
	t.Parallel()

	r := &Resource{}
	resp := &resource.ImportStateResponse{}

	r.ImportState(context.Background(), resource.ImportStateRequest{ID: "123"}, resp)

	if !resp.Diagnostics.HasError() {
		t.Fatal("expected import to be refused, got no error")
	}

	detail := resp.Diagnostics.Errors()[0].Detail()

	// The reason must name the attribute that makes import unsafe, so the
	// message stays useful if the refusal is ever revisited.
	if !strings.Contains(detail, "pre-provisioned") {
		t.Errorf(
			"refusal should explain that a node's origin cannot be read back; got: %s",
			detail,
		)
	}
}
