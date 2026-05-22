// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package image_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"

	"github.com/HPE/terraform-provider-hpe/morpheus"
	sdkv2morpheus "github.com/HPE/terraform-provider-hpe/morpheus/sdkv2"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers/capabilities"
)

// nolint: lll
var url = "https://dl-cdn.alpinelinux.org/alpine/v3.22/releases/cloud/generic_alpine-3.22.2-x86_64-bios-cloudinit-r0.qcow2"

func TestAccMorpheusImageUpdate(t *testing.T) {
	if capabilities.Missing(t, capabilities.Alletra) {
		t.Log("Skipping test due to missing capabilities")

		return
	}
	defer testhelpers.RecordResult(t)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	t.Parallel()

	providerConfig := testhelpers.ProviderBlock()
	name := acctest.RandomWithPrefix(t.Name())

	datasourceConfig := `
data "hpe_morpheus_os_type" "test" {
	name = "linux"
}

data "hpe_morpheus_storage_bucket" "test" {
	name = "Local Storage"
}
`

	checks := []resource.TestCheckFunc{
		resource.TestCheckResourceAttr(
			"hpe_morpheus_image.example",
			"name",
			name+"2",
		),
	}

	checkFn := resource.ComposeAggregateTestCheckFunc(
		checks...,
	)

	checkReplace := resource.ConfigPlanChecks{
		PreApply: []plancheck.PlanCheck{
			plancheck.ExpectResourceAction(
				"hpe_morpheus_image.example",
				plancheck.ResourceActionReplace,
			),
		},
	}

	// this check verifies that the resource is going to be updated in place
	checkInPlaceUpdate := resource.ConfigPlanChecks{
		PreApply: []plancheck.PlanCheck{
			plancheck.ExpectResourceAction(
				"hpe_morpheus_image.example",
				plancheck.ResourceActionUpdate,
			),
		},
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, morpheus.New(), sdkv2morpheus.Provider()),
		Steps: []resource.TestStep{
			{
				Config: providerConfig + datasourceConfig + `
					resource "hpe_morpheus_image" "example" {
						name = "` + name + `"
						description = "this is a test image"
						labels = ["terraform-image"]
						image_type = "qcow2"
						url = "` + url + `"
						storage_provider_id = data.hpe_morpheus_storage_bucket.test.id
						install_agent = false
						cloud_init = false
						os_type_id = data.hpe_morpheus_os_type.test.id
						min_ram = 2
						uefi = false
						min_disk = 25
						trial_version = false
						ssh_username = "test-user"
						user_data = trimspace(
						<<-EOT
						#!/bin/sh
						apk add --no-cache bash
						EOT
						)
						tags = [
						{
							name = "test-tag"
							value = "true"
						}
						]
						virtio_supported = false
						visibility = "private"
					}`,
				PlanOnly: false,
			},
			// change fields which don't require a replace
			{
				Config: providerConfig + datasourceConfig + `
					resource "hpe_morpheus_image" "example" {
						name = "` + name + `2"
						description = "this is a test image" # requires replace
						labels = ["terraform-image", "terraform-test"]
						image_type = "qcow2" # requires replace
						url = "` + url + `" # requires replace
						storage_provider_id = data.hpe_morpheus_storage_bucket.test.id # requires replace
						install_agent = true
						cloud_init = true
						os_type_id = data.hpe_morpheus_os_type.test.id
						min_ram = 1
						uefi = true
						min_disk = 20
						trial_version = true
						ssh_username = "test-user2"
						user_data = trimspace(
						<<-EOT
						#!/bin/sh
						apk add --no-cache bash go
						EOT
						)
						tags = [
							{
								name = "test-tag"
								value = "true"
							},
							{
								name = "test-tag2"
								value = "false"
							}
						]
						virtio_supported = true
						visibility = "public"
					}`,
				Check:            checkFn,
				PlanOnly:         false,
				ConfigPlanChecks: checkInPlaceUpdate,
			},
			{
				Config: providerConfig + datasourceConfig + `
					resource "hpe_morpheus_image" "example" {
						name = "` + name + `2"
						description = "this is a test image" # requires replace
						image_type = "qcow2" # requires replace
						url = "` + url + `" # requires replace
						storage_provider_id = data.hpe_morpheus_storage_bucket.test.id # requires replace
					}`,
				Check:              checkFn,
				ExpectNonEmptyPlan: true,
				PlanOnly:           false,
				ConfigPlanChecks:   checkInPlaceUpdate,
			},
			{
				Config: providerConfig + datasourceConfig + `
					resource "hpe_morpheus_image" "example" {
						name = "` + name + `2"
						description = "this is a test image" # requires replace
						image_type = "qcow2" # requires replace
						url = "` + url + `" # requires replace
						storage_provider_id = data.hpe_morpheus_storage_bucket.test.id # requires replace
						ssh_password_wo = "this-is-a-test-password2"
						ssh_password_wo_version = 1
					}`,
				Check:              checkFn,
				ExpectNonEmptyPlan: true,
				PlanOnly:           false,
				ConfigPlanChecks:   checkInPlaceUpdate,
			},
			// change description to force a replace
			{
				Config: providerConfig + datasourceConfig + `
					resource "hpe_morpheus_image" "example" {
						name = "` + name + `2"
						description = "this is a test image2" # requires replace
						labels = ["terraform-image", "terraform-test"]
						image_type = "qcow2" # requires replace
						url = "` + url + `" # requires replace
						storage_provider_id = data.hpe_morpheus_storage_bucket.test.id # requires replace
						install_agent = true
						cloud_init = true
						os_type_id = data.hpe_morpheus_os_type.test.id
						min_ram = 1
						uefi = true
						min_disk = 20
						trial_version = true
						ssh_username = "test-user2"
						user_data = trimspace(
						<<-EOT
						#!/bin/sh
						apk add --no-cache bash go
						EOT
						)
						tags = [
						{
							name = "test-tag"
							value = "true"
						}
						]
						virtio_supported = true
						visibility = "public"
					}`,
				Check:            checkFn,
				PlanOnly:         false,
				ConfigPlanChecks: checkReplace,
			},
		},
	})
}

func TestAccMorpheusImageUpdatePassword(t *testing.T) {
	if capabilities.Missing(t, capabilities.Alletra) {
		t.Log("Skipping test due to missing capabilities")

		return
	}
	defer testhelpers.RecordResult(t)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	t.Parallel()

	providerConfig := testhelpers.ProviderBlock()
	name := acctest.RandomWithPrefix(t.Name())

	datasourceConfig := `
data "hpe_morpheus_storage_bucket" "test" {
	name = "Local Storage"
}
`

	checks := []resource.TestCheckFunc{
		resource.TestCheckResourceAttr(
			"hpe_morpheus_image.example",
			"name",
			name,
		),
		resource.TestCheckNoResourceAttr(
			"hpe_morpheus_image.example",
			"ssh_password_wo",
		),
	}

	checkFn := resource.ComposeAggregateTestCheckFunc(
		checks...,
	)
	// no changes in plan
	checkNoOp := resource.ConfigPlanChecks{
		PreApply: []plancheck.PlanCheck{
			plancheck.ExpectResourceAction(
				"hpe_morpheus_image.example",
				plancheck.ResourceActionNoop,
			),
		},
	}

	// this check verifies that the resource is going to be updated in place
	checkInPlaceUpdate := resource.ConfigPlanChecks{
		PreApply: []plancheck.PlanCheck{
			plancheck.ExpectResourceAction(
				"hpe_morpheus_image.example",
				plancheck.ResourceActionUpdate,
			),
		},
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, morpheus.New(), sdkv2morpheus.Provider()),
		Steps: []resource.TestStep{
			{
				Config: providerConfig + datasourceConfig + `
					resource "hpe_morpheus_image" "example" {
						name = "` + name + `"
						description = "this is a test image" # requires replace
						image_type = "qcow2" # requires replace
						storage_provider_id = data.hpe_morpheus_storage_bucket.test.id # requires replace
						ssh_password_wo = "this-is-a-test-password"
						ssh_password_wo_version = 1
					}`,
				Check:    checkFn,
				PlanOnly: false,
			},
			// changing the password without bumping the `wo_version` should not
			// change the password
			{
				Config: providerConfig + datasourceConfig + `
					resource "hpe_morpheus_image" "example" {
						name = "` + name + `"
						description = "this is a test image" # requires replace
						image_type = "qcow2" # requires replace
						storage_provider_id = data.hpe_morpheus_storage_bucket.test.id # requires replace
						ssh_password_wo = "this-is-a-test-password2"
						ssh_password_wo_version = 1
					}`,
				Check:            checkFn,
				ConfigPlanChecks: checkNoOp,
				PlanOnly:         false,
			},
			// removing the password field should not cause any changes
			{
				Config: providerConfig + datasourceConfig + `
					resource "hpe_morpheus_image" "example" {
						name = "` + name + `"
						description = "this is a test image" # requires replace
						image_type = "qcow2" # requires replace
						storage_provider_id = data.hpe_morpheus_storage_bucket.test.id # requires replace
						ssh_password_wo_version = 1
					}`,
				Check:            checkFn,
				ConfigPlanChecks: checkNoOp,
				PlanOnly:         false,
			},
			// bumping the password version and changing the password attribute
			// should trigger an in place update
			{
				Config: providerConfig + datasourceConfig + `
					resource "hpe_morpheus_image" "example" {
						name = "` + name + `"
						description = "this is a test image" # requires replace
						image_type = "qcow2" # requires replace
						storage_provider_id = data.hpe_morpheus_storage_bucket.test.id # requires replace
						ssh_password_wo = "this-is-a-test-password2"
						ssh_password_wo_version = 2
					}`,
				Check:            checkFn,
				PlanOnly:         false,
				ConfigPlanChecks: checkInPlaceUpdate,
			},
		},
	})
}
