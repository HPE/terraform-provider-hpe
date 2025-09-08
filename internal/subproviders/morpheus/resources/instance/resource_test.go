// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

//go:generate go run ../../../../../cmd/render resources/hpe_morpheus_instance example.tf.tmpl Name "TestInstance" InstanceType "9" ResourcePool "pool-62299"

package instance_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/HPE/terraform-provider-hpe/internal/provider"
	"github.com/HPE/terraform-provider-hpe/internal/subproviders/morpheus"
	"github.com/HPE/terraform-provider-hpe/internal/subproviders/morpheus/testhelpers"
)

func newProviderWithError() (tfprotov6.ProviderServer, error) {
	providerInstance := provider.New("test", morpheus.New())()

	return providerserver.NewProtocol6WithError(providerInstance)()
}

var testAccProtoV6ProviderFactories = map[string]func() (
	tfprotov6.ProviderServer, error,
){
	"hpe": newProviderWithError,
}

// Tests that our example file template used for docs is a valid config
func TestAccMorpheusInstanceExampleOk(t *testing.T) {
	defer testhelpers.RecordResult(t)
	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	providerConfig := testhelpers.ProviderBlock()

	name := acctest.RandomWithPrefix(t.Name())
	instanceTypeID := "9"
	resourcePool := "pool-62299"

	resourceConfig, err := testhelpers.RenderExample(t, "example.tf.tmpl",
		"Name", name,
		"InstanceType", instanceTypeID,
		"ResourcePool", resourcePool,
	)
	if err != nil {
		t.Fatal(err)
	}

	checks := []resource.TestCheckFunc{
		resource.TestCheckResourceAttr(
			"hpe_morpheus_instance.example",
			"name",
			name,
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_instance.example",
			"instance_type_id",
			instanceTypeID,
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_instance.example",
			"config.resourcePoolId",
			resourcePool,
		),
	}

	checkFn := resource.ComposeAggregateTestCheckFunc(checks...)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:             providerConfig + resourceConfig,
				ExpectNonEmptyPlan: false,
				PlanOnly:           false,
				Check:              checkFn,
			},
		},
	})
}
