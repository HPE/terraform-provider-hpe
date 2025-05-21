package group_test

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"regexp"
	"testing"

	"github.com/HewlettPackard/hpe-morpheus-go-sdk/sdk"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/HPE/terraform-provider-hpe/internal/provider"
	"github.com/HPE/terraform-provider-hpe/internal/subproviders/morpheus"
	"github.com/HPE/terraform-provider-hpe/internal/subproviders/morpheus/clientfactory"
	"github.com/HPE/terraform-provider-hpe/internal/subproviders/morpheus/datasources/group/consts"
)

const providerConfig = `
variable "testacc_morpheus_url" {}
variable "testacc_morpheus_insecure" {}
variable "testacc_morpheus_username" {}
variable "testacc_morpheus_password" {}

provider "hpe" {
  morpheus {
    url          = var.testacc_morpheus_url
    insecure     = var.testacc_morpheus_insecure
    username     = var.testacc_morpheus_username
    password     = var.testacc_morpheus_password
  }
}
`

const providerConfigOffline = `
provider "hpe" {
  morpheus {
    url          = ""
    username     = ""
    password     = ""
  }
}
`

var group sdk.ListGroups200ResponseAllOfGroupsInner

func newClient(t *testing.T, ctx context.Context) *sdk.APIClient {
	t.Helper()

	return clientfactory.NewAPIClient(
		ctx,
		os.Getenv("TF_VAR_testacc_morpheus_url"),
		os.Getenv("TF_VAR_testacc_morpheus_username"),
		os.Getenv("TF_VAR_testacc_morpheus_password"),
		"",
		clientfactory.WithInsecureTLS())
}

func createGroup(t *testing.T) sdk.ListGroups200ResponseAllOfGroupsInner {
	t.Helper()

	ctx := context.TODO()

	addGroup := sdk.NewAddGroupsRequestGroupWithDefaults()
	addGroup.SetName("testacc-" + t.Name())
	addGroup.SetCode("jason")
	addGroup.SetLocation("here")

	addGroupReq := sdk.NewAddGroupsRequest(*addGroup)

	client := newClient(t, ctx)

	g, hresp, err := client.GroupsAPI.AddGroups(ctx).AddGroupsRequest(*addGroupReq).Execute()
	if err != nil || hresp.StatusCode != http.StatusOK {
		t.Fatalf("POST failed for group %d - %v", 1, err)
	}

	group := g.GetGroup()

	t.Cleanup(func() {
		deleteGroup(t, group.GetId())
	})

	return group
}

func deleteGroup(t *testing.T, id int64) {
	t.Helper()

	ctx := context.TODO()

	client := newClient(t, ctx)

	_, hresp, err := client.GroupsAPI.RemoveGroups(ctx, id).Execute()
	if err != nil || hresp.StatusCode != http.StatusOK {
		t.Fatalf("POST failed for group %d - %v", 1, err)
	}
}

func newProviderWithError() (tfprotov6.ProviderServer, error) {
	providerInstance := provider.New("test", morpheus.New())()

	return providerserver.NewProtocol6WithError(providerInstance)()
}

var testAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
	"hpe": newProviderWithError,
}

func TestAccMorpheusFindGroupById(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	group := createGroup(t)
	groupID := fmt.Sprintf("%d", group.GetId())
	groupName := group.GetName()

	config := providerConfig + `
      data "hpe_morpheus_group" "test" {
        id = ` + groupID + `
      }`

	checks := []resource.TestCheckFunc{
		resource.TestCheckResourceAttr(
			"data.hpe_morpheus_group.test",
			"name",
			groupName,
		),
		resource.TestCheckResourceAttr(
			"data.hpe_morpheus_group.test",
			"id",
			groupID,
		),
	}

	checkFn := resource.ComposeAggregateTestCheckFunc(checks...)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check:  checkFn,
			},
		},
	})
}

func TestAccMorpheusFindGroupByName(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	group := createGroup(t)
	groupID := fmt.Sprintf("%d", group.GetId())
	groupName := group.GetName()

	config := providerConfig + `
      data "hpe_morpheus_group" "test" {
        name = "` + groupName + `" 
      }`

	checks := []resource.TestCheckFunc{
		resource.TestCheckResourceAttr(
			"data.hpe_morpheus_group.test",
			"name",
			groupName,
		),
		resource.TestCheckResourceAttr(
			"data.hpe_morpheus_group.test",
			"id",
			groupID,
		),
	}

	checkFn := resource.ComposeAggregateTestCheckFunc(checks...)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check:  checkFn,
			},
		},
	})
}

func TestAccMorpheusFindGroupNotFound(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	config := providerConfig + `
      data "hpe_morpheus_group" "test" {
        name = "______" 
      }`

	checks := []resource.TestCheckFunc{
		resource.TestCheckNoResourceAttr(
			"data.hpe_morpheus_group.test",
			"id",
		),
	}

	checkFn := resource.ComposeAggregateTestCheckFunc(checks...)

	expected := consts.ErrorNoGroupFound

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      config,
				Check:       checkFn,
				ExpectError: regexp.MustCompile(expected),
			},
		},
	})
}

func TestAccMorpheusFindGroupNoSearchAttrs(t *testing.T) {
	config := providerConfigOffline + `
      data "hpe_morpheus_group" "test" {
      }`

	checks := []resource.TestCheckFunc{
		resource.TestCheckNoResourceAttr(
			"data.hpe_morpheus_group.test",
			"id",
		),
	}

	checkFn := resource.ComposeAggregateTestCheckFunc(checks...)

	expected := consts.ErrorNoValidSearchTerms

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      config,
				Check:       checkFn,
				ExpectError: regexp.MustCompile(expected),
			},
		},
	})
}

func TestAccMorpheusFindGroupBothSearchAttrs(t *testing.T) {
	config := providerConfigOffline + `
      data "hpe_morpheus_group" "test" {
        id = 1
        name = "______" 
      }`

	checks := []resource.TestCheckFunc{
		resource.TestCheckNoResourceAttr(
			"data.hpe_morpheus_group.test",
			"id",
		),
	}

	checkFn := resource.ComposeAggregateTestCheckFunc(checks...)

	expected := consts.ErrorRunningPreApply

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      config,
				Check:       checkFn,
				ExpectError: regexp.MustCompile(expected),
			},
		},
	})
}
