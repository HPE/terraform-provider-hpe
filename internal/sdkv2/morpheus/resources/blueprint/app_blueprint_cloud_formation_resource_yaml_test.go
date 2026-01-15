// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package blueprint_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/HPE/terraform-provider-hpe/internal/framework/subproviders/morpheus"
	"github.com/HPE/terraform-provider-hpe/internal/framework/subproviders/morpheus/testhelpers"
	sdkv2morpheus "github.com/HPE/terraform-provider-hpe/internal/sdkv2/morpheus"
	"github.com/HPE/terraform-provider-hpe/internal/sdkv2/morpheus/resources/blueprint"
)

func TestAccMorpheusAppBlueprintCloudFormationYamlExampleOk(t *testing.T) {
	t.Parallel()

	defer testhelpers.RecordResult(t)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	t.Skip("Skipping due to missing infrastructure in test environment")

	providerConfig := testhelpers.ProviderBlock()

	name := acctest.RandomWithPrefix(t.Name())

	resourceConfig, err := blueprint.RenderAppBlueprintCloudFormationYamlConfig(t, map[string]string{
		"Name": name,
	})
	if err != nil {
		t.Fatal(err)
	}

	checks := []resource.TestCheckFunc{
		resource.TestCheckResourceAttr(
			"hpe_morpheus_app_blueprint_cloud_formation.example",
			"blueprint_content",
			"---\nAWSTemplateFormatVersion: '2010-09-09'\nDescription: 'AWS CloudFormation Sample Template "+
				"S3_Website_Bucket_With_Retain_On_Delete:\n Sample template showing how to create a publicly "+
				"accessible S3 bucket configured\n for website access with a deletion policy of retain on delete. "+
				"**WARNING** This\n template creates an S3 bucket that will NOT be deleted when the stack is deleted.\n "+
				"You will be billed for the AWS resources used if you create a stack from this template.'\nResources:\n "+
				"S3Bucket:\n Type: AWS::S3::Bucket\n Properties:\n AccessControl: PublicRead\n WebsiteConfiguration:\n "+
				"IndexDocument: index.html\n ErrorDocument: error.html\n DeletionPolicy: Retain\nOutputs:\n WebsiteURL:\n "+
				"Value:\n Fn::GetAtt:\n - S3Bucket\n - WebsiteURL\n Description: URL for website hosted on S3\n "+
				"S3BucketSecureURL:\n Value:\n Fn::Join:\n - ''\n - - https://\n - Fn::GetAtt:\n - S3Bucket\n - "+
				"DomainName\n Description: Name of S3 bucket to hold website content",
		),

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
			"example_cloud_formation_app_blueprint_yaml",
		),

		resource.TestCheckResourceAttr(
			"hpe_morpheus_app_blueprint_cloud_formation.example",
			"source_type",
			"yaml",
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
