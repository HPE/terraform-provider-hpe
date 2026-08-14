// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package storagecontrollertype_test

import (
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers/capabilities"
	"github.com/HPE/terraform-provider-hpe/provider/adapter"
)

// controllerName is the SCSI VMware Paravirtual controller, which the "vmware"
// provision type seeds on every appliance (independent of whether a VMware
// cloud is configured), so no environment fixture is required.
const controllerName = "SCSI VMware Paravirtual"

func TestMain(m *testing.M) {
	code := m.Run()
	testhelpers.WriteMergedResults()
	os.Exit(code)
}

func TestAccMorpheusFindStorageControllerType(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.All)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	t.Parallel()

	providerConfig := testhelpers.ProviderBlock()

	dataSourceConfig, err := testhelpers.RenderExample(t,
		"example.tf.tmpl",
		"ControllerName", controllerName,
		"BusNumber", "1",
		"InterfaceNumber", "0",
	)
	if err != nil {
		t.Fatal(err)
	}

	checks := []resource.TestCheckFunc{
		resource.TestCheckResourceAttrSet(
			"data.hpe_morpheus_storage_controller_type.example",
			"id",
		),
		resource.TestCheckResourceAttrSet(
			"data.hpe_morpheus_storage_controller_type.example",
			"category",
		),
		// The mount point must have the format -1:<bus>:<typeId>:<unit>. bus and
		// unit are the inputs (1 and 0); typeId is not pinned as it varies by
		// appliance.
		resource.TestMatchResourceAttr(
			"data.hpe_morpheus_storage_controller_type.example",
			"controller_mount_point",
			regexp.MustCompile(`^-1:1:[0-9]+:0$`),
		),
		resource.TestCheckResourceAttr(
			"data.hpe_morpheus_storage_controller_type.example",
			"provision_type_code",
			"vmware",
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

func TestAccMorpheusFindStorageControllerTypeNotFound(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.All)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	t.Parallel()

	providerConfig := testhelpers.ProviderBlock()

	dataSourceConfig, err := testhelpers.RenderExample(t,
		"example.tf.tmpl",
		"ControllerName", "__nonexistent_controller__",
		"BusNumber", "1",
		"InterfaceNumber", "0",
	)
	if err != nil {
		t.Fatal(err)
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, adapter.NewMorpheus(), nil),
		Steps: []resource.TestStep{
			{
				Config:      providerConfig + dataSourceConfig,
				ExpectError: regexp.MustCompile(`no storage controller type found`),
			},
		},
	})
}
