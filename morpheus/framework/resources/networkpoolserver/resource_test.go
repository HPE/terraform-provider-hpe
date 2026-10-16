package networkpoolserver_test

import (
	"os"
	"regexp"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"

	"github.com/HPE/terraform-provider-hpe/morpheus/framework/resources/networkpoolserver"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers/capabilities"
	"github.com/HPE/terraform-provider-hpe/provider/adapter"
)

func TestMain(m *testing.M) {
	code := testhelpers.TestMain(m)
	testhelpers.WriteMergedResults()
	os.Exit(code)
}

func TestAccMorpheusNetworkPoolServerResourceExampleOk(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.NetworkPool)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	t.Parallel()

	providerConfig := testhelpers.ProviderBlock()
	name := acctest.RandomWithPrefix(t.Name())
	resourceName := "hpe_morpheus_network_pool_server.example"

	resourceConfig, err := networkpoolserver.RenderNetworkPoolServerInfobloxConfig(t, map[string]string{
		"Name": name,
	})
	if err != nil {
		t.Fatal(err)
	}

	resourceConfig = strings.Replace(
		resourceConfig,
		`resource "hpe_morpheus_network_pool_server" "infoblox" {`,
		`resource "hpe_morpheus_network_pool_server" "example" {`,
		1,
	)

	checks := resource.ComposeAggregateTestCheckFunc(
		resource.TestCheckResourceAttr(resourceName, "name", name),
		resource.TestCheckResourceAttr(resourceName, "type_code", "infoblox"),
		resource.TestCheckResourceAttrSet(resourceName, "type_id"),
		resource.TestCheckResourceAttr(resourceName, "enabled", "true"),
		resource.TestCheckResourceAttr(resourceName, "service_url", "https://infoblox.example.com/wapi/v2.12"),
		resource.TestCheckResourceAttr(resourceName, "service_username", "admin"),
		resource.TestCheckResourceAttr(resourceName, "ignore_ssl", "true"),
		resource.TestCheckResourceAttr(resourceName, "network_filter", "10.0.0.0/8"),
		resource.TestCheckResourceAttr(resourceName, "zone_filter", "example.com"),
		resource.TestCheckResourceAttr(resourceName, "tenant_match", ".*"),
		resource.TestCheckResourceAttr(resourceName, "service_mode", "static"),
		resource.TestCheckResourceAttr(resourceName, "service_throttle_rate", "0"),
		resource.TestCheckResourceAttrSet(resourceName, "id"),
	)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, adapter.NewMorpheus(), nil),
		Steps: []resource.TestStep{
			{
				Config: providerConfig + resourceConfig,
				Check:  checks,
			},
			{
				Config:             providerConfig + resourceConfig,
				ExpectNonEmptyPlan: false,
				PlanOnly:           true,
			},
			{
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"service_username", "service_password_wo", "service_password_wo_version", "write_only_revision",
				},
				ResourceName: "hpe_morpheus_network_pool_server.example",
			},
		},
	})
}

func TestAccMorpheusNetworkPoolServerResourceUpdateOk(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.NetworkPool)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	t.Parallel()

	providerConfig := testhelpers.ProviderBlock()
	name := acctest.RandomWithPrefix(t.Name())
	resourceName := "hpe_morpheus_network_pool_server.example"

	createConfig, err := networkpoolserver.RenderNetworkPoolServerInfobloxConfig(t, map[string]string{
		"Name": name,
	})
	if err != nil {
		t.Fatal(err)
	}

	createConfig = strings.Replace(
		createConfig,
		`resource "hpe_morpheus_network_pool_server" "infoblox" {`,
		`resource "hpe_morpheus_network_pool_server" "example" {`,
		1,
	)

	updateConfig := `
resource "hpe_morpheus_network_pool_server" "example" {
  name                        = "` + name + `"
  type_code                   = "infoblox"
  enabled                     = false
  service_url                 = "https://infoblox.example.com/wapi/v2.12"
  service_username            = "admin"
  service_password_wo         = "changeme"
  service_password_wo_version = 1
  ignore_ssl                  = true
  network_filter              = "192.168.0.0/16"
  zone_filter                 = "example.com"
  tenant_match                = ".*"
  service_mode                = "static"
  service_throttle_rate       = 25
}
`

	createChecks := resource.ComposeAggregateTestCheckFunc(
		resource.TestCheckResourceAttr(resourceName, "name", name),
		resource.TestCheckResourceAttr(resourceName, "type_code", "infoblox"),
		resource.TestCheckResourceAttr(resourceName, "enabled", "true"),
		resource.TestCheckResourceAttr(resourceName, "service_url", "https://infoblox.example.com/wapi/v2.12"),
		resource.TestCheckResourceAttr(resourceName, "service_username", "admin"),
		resource.TestCheckResourceAttr(resourceName, "ignore_ssl", "true"),
		resource.TestCheckResourceAttr(resourceName, "network_filter", "10.0.0.0/8"),
		resource.TestCheckResourceAttr(resourceName, "zone_filter", "example.com"),
		resource.TestCheckResourceAttr(resourceName, "tenant_match", ".*"),
		resource.TestCheckResourceAttr(resourceName, "service_mode", "static"),
		resource.TestCheckResourceAttr(resourceName, "service_throttle_rate", "0"),
	)

	updateChecks := resource.ComposeAggregateTestCheckFunc(
		resource.TestCheckResourceAttr(resourceName, "name", name),
		resource.TestCheckResourceAttr(resourceName, "type_code", "infoblox"),
		resource.TestCheckResourceAttr(resourceName, "enabled", "false"),
		resource.TestCheckResourceAttr(resourceName, "service_url", "https://infoblox.example.com/wapi/v2.12"),
		resource.TestCheckResourceAttr(resourceName, "service_username", "admin"),
		resource.TestCheckResourceAttr(resourceName, "ignore_ssl", "true"),
		resource.TestCheckResourceAttr(resourceName, "network_filter", "192.168.0.0/16"),
		resource.TestCheckResourceAttr(resourceName, "zone_filter", "example.com"),
		resource.TestCheckResourceAttr(resourceName, "tenant_match", ".*"),
		resource.TestCheckResourceAttr(resourceName, "service_mode", "static"),
		resource.TestCheckResourceAttr(resourceName, "service_throttle_rate", "25"),
	)

	checkInPlaceUpdate := resource.ConfigPlanChecks{
		PreApply: []plancheck.PlanCheck{
			plancheck.ExpectResourceAction(resourceName, plancheck.ResourceActionUpdate),
		},
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, adapter.NewMorpheus(), nil),
		Steps: []resource.TestStep{
			{Config: providerConfig + createConfig, Check: createChecks},
			{Config: providerConfig + updateConfig, Check: updateChecks, ConfigPlanChecks: checkInPlaceUpdate},
			{Config: providerConfig + updateConfig, ExpectNonEmptyPlan: false, PlanOnly: true},
		},
	})
}

// TestAccMorpheusNetworkPoolServerResourcePhpIpamOk creates a phpIPAM pool
// server against a real phpIPAM service, which validates the app ID (the
// create is refused with "App ID is required" without config.appId). It runs
// only when the service is configured through the TF_ACC_PHPIPAM_* variables.
func TestAccMorpheusNetworkPoolServerResourcePhpIpamOk(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.NetworkPool)

	serviceURL := os.Getenv("TF_ACC_PHPIPAM_URL")
	appID := os.Getenv("TF_ACC_PHPIPAM_APP_ID")
	username := os.Getenv("TF_ACC_PHPIPAM_USERNAME")
	password := os.Getenv("TF_ACC_PHPIPAM_PASSWORD")
	if serviceURL == "" || appID == "" || username == "" || password == "" {
		t.Skip("phpIPAM service not configured; set TF_ACC_PHPIPAM_URL, TF_ACC_PHPIPAM_APP_ID, " +
			"TF_ACC_PHPIPAM_USERNAME and TF_ACC_PHPIPAM_PASSWORD")
	}

	t.Parallel()

	providerConfig := testhelpers.ProviderBlock()
	name := acctest.RandomWithPrefix(t.Name())
	resourceName := "hpe_morpheus_network_pool_server.phpipam"

	// Render every configuration up front: TestStep may evaluate Config off
	// the test goroutine, where t.Fatal must not be called.
	render := func(name, appID string) string {
		cfg, err := networkpoolserver.RenderNetworkPoolServerPhpIPAMConfig(t, map[string]string{
			"Name":            name,
			"ServiceUrl":      serviceURL,
			"AppId":           appID,
			"ServiceUsername": username,
			"ServicePassword": password,
			"IgnoreSsl":       "true",
		})
		if err != nil {
			t.Fatal(err)
		}

		return providerConfig + cfg
	}
	noAppIDConfig := render(name, "")
	createConfig := render(name, appID)
	renamedConfig := render(name+"-renamed", appID)
	appID2 := os.Getenv("TF_ACC_PHPIPAM_APP_ID_2")
	var appID2Config string
	if appID2 != "" {
		appID2Config = render(name+"-renamed", appID2)
	}

	steps := []resource.TestStep{
		{
			// An empty app_id is refused by the appliance ("App ID is required").
			// The provider itself refuses a missing app_id for phpIPAM at plan
			// time (see validate.go), covered by the unit tests.
			Config:      noAppIDConfig,
			ExpectError: regexp.MustCompile(`(?i)app ?id`),
		},
		{
			Config: createConfig,
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr(resourceName, "name", name),
				resource.TestCheckResourceAttr(resourceName, "type_code", "phpipam"),
				resource.TestCheckResourceAttr(resourceName, "service_url", serviceURL),
				resource.TestCheckResourceAttr(resourceName, "app_id", appID),
				resource.TestCheckResourceAttrSet(resourceName, "id"),
			),
		},
		{
			Config:   createConfig,
			PlanOnly: true,
		},
		{
			// Update sends config.appId alongside the change; the appliance
			// re-validates it against phpIPAM.
			Config: renamedConfig,
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr(resourceName, "name", name+"-renamed"),
				resource.TestCheckResourceAttr(resourceName, "app_id", appID),
			),
		},
	}

	// Changing the app ID itself needs a second application on the phpIPAM
	// service, which not every lab has.
	if appID2 != "" {
		steps = append(steps, resource.TestStep{
			Config: appID2Config,
			Check:  resource.TestCheckResourceAttr(resourceName, "app_id", appID2),
		})
	}

	steps = append(steps, resource.TestStep{
		ImportState:       true,
		ImportStateVerify: true,
		ImportStateVerifyIgnore: []string{
			"service_username", "service_password_wo", "service_password_wo_version", "write_only_revision",
			// The status moves as the appliance syncs the real service.
			"status",
		},
		ResourceName: resourceName,
	})

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, adapter.NewMorpheus(), nil),
		Steps:                    steps,
	})
}
