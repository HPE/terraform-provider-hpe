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

	// Morpheus 9.1.2 verifies the service URL on create, so a reachable IPAM
	// is required; without the fixture there is no appliance to create against.
	f := testhelpers.RequireIpamFixture(t)

	t.Parallel()

	providerConfig := testhelpers.ProviderBlock()
	name := acctest.RandomWithPrefix(t.Name())
	resourceName := "hpe_morpheus_network_pool_server.example"

	resourceConfig, err := networkpoolserver.RenderNetworkPoolServerInfobloxConfig(t, map[string]string{
		"Name":            name,
		"TypeCode":        f.TypeCode,
		"ServiceUrl":      f.URL,
		"ServiceUsername": f.Username,
		"ServicePassword": f.Password,
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

	// zone_filter, tenant_match and service_mode are Infoblox-specific; the
	// shared Infoblox template always emits them, but a non-Infoblox fixture
	// must not be sent attributes its type does not accept, so strip them
	// unless the fixture is an Infoblox server.
	if f.TypeCode != "infoblox" {
		resourceConfig = stripInfobloxOnlyAttrs(resourceConfig)
	}
	resourceConfig = f.WithAppID(resourceConfig)

	checkList := []resource.TestCheckFunc{
		resource.TestCheckResourceAttr(resourceName, "name", name),
		resource.TestCheckResourceAttr(resourceName, "type_code", f.TypeCode),
		resource.TestCheckResourceAttrSet(resourceName, "type_id"),
		resource.TestCheckResourceAttr(resourceName, "enabled", "true"),
		resource.TestCheckResourceAttr(resourceName, "service_url", f.URL),
		resource.TestCheckResourceAttr(resourceName, "service_username", f.Username),
		resource.TestCheckResourceAttr(resourceName, "ignore_ssl", "true"),
		resource.TestCheckResourceAttr(resourceName, "network_filter", "10.0.0.0/8"),
		resource.TestCheckResourceAttr(resourceName, "service_throttle_rate", "0"),
		resource.TestCheckResourceAttrSet(resourceName, "id"),
	}

	// zone_filter, tenant_match and service_mode are Infoblox-specific
	// attributes (part of the Infoblox pool server schema); the example
	// template always emits them, but only Infoblox echoes them back, so
	// assert them only when the fixture is an Infoblox server.
	if f.TypeCode == "infoblox" {
		checkList = append(checkList,
			resource.TestCheckResourceAttr(resourceName, "zone_filter", "example.com"),
			resource.TestCheckResourceAttr(resourceName, "tenant_match", ".*"),
			resource.TestCheckResourceAttr(resourceName, "service_mode", "static"),
		)
	}

	checks := resource.ComposeAggregateTestCheckFunc(checkList...)

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

	// Morpheus 9.1.2 verifies the service URL on create, so a reachable IPAM
	// is required; without the fixture there is no appliance to create against.
	f := testhelpers.RequireIpamFixture(t)

	t.Parallel()

	providerConfig := testhelpers.ProviderBlock()
	name := acctest.RandomWithPrefix(t.Name())
	resourceName := "hpe_morpheus_network_pool_server.example"

	createConfig, err := networkpoolserver.RenderNetworkPoolServerInfobloxConfig(t, map[string]string{
		"Name":            name,
		"TypeCode":        f.TypeCode,
		"ServiceUrl":      f.URL,
		"ServiceUsername": f.Username,
		"ServicePassword": f.Password,
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

	// zone_filter, tenant_match and service_mode are Infoblox-specific; strip
	// them from the create config unless the fixture is an Infoblox server so a
	// non-Infoblox type is not sent attributes it does not accept.
	if f.TypeCode != "infoblox" {
		createConfig = stripInfobloxOnlyAttrs(createConfig)
	}
	createConfig = f.WithAppID(createConfig)

	// zone_filter, tenant_match and service_mode are Infoblox-specific; include
	// them in the update config only for an Infoblox fixture so a non-Infoblox
	// type is not sent attributes it does not accept.
	infobloxExtra := ""
	if f.TypeCode == "infoblox" {
		infobloxExtra = `
  zone_filter                 = "example.com"
  tenant_match                = ".*"
  service_mode                = "static"`
	}

	updateConfig := `
resource "hpe_morpheus_network_pool_server" "example" {
  name                        = "` + name + `"
  type_code                   = "` + f.TypeCode + `"
  enabled                     = false
  service_url                 = "` + f.URL + `"
  service_username            = "` + f.Username + `"
  service_password_wo         = "` + f.Password + `"
  service_password_wo_version = 1
  ignore_ssl                  = true
  network_filter              = "192.168.0.0/16"` + infobloxExtra + `
  service_throttle_rate       = 25
}
`
	updateConfig = f.WithAppID(updateConfig)

	createCheckList := []resource.TestCheckFunc{
		resource.TestCheckResourceAttr(resourceName, "name", name),
		resource.TestCheckResourceAttr(resourceName, "type_code", f.TypeCode),
		resource.TestCheckResourceAttr(resourceName, "enabled", "true"),
		resource.TestCheckResourceAttr(resourceName, "service_url", f.URL),
		resource.TestCheckResourceAttr(resourceName, "service_username", f.Username),
		resource.TestCheckResourceAttr(resourceName, "ignore_ssl", "true"),
		resource.TestCheckResourceAttr(resourceName, "network_filter", "10.0.0.0/8"),
		resource.TestCheckResourceAttr(resourceName, "service_throttle_rate", "0"),
	}

	updateCheckList := []resource.TestCheckFunc{
		resource.TestCheckResourceAttr(resourceName, "name", name),
		resource.TestCheckResourceAttr(resourceName, "type_code", f.TypeCode),
		resource.TestCheckResourceAttr(resourceName, "enabled", "false"),
		resource.TestCheckResourceAttr(resourceName, "service_url", f.URL),
		resource.TestCheckResourceAttr(resourceName, "service_username", f.Username),
		resource.TestCheckResourceAttr(resourceName, "ignore_ssl", "true"),
		resource.TestCheckResourceAttr(resourceName, "network_filter", "192.168.0.0/16"),
		resource.TestCheckResourceAttr(resourceName, "service_throttle_rate", "25"),
	}

	if f.TypeCode == "infoblox" {
		infobloxChecks := []resource.TestCheckFunc{
			resource.TestCheckResourceAttr(resourceName, "zone_filter", "example.com"),
			resource.TestCheckResourceAttr(resourceName, "tenant_match", ".*"),
			resource.TestCheckResourceAttr(resourceName, "service_mode", "static"),
		}
		createCheckList = append(createCheckList, infobloxChecks...)
		updateCheckList = append(updateCheckList, infobloxChecks...)
	}

	createChecks := resource.ComposeAggregateTestCheckFunc(createCheckList...)
	updateChecks := resource.ComposeAggregateTestCheckFunc(updateCheckList...)

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

// stripInfobloxOnlyAttrs removes the Infoblox-specific HCL lines (zone_filter,
// tenant_match, service_mode) from a rendered network_pool_server config. The
// shared Infoblox example template always emits them; a non-Infoblox fixture
// must not be sent attributes its type does not accept.
func stripInfobloxOnlyAttrs(config string) string {
	var kept []string
	for _, line := range strings.Split(config, "\n") {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "zone_filter") ||
			strings.HasPrefix(trimmed, "tenant_match") ||
			strings.HasPrefix(trimmed, "service_mode") {
			continue
		}
		kept = append(kept, line)
	}

	return strings.Join(kept, "\n")
}
