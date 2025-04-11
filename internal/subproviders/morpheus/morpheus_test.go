package morpheus_test

import (
	"regexp"
	"testing"

	"github.com/HPE/terraform-provider-hpe/internal/provider"
	"github.com/HPE/terraform-provider-hpe/internal/subproviders/morpheus"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

var testAccProtoV6ProviderFactories = map[string]func() (
	tfprotov6.ProviderServer, error,
){
	"hpe": providerserver.NewProtocol6WithError(
		provider.New("test", morpheus.New())(),
	),
}

func TestAccMorpheusSubProviderDuplicateBlock(t *testing.T) {
	duplicateProviderConfig := `
provider "hpe" {
	morpheus {}
        # bad: duplicate provider block
	morpheus {}
}

# pseudo resource, needed to trigger parsing the provider block
resource "hpe_fake_resource" "test" {
}
`

	expected := "list must contain at least 0 elements and at most 1"
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				ExpectError:        regexp.MustCompile(expected),
				Config:             duplicateProviderConfig,
				PlanOnly:           true,
				ExpectNonEmptyPlan: false,
			},
		},
	})
}
