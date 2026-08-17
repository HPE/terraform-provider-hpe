// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package instancedisktype_test

import (
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers/capabilities"
	"github.com/HPE/terraform-provider-hpe/provider/adapter"
)

// These are long-lived, seeded resources on the test appliance, matching the
// fixtures the hpe_morpheus_instance acceptance test uses: the "hvm" cloud, the
// "Single KVM VM" layout (id 77) and the default group (id 1). The "Standard"
// disk type is available for provisioning in every cloud, so it is a stable
// lookup target. The cloud id is resolved by name through a hpe_morpheus_cloud
// data source (rather than hard-coded) so the test survives cloud id changes.
const (
	cloudName    = "hvm"
	layoutID     = "77"
	groupID      = "1"
	diskTypeName = "Standard"
)

func TestMain(m *testing.M) {
	code := testhelpers.TestMain(m)
	testhelpers.WriteMergedResults()
	os.Exit(code)
}

func TestAccMorpheusFindInstanceDiskType(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.All)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	t.Parallel()

	providerConfig := testhelpers.ProviderBlock()

	dataSourceConfig, err := testhelpers.RenderExample(t,
		"example.tf.tmpl",
		"CloudName", cloudName,
		"Name", diskTypeName,
		"LayoutId", layoutID,
		"GroupId", groupID,
	)
	if err != nil {
		t.Fatal(err)
	}

	checks := []resource.TestCheckFunc{
		resource.TestCheckResourceAttrSet(
			"data.hpe_morpheus_instance_disk_type.example",
			"id",
		),
		resource.TestCheckResourceAttr(
			"data.hpe_morpheus_instance_disk_type.example",
			"name",
			diskTypeName,
		),
	}

	checkFn := resource.ComposeAggregateTestCheckFunc(checks...)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, adapter.NewMorpheus(), nil),
		Steps: []resource.TestStep{
			{
				Config: providerConfig + dataSourceConfig,
				Check:  checkFn,
			},
		},
	})
}

func TestAccMorpheusFindInstanceDiskTypeNotFound(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.All)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	t.Parallel()

	providerConfig := testhelpers.ProviderBlock()

	dataSourceConfig, err := testhelpers.RenderExample(t,
		"example.tf.tmpl",
		"CloudName", cloudName,
		"Name", "__nonexistent_disk_type__",
		"LayoutId", layoutID,
		"GroupId", groupID,
	)
	if err != nil {
		t.Fatal(err)
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, adapter.NewMorpheus(), nil),
		Steps: []resource.TestStep{
			{
				Config:      providerConfig + dataSourceConfig,
				ExpectError: regexp.MustCompile(`no instance disk type found`),
			},
		},
	})
}
