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
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers/capabilities"
)

func TestAccMorpheusAppBlueprintCloudFormationJsonExampleOk(t *testing.T) {
	if capabilities.Missing(t, capabilities.All, capabilities.AWS) {
		t.Log("Skipping test due to missing capabilities")

		return
	}
	t.Parallel()

	defer testhelpers.RecordResult(t)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	// t.Skip("Skipping due to missing infrastructure in test environment")

	providerConfig := testhelpers.ProviderBlock()

	name := acctest.RandomWithPrefix(t.Name())

	resourceConfig, err := blueprint.RenderAppBlueprintCloudFormationJSONConfig(t, map[string]string{
		"Name": name,
	})
	if err != nil {
		t.Fatal(err)
	}

	checks := []resource.TestCheckFunc{
		// Need to check diff suppress funcs
		// resource.TestCheckResourceAttr(
		// 	"hpe_morpheus_app_blueprint_cloud_formation.example",
		// 	"blueprint_content",
		// 	"{\n \"AWSTemplateFormatVersion\" : \"2010-09-09\",\n\n \"Description\" : \"AWS CloudFormation Sample "+
		// 		"Template S3_Website_Bucket_With_Retain_On_Delete: Sample template showing how to create a publicly "+
		// 		"accessible S3 bucket configured for website access with a deletion policy of retain on delete. **WARNING** "+
		// 		"This template creates an S3 bucket that will NOT be deleted when the stack is deleted. You will be billed "+
		// 		"for the AWS resources used if you create a stack from this template.\",\n\n \"Resources\" : {\n "+
		// 		"\"S3Bucket\" : {\n \"Type\" : \"AWS::S3::Bucket\",\n \"Properties\" : {\n \"AccessControl\" : "+
		// 		"\"PublicRead\",\n \"WebsiteConfiguration\" : {\n \"IndexDocument\" : \"index.html\",\n \"ErrorDocument\" : "+
		// 		"\"error.html\"\n }\n },\n \"DeletionPolicy\" : \"Retain\"\n }\n },\n\n \"Outputs\" : {\n \"WebsiteURL\" : "+
		// 		"{\n \"Value\" : { \"Fn::GetAtt\" : [ \"S3Bucket\", \"WebsiteURL\" ] },\n \"Description\" : \"URL for "+
		// 		"website hosted on S3\"\n },\n \"S3BucketSecureURL\" : {\n \"Value\" : { \"Fn::Join\" : [ \"\", [ "+
		// 		"\"https://\", { \"Fn::GetAtt\" : [ \"S3Bucket\", \"DomainName\" ] } ] ] },\n \"Description\" : \"Name of "+
		// 		"S3 bucket to hold website content\"\n }\n }\n}",
		// ),

		resource.TestCheckResourceAttr(
			"hpe_morpheus_app_blueprint_cloud_formation.example",
			"capability_auto_expand",
			"true",
		),

		resource.TestCheckResourceAttr(
			"hpe_morpheus_app_blueprint_cloud_formation.example",
			"capability_iam",
			"true",
		),

		resource.TestCheckResourceAttr(
			"hpe_morpheus_app_blueprint_cloud_formation.example",
			"capability_named_iam",
			"true",
		),

		resource.TestCheckResourceAttr(
			"hpe_morpheus_app_blueprint_cloud_formation.example",
			"category",
			"cloudformation",
		),

		resource.TestCheckResourceAttr(
			"hpe_morpheus_app_blueprint_cloud_formation.example",
			"cloud_init_enabled",
			"true",
		),

		resource.TestCheckResourceAttr(
			"hpe_morpheus_app_blueprint_cloud_formation.example",
			"description",
			"Example cloud formation app blueprint",
		),

		resource.TestCheckResourceAttr(
			"hpe_morpheus_app_blueprint_cloud_formation.example",
			"install_agent",
			"true",
		),

		resource.TestCheckResourceAttr(
			"hpe_morpheus_app_blueprint_cloud_formation.example",
			"name",
			name,
		),

		resource.TestCheckResourceAttr(
			"hpe_morpheus_app_blueprint_cloud_formation.example",
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
				ExpectNonEmptyPlan: false,
				Check:              checkFn,
			},
			// Plan after apply
			{
				Config:             providerConfig + resourceConfig,
				ExpectNonEmptyPlan: false,
				PlanOnly:           true,
			},
		},
	})
}
