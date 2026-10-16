// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package cloud_test

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"

	"github.com/HPE/terraform-provider-hpe/morpheus/framework/resources/cloud"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers/capabilities"
)

// rotateZoneConfigSecret rewrites one config secret of a cloud directly through
// the API (PUT /api/zones/{id}), reproducing an out-of-band change made in the
// UI or by another tool. It is used in a PreConfig so the following plan sees
// appliance-side drift. The hash Morpheus renders beside the secret changes as
// a result, which is what the provider's drift detection compares.
func rotateZoneConfigSecret(t *testing.T, id int64, configKey, value string) {
	t.Helper()

	// GET the zone's existing config first and PUT the whole map back with only
	// the secret key changed. A PUT carrying just the one key can replace the
	// config map and wipe the cloud's other config values (which would make the
	// following plan show unrelated drift); preserving the full map avoids that.
	cfg := zoneConfig(t, id)
	cfg[configKey] = value

	base, _ := testhelpers.LookupProviderEnv("", "url")
	token, _ := testhelpers.LookupProviderEnv("", "access_token")

	body, err := json.Marshal(map[string]any{
		"zone": map[string]any{"config": cfg},
	})
	if err != nil {
		t.Fatalf("rotateZoneConfigSecret: marshal: %v", err)
	}

	url := base + "/api/zones/" + strconv.FormatInt(id, 10)

	req, err := http.NewRequestWithContext(context.Background(), http.MethodPut, url, bytes.NewReader(body))
	if err != nil {
		t.Fatalf("rotateZoneConfigSecret: request: %v", err)
	}

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("rotateZoneConfigSecret: PUT %s: %v", url, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(resp.Body)
		t.Fatalf("rotateZoneConfigSecret: PUT %s returned %d: %s", url, resp.StatusCode, string(b))
	}
}

// zoneConfig reads a cloud's config map through the API, so the rotate PUT can
// preserve every key but the one it changes.
func zoneConfig(t *testing.T, id int64) map[string]any {
	t.Helper()

	base, ok := testhelpers.LookupProviderEnv("", "url")
	if !ok {
		t.Fatal("TF_VAR_testacc_morpheus_url not set")
	}

	token, ok := testhelpers.LookupProviderEnv("", "access_token")
	if !ok {
		t.Fatal("TF_VAR_testacc_morpheus_access_token not set")
	}

	url := base + "/api/zones/" + strconv.FormatInt(id, 10)

	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, url, nil)
	if err != nil {
		t.Fatalf("zoneConfig: request: %v", err)
	}

	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("zoneConfig: GET %s: %v", url, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(resp.Body)
		t.Fatalf("zoneConfig: GET %s returned %d: %s", url, resp.StatusCode, string(b))
	}

	var out struct {
		Zone struct {
			Config map[string]any `json:"config"`
		} `json:"zone"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		t.Fatalf("zoneConfig: decode: %v", err)
	}

	if out.Zone.Config == nil {
		return map[string]any{}
	}

	return out.Zone.Config
}

// zoneConfigHashValue reads one hash key from a cloud's config map, so a test
// can prove the rotate PUT actually changed the appliance hash.
func zoneConfigHashValue(t *testing.T, id int64, hashKey string) string {
	t.Helper()

	cfg := zoneConfig(t, id)
	if v, ok := cfg[hashKey].(string); ok {
		return v
	}

	return ""
}

// TestAccMorpheusCloudResourceWriteOnlyDriftDetected asserts appliance-side
// drift detection for the VMware password: after the cloud is created, the
// config password is rotated through the API, so the next plan must see the
// drift and plan an update that re-sends the configured value; the apply then
// restores it and the following plan is clean. Needs a reachable vCenter
// (Morpheus validates the connection on create even for a disabled cloud).
func TestAccMorpheusCloudResourceWriteOnlyDriftDetected(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.VMware)
	testhelpers.SkipUnlessAcc(t)

	apiURL, ok := testhelpers.LookupProviderEnv("", "vcenter_url")
	if !ok {
		t.Skip("TF_VAR_testacc_morpheus_vcenter_url (and _vcenter_username/_vcenter_datacenter) not set")
	}

	vcUser, _ := testhelpers.LookupProviderEnv("", "vcenter_username")
	vcDC, _ := testhelpers.LookupProviderEnv("", "vcenter_datacenter")

	name := acctest.RandomWithPrefix(t.Name())
	visibility := testhelpers.TenantVisibility(t)

	const password = "Secret123!"

	config := testhelpers.ProviderBlock() + testhelpers.WhoamiBlock() + `
resource "hpe_morpheus_cloud" "wo" {
  name       = "` + name + `"
  code       = "` + name + `"
  tenant_id  = ` + testhelpers.WhoamiTenantIDRef + `
  group_id   = 1
  enabled    = false
  location   = "somewhere"
  visibility = "` + visibility + `"

  config_vmware = {
    api_url    = "` + apiURL + `"
    username   = "` + vcUser + `"
    password   = "` + password + `"
    datacenter = "` + vcDC + `"
  }
}
`

	var stashed *terraform.State

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check:  testhelpers.StashState(&stashed),
			},
			{
				// Rotate the secret on the appliance, then plan: drift must be
				// detected and an update (re-send) planned.
				PreConfig: func() {
					id := idFromState(t, stashed, "hpe_morpheus_cloud.wo")
					before := zoneConfigHashValue(t, id, "passwordHash")
					rotateZoneConfigSecret(t, id, "password", "Rotated456!")
					if after := zoneConfigHashValue(t, id, "passwordHash"); after == before {
						t.Fatalf("rotate PUT did not change passwordHash (%q); "+
							"drift cannot be detected", before)
					}
				},
				Config:             config,
				PlanOnly:           true,
				ExpectNonEmptyPlan: true,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("hpe_morpheus_cloud.wo", plancheck.ResourceActionUpdate),
					},
				},
			},
			{
				// Apply the re-send.
				Config: config,
			},
			{
				// The configured value has been restored: the plan is clean.
				Config:   config,
				PlanOnly: true,
			},
		},
	})
}

// TestAccMorpheusCloudResourceAzureWriteOnlyDriftDetected asserts appliance-side
// drift detection for the Azure client secret: after the cloud is created, the
// config clientSecret is rotated through the API, so the next plan must see the
// drift and plan an update that re-sends the configured value; the apply then
// restores it and the following plan is clean. A global Azure cloud is accepted
// without a reachable endpoint, so no extra fixture env is needed beyond the
// Azure capability.
func TestAccMorpheusCloudResourceAzureWriteOnlyDriftDetected(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.Azure)
	testhelpers.SkipUnlessAcc(t)

	name := acctest.RandomWithPrefix(t.Name())
	code := strings.ToLower(name)

	resourceConfig, err := cloud.RenderCloudAzureConfig(t, map[string]string{
		"Name":          name,
		"Code":          code,
		"TenantId":      testhelpers.WhoamiTenantIDRef,
		"GroupId":       "1",
		"Label":         "aLabel",
		"ApplianceUrl":  "https://somewhere.com",
		"AzureRegion":   "eastus",
		"SubscriberId":  "sub-12345-67890",
		"AzureTenantId": "tenant-12345-67890",
		"ClientId":      "client-abc-123",
		"ClientSecret":  "example",
		"ResourceGroup": "my-resource-group",
	})
	if err != nil {
		t.Fatal(err)
	}

	config := testhelpers.ProviderBlock() + testhelpers.WhoamiBlock() + resourceConfig

	var stashed *terraform.State

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check:  testhelpers.StashState(&stashed),
			},
			{
				PreConfig: func() {
					id := idFromState(t, stashed, "hpe_morpheus_cloud.example")
					before := zoneConfigHashValue(t, id, "clientSecretHash")
					rotateZoneConfigSecret(t, id, "clientSecret", "rotated-secret-000")
					if after := zoneConfigHashValue(t, id, "clientSecretHash"); after == before {
						t.Fatalf("rotate PUT did not change clientSecretHash (%q); "+
							"drift cannot be detected", before)
					}
				},
				Config:             config,
				PlanOnly:           true,
				ExpectNonEmptyPlan: true,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("hpe_morpheus_cloud.example", plancheck.ResourceActionUpdate),
					},
				},
			},
			{
				Config: config,
			},
			{
				Config:   config,
				PlanOnly: true,
			},
		},
	})
}

// idFromState reads the int64 "id" attribute of a resource from stashed state.
func idFromState(t *testing.T, state *terraform.State, address string) int64 {
	t.Helper()

	if state == nil {
		t.Fatal("no state stashed; add testhelpers.StashState to the preceding step")
	}

	rs, ok := state.RootModule().Resources[address]
	if !ok {
		t.Fatalf("resource not found in state: %s", address)
	}

	id, err := strconv.ParseInt(rs.Primary.Attributes["id"], 10, 64)
	if err != nil {
		t.Fatalf("resource %s id %q is not an integer: %v", address, rs.Primary.Attributes["id"], err)
	}

	return id
}
