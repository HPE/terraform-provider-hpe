// (C) Copyright 2025-2026 Hewlett Packard Enterprise Development LP

package usergroup_test

import (
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/HPE/terraform-provider-hpe/morpheus"
	"github.com/HPE/terraform-provider-hpe/morpheus/framework/resources/role"
	"github.com/HPE/terraform-provider-hpe/morpheus/framework/resources/user"
	sdkv2morpheus "github.com/HPE/terraform-provider-hpe/morpheus/sdkv2"
	dsusergroup "github.com/HPE/terraform-provider-hpe/morpheus/sdkv2/datasources/usergroup"
	"github.com/HPE/terraform-provider-hpe/morpheus/sdkv2/resources/usergroup"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers/systemoverride"
)

func TestAccMorpheusDataSourceUserGroupsExampleOk(t *testing.T) {
	t.Parallel()

	defer testhelpers.RecordResult(t)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	testSystem := systemoverride.GetPreferred(t, "feature")
	providerConfig := testhelpers.ProviderBlock(testSystem)

	name := acctest.RandomWithPrefix(t.Name())

	var dependenciesConfig string

	if currentDependency, err := role.RenderRoleUserConfig(t, map[string]string{
		"Name": name,
	}); err != nil {
		t.Fatal(err)
	} else {
		dependenciesConfig += currentDependency
	}

	if currentDependency, err := user.RenderUserConfig(t, map[string]string{
		"Name":     name,
		"Username": strings.ToLower(name),
		"RoleIds":  "resource.hpe_morpheus_role.example.id",
	}); err != nil {
		t.Fatal(err)
	} else {
		dependenciesConfig += currentDependency
	}

	if currentDependency, err := usergroup.RenderUserGroupConfig(t, map[string]string{
		"Name":       name,
		"SudoAccess": "false",
		"UserIds":    "[resource.hpe_morpheus_user.example.id]",
	}); err != nil {
		t.Fatal(err)
	} else {
		dependenciesConfig += currentDependency
	}

	// find the user group we just created - filter on generated name
	datasourceConfig, err := dsusergroup.RenderUserGroupsConfig(t, map[string]string{
		"Values": `[resource.hpe_morpheus_user_group.example.name]`,
	})
	if err != nil {
		t.Fatal(err)
	}

	checks := []resource.TestCheckFunc{
		resource.TestCheckResourceAttrPair(
			"data.hpe_morpheus_user_groups.example",
			"ids.0",
			"hpe_morpheus_user_group.example",
			"id",
		),

		resource.TestCheckResourceAttr(
			"data.hpe_morpheus_user_groups.example",
			"filter.0.name",
			"name",
		),

		resource.TestCheckResourceAttr(
			"data.hpe_morpheus_user_groups.example",
			"sort_ascending",
			"true",
		),

		resource.TestCheckResourceAttr(
			"data.hpe_morpheus_user_groups.example",
			"filter.0.values.0",
			name,
		),
	}

	checkFn := resource.ComposeAggregateTestCheckFunc(checks...)
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, morpheus.New(), sdkv2morpheus.Provider()),
		Steps: []resource.TestStep{
			{
				Config:             providerConfig + dependenciesConfig + datasourceConfig,
				ExpectNonEmptyPlan: false,
				Check:              checkFn,
			},
		},
	})
}
