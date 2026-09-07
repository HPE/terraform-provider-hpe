// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package task_test

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	sdkv2morpheus "github.com/HPE/terraform-provider-hpe/morpheus/sdkv2"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers/capabilities"
	"github.com/HPE/terraform-provider-hpe/provider/adapter"
)

// TestAccMorpheusTaskAnsiblePlaybookValidationEmptyPlaybook verifies that an
// empty playbook value is rejected at plan time rather than being sent to the
// API (MORPH-14579).
func TestAccMorpheusTaskAnsiblePlaybookValidationEmptyPlaybook(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.All)

	providerConfig := testhelpers.ProviderBlock()

	config := providerConfig + `
resource "hpe_morpheus_task_ansible_playbook" "test" {
	name     = "TestAccMorpheusTaskAnsiblePlaybookValidationEmptyPlaybook"
	playbook = ""
}
`

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, adapter.NewMorpheus(), sdkv2morpheus.Provider()),
		Steps: []resource.TestStep{
			{
				Config:      config,
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`to not be an empty string`),
			},
		},
	})
}
