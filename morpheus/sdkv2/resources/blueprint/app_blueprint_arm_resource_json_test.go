// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package blueprint_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/HPE/terraform-provider-hpe/morpheus"
	sdkv2morpheus "github.com/HPE/terraform-provider-hpe/morpheus/sdkv2"
	"github.com/HPE/terraform-provider-hpe/morpheus/sdkv2/resources/blueprint"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers/systemoverride"
)

func TestAccMorpheusAppBlueprintArmJsonExampleOk(t *testing.T) {
	t.Parallel()

	defer testhelpers.RecordResult(t)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	// t.Skip("Skipping due to missing infrastructure in test environment")

	testSystem := systemoverride.GetPreferred(t, "feature")
	providerConfig := testhelpers.ProviderBlock(testSystem)

	name := acctest.RandomWithPrefix(t.Name())

	resourceConfig, err := blueprint.RenderAppBlueprintArmJSONConfig(t, map[string]string{
		"Name": name,
	})
	if err != nil {
		t.Fatal(err)
	}

	checks := []resource.TestCheckFunc{
		// check diff suppress functions
		// resource.TestCheckResourceAttr(
		// 	"hpe_morpheus_app_blueprint_arm.example",
		// 	"blueprint_content",
		//nolint:lll // 	"{\n \"type\": \"Microsoft.Storage/storageAccounts\",\n \"apiVersion\": \"2019-04-01\",\n \"name\": \"string\",\n \"location\": \"string\",\n \"tags\": {\n \"tagName1\": \"tagValue1\",\n \"tagName2\": \"tagValue2\"\n },\n \"sku\": {\n \"name\": \"string\",\n \"restrictions\": [\n {\n \"reasonCode\": \"string\"\n }\n ]\n },\n \"kind\": \"string\",\n \"identity\": {\n \"type\": \"SystemAssigned\"\n },\n \"properties\": {\n \"accessTier\": \"string\",\n \"allowBlobPublicAccess\": \"bool\",\n \"allowSharedKeyAccess\": \"bool\",\n \"azureFilesIdentityBasedAuthentication\": {\n \"activeDirectoryProperties\": {\n \"azureStorageSid\": \"string\",\n \"domainGuid\": \"string\",\n \"domainName\": \"string\",\n \"domainSid\": \"string\",\n \"forestName\": \"string\",\n \"netBiosDomainName\": \"string\"\n },\n \"directoryServiceOptions\": \"string\"\n },\n \"customDomain\": {\n \"name\": \"string\",\n \"useSubDomainName\": \"bool\"\n },\n \"encryption\": {\n \"keySource\": \"string\",\n \"keyvaultproperties\": {\n \"keyname\": \"string\",\n \"keyvaulturi\": \"string\",\n \"keyversion\": \"string\"\n },\n \"services\": {\n \"blob\": {\n \"enabled\": \"bool\"\n },\n \"file\": {\n \"enabled\": \"bool\"\n }\n }\n },\n \"isHnsEnabled\": \"bool\",\n \"largeFileSharesState\": \"string\",\n \"minimumTlsVersion\": \"string\",\n \"networkAcls\": {\n \"bypass\": \"string\",\n \"defaultAction\": \"string\",\n \"ipRules\": [\n {\n \"action\": \"Allow\",\n \"value\": \"string\"\n }\n ],\n \"virtualNetworkRules\": [\n {\n \"action\": \"Allow\",\n \"id\": \"string\",\n \"state\": \"string\"\n }\n ]\n },\n \"supportsHttpsTrafficOnly\": \"bool\"\n }\n}",
		// ),

		resource.TestCheckResourceAttr(
			"hpe_morpheus_app_blueprint_arm.example",
			"category",
			"armtemplates",
		),

		resource.TestCheckResourceAttr(
			"hpe_morpheus_app_blueprint_arm.example",
			"cloud_init_enabled",
			"true",
		),

		resource.TestCheckResourceAttr(
			"hpe_morpheus_app_blueprint_arm.example",
			"description",
			"example arm app blueprint",
		),

		resource.TestCheckResourceAttr(
			"hpe_morpheus_app_blueprint_arm.example",
			"install_agent",
			"true",
		),

		resource.TestCheckResourceAttr(
			"hpe_morpheus_app_blueprint_arm.example",
			"name",
			name,
		),

		resource.TestCheckResourceAttr(
			"hpe_morpheus_app_blueprint_arm.example",
			"os_type",
			"linux",
		),

		resource.TestCheckResourceAttr(
			"hpe_morpheus_app_blueprint_arm.example",
			"source_type",
			"json",
		),
	}

	checkFn := resource.ComposeAggregateTestCheckFunc(checks...)
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, morpheus.New(), sdkv2morpheus.Provider()),
		Steps: []resource.TestStep{
			// Apply
			{
				Config:             providerConfig + resourceConfig,
				ExpectNonEmptyPlan: true,
				Check:              checkFn,
			},
			// Plan after apply
			{
				Config:             providerConfig + resourceConfig,
				ExpectNonEmptyPlan: true,
				PlanOnly:           true,
			},
		},
	})
}
