// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package instanceclone_test

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers/capabilities"
	"github.com/HPE/terraform-provider-hpe/provider/adapter"
)

// TestAccMorpheusInstanceCloneResourceWriteOnlyOk: the write-only
// source_instance_id is used only at creation (it is the instance the clone is
// taken from), so a change plans a replacement (MORPH-18005); the deprecated
// marker still replaces when bumped and does nothing when removed; both import
// forms leave a clean plan. Needs two source instances: TF_ACC_INSTANCE_ID and
// TF_ACC_INSTANCE_ID_2 (and TF_ACC_NETWORK_ID), as the other clone tests do.
func TestAccMorpheusInstanceCloneResourceWriteOnlyOk(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.All)
	testhelpers.SkipUnlessAcc(t)

	src1, src2, networkID := os.Getenv("TF_ACC_INSTANCE_ID"), os.Getenv("TF_ACC_INSTANCE_ID_2"), os.Getenv("TF_ACC_NETWORK_ID")
	if src1 == "" || src2 == "" || networkID == "" {
		t.Skip("TF_ACC_INSTANCE_ID, TF_ACC_INSTANCE_ID_2 and TF_ACC_NETWORK_ID must be set")
	}

	name := acctest.RandomWithPrefix("tfwo-clone")

	c := testhelpers.WriteOnlyCase{
		Address:   "hpe_morpheus_instance_clone.wo",
		Attribute: "source_instance_id",
		Marker:    "source_instance_id_version",
		// The source instance is in the URL (/api/instances/{id}/clone), not the
		// body: any write to the clone path counts as "sent".
		APIPath: "/clone",
		APIKey:  "",
		// The clone updates only its plan and volume sizes in place, so the
		// matrix's unrelated-change step is left out (unrelated is ignored).
		Config: func(value, extra string, _ bool) string {
			return testhelpers.ProviderBlock() + `
resource "hpe_morpheus_instance_clone" "wo" {
  source_instance_id = ` + value + extra + `
  name               = "` + name + `"

  volumes = [{ name = "root", size = 20, root_volume = true }]

  network_interfaces = [{ network_id = ` + networkID + `, ip_mode = "" }]
}
`
		},
	}

	rec := testhelpers.RecordRequests(t)

	steps := testhelpers.WriteOnlyReplaceSteps(c, rec, src1, src2)
	steps = append(steps, testhelpers.WriteOnlyImportSteps(c, src2, "volumes", "network_interfaces")...)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, adapter.NewMorpheus(), nil),
		Steps:                    steps,
	})
}
