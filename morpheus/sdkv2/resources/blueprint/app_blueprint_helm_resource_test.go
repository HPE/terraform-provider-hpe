// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package blueprint_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	sdkv2morpheus "github.com/HPE/terraform-provider-hpe/morpheus/sdkv2"
	"github.com/HPE/terraform-provider-hpe/morpheus/sdkv2/resources/blueprint"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers/capabilities"
	"github.com/HPE/terraform-provider-hpe/provider/adapter"
)

func TestAccMorpheusAppBlueprintHelmExampleOk(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.All, capabilities.Git)

	integrationName := testhelpers.GitIntegrationName(t)
	repositoryName := testhelpers.GitRepositoryName(t)
	branch := testhelpers.GitRepositoryBranch(t)

	t.Parallel()

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	// t.Skip("Skipping due to missing infrastructure in test environment")

	providerConfig := testhelpers.ProviderBlock()

	name := acctest.RandomWithPrefix(t.Name())

	resourceConfig, err := blueprint.RenderAppBlueprintHelmConfig(t, map[string]string{
		"Name":          name,
		"IntegrationId": testhelpers.GitIntegrationIDRef,
		"RepositoryId":  testhelpers.GitRepositoryIDRef(repositoryName),
		"VersionRef":    testhelpers.HCLEscape(branch),
	})

	dependenciesConfig := testhelpers.GitIntegrationBlock(integrationName)
	if err != nil {
		t.Fatal(err)
	}

	checks := []resource.TestCheckFunc{
		resource.TestCheckResourceAttr(
			"hpe_morpheus_app_blueprint_helm.example",
			"category",
			"helm",
		),

		resource.TestCheckResourceAttr(
			"hpe_morpheus_app_blueprint_helm.example",
			"description",
			"tf example helm app blueprint",
		),

		resource.TestCheckResourceAttrPair(
			"hpe_morpheus_app_blueprint_helm.example", "integration_id",
			"data.hpe_morpheus_integration_git.git", "id",
		),

		resource.TestCheckResourceAttr(
			"hpe_morpheus_app_blueprint_helm.example",
			"name",
			name,
		),

		resource.TestCheckResourceAttrPair(
			"hpe_morpheus_app_blueprint_helm.example", "repository_id",
			"data.hpe_morpheus_integration_git.git", "repository_ids."+repositoryName,
		),

		resource.TestCheckResourceAttr(
			"hpe_morpheus_app_blueprint_helm.example",
			"version_ref",
			branch,
		),

		resource.TestCheckResourceAttr(
			"hpe_morpheus_app_blueprint_helm.example",
			"working_path",
			"./test",
		),
	}

	checkFn := resource.ComposeAggregateTestCheckFunc(checks...)
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, adapter.NewMorpheus(), sdkv2morpheus.Provider()),
		Steps: []resource.TestStep{
			// Apply
			{
				Config:             providerConfig + dependenciesConfig + resourceConfig,
				ExpectNonEmptyPlan: false,
				Check:              checkFn,
			},
			// Plan after apply
			{
				Config:             providerConfig + dependenciesConfig + resourceConfig,
				ExpectNonEmptyPlan: false,
				PlanOnly:           true,
			},
		},
	})
}
