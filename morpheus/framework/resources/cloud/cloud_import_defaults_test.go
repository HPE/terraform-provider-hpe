// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package cloud_test

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"

	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers/capabilities"
)

// TestAccMorpheusCloudResourceImportOutOfBandFillsDefaults proves that
// schemadefaults.Apply recovers a schema default the API omits, on import of a
// resource that was created outside Terraform.
//
// This is the scenario the null-in-state-after-import defect (MORPH-16192)
// actually bites: a resource created out of band via the UI or API, then
// imported. A resource created by Terraform round-trips its defaults because the
// plan resolves them and create sends them, so it cannot exercise the defect --
// which is why this test creates the cloud directly through the API, omitting
// config_hvm.certificate_provider (schema default "internal"). The Morpheus API
// omits certificateProvider from the GET response when it was never set, so on
// import the cloud Read maps it to null; schemadefaults.Apply fills it from the
// schema instead.
//
// With the fix, the imported config_hvm.certificate_provider is "internal".
// Without it, the value is empty and the ImportStateCheck below fails -- the
// test discriminates.
func TestAccMorpheusCloudResourceImportOutOfBandFillsDefaults(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.All)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	t.Parallel()

	name := acctest.RandomWithPrefix(t.Name())
	code := strings.ToLower(name)

	config := testhelpers.ProviderBlock() + `
resource "hpe_morpheus_cloud" "example" {
  name      = "` + name + `"
  tenant_id = 1
  group_id  = 1
  code      = "` + code + `"

  config_hvm = {
    enable_network_type_selection = false
  }
}`

	var cloudID int

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					cloudID = createCloudOmittingCertProvider(t, name, code)
					t.Cleanup(func() { deleteCloudOutOfBand(t, cloudID) })
				},
				Config:            config,
				ResourceName:      "hpe_morpheus_cloud.example",
				ImportState:       true,
				ImportStateIdFunc: func(*terraform.State) (string, error) { return strconv.Itoa(cloudID), nil },
				ImportStateCheck: func(states []*terraform.InstanceState) error {
					for _, is := range states {
						if is.Attributes["name"] != name {
							continue
						}

						got := is.Attributes["config_hvm.certificate_provider"]
						if got != "internal" {
							return fmt.Errorf(
								"config_hvm.certificate_provider = %q after import, want "+
									"\"internal\": the API omits it and schemadefaults.Apply "+
									"should have filled the schema default", got)
						}

						return nil
					}

					return fmt.Errorf("imported cloud %q not found in state", name)
				},
			},
		},
	})
}

// --- out-of-band API helpers (raw HTTP, so the resource is created the way the
// UI would, not through Terraform) ---

func featureHTTPClient() *http.Client {
	if os.Getenv("TF_VAR_testacc_morpheus_insecure") == "true" {
		return &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{InsecureSkipVerify: true}, //nolint:gosec // test appliance, self-signed
			},
		}
	}

	return &http.Client{}
}

func featureToken(t *testing.T, base string, hc *http.Client) string {
	t.Helper()

	form := url.Values{}
	form.Set("username", os.Getenv("TF_VAR_testacc_morpheus_username"))
	form.Set("password", os.Getenv("TF_VAR_testacc_morpheus_password"))

	req, err := http.NewRequest( //nolint:gosec // G704: URL is a controlled test-appliance env var
		http.MethodPost,
		base+"/oauth/token?grant_type=password&scope=write&client_id=morph-api",
		strings.NewReader(form.Encode()),
	)
	if err != nil {
		t.Fatalf("oauth request: %v", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := hc.Do(req) //nolint:gosec // G704: URL is a controlled test-appliance env var
	if err != nil {
		t.Fatalf("oauth: %v", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		t.Fatalf("oauth: status %d: %s", resp.StatusCode, string(body))
	}

	var out struct {
		AccessToken string `json:"access_token"`
	}
	if err := json.Unmarshal(body, &out); err != nil || out.AccessToken == "" {
		t.Fatalf("oauth: no token (status %d, err=%v): %s", resp.StatusCode, err, string(body))
	}

	return out.AccessToken
}

func createCloudOmittingCertProvider(t *testing.T, name, code string) int {
	t.Helper()

	base := strings.TrimRight(os.Getenv("TF_VAR_testacc_morpheus_url"), "/")
	hc := featureHTTPClient()
	token := featureToken(t, base, hc)

	// A standard cloud whose config omits certificateProvider. The API stores it
	// without that field and omits it from the GET response.
	body := `{"zone":{"name":"` + name + `","code":"` + code +
		`","zoneType":{"code":"standard"},"groupId":1,` +
		`"config":{"enableNetworkTypeSelection":"off"}}}`

	req, err := http.NewRequest( //nolint:gosec // G704: controlled test-appliance URL
		http.MethodPost,
		base+"/api/zones",
		strings.NewReader(body),
	)
	if err != nil {
		t.Fatalf("create cloud request: %v", err)
	}
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := hc.Do(req) //nolint:gosec // G704: URL is a controlled test-appliance env var
	if err != nil {
		t.Fatalf("create cloud: %v", err)
	}
	defer resp.Body.Close()

	b, _ := io.ReadAll(resp.Body)
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		t.Fatalf("create cloud: status %d: %s", resp.StatusCode, string(b))
	}

	var out struct {
		Success bool   `json:"success"`
		Msg     string `json:"msg"`
		Zone    struct {
			ID int `json:"id"`
		} `json:"zone"`
	}
	if err := json.Unmarshal(b, &out); err != nil {
		t.Fatalf("create cloud decode: %v (body: %s)", err, string(b))
	}
	if !out.Success || out.Zone.ID == 0 {
		t.Fatalf("create cloud failed: %s", out.Msg)
	}

	return out.Zone.ID
}

func deleteCloudOutOfBand(t *testing.T, id int) {
	t.Helper()

	if id == 0 {
		return
	}

	base := strings.TrimRight(os.Getenv("TF_VAR_testacc_morpheus_url"), "/")
	hc := featureHTTPClient()
	token := featureToken(t, base, hc)

	req, err := http.NewRequest( //nolint:gosec // G704: URL is a controlled test-appliance env var
		http.MethodDelete,
		base+"/api/zones/"+strconv.Itoa(id)+"?force=true",
		nil,
	)
	if err != nil {
		t.Logf("cleanup: delete request for cloud %d: %v", id, err)

		return
	}
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := hc.Do(req) //nolint:gosec // G704: URL is a controlled test-appliance env var
	if err != nil {
		t.Logf("cleanup: delete cloud %d: %v", id, err)

		return
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, _ := io.ReadAll(resp.Body)
		t.Logf("cleanup: delete cloud %d returned status %d: %s", id, resp.StatusCode, string(body))
	}
}
