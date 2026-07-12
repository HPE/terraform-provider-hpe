// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package acctest

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/HPE/terraform-provider-hpe/opsramp/client"
	"github.com/HPE/terraform-provider-hpe/opsramp/utils/clientfactory"

	"github.com/HPE/terraform-provider-hpe/provider"
	"github.com/HPE/terraform-provider-hpe/provider/adapter"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
)

// ProtoV6ProviderFactories instantiates the provider for acceptance tests.
func ProtoV6ProviderFactories() map[string]func() (tfprotov6.ProviderServer, error) {
	return map[string]func() (tfprotov6.ProviderServer, error){
		"hpe": providerserver.NewProtocol6WithError(provider.New("test", adapter.NewOpsRamp())()),
	}
}

var envSuffix = []string{
	"client_id",
	"client_secret",
	"endpoint",
	"tenant",
}

// PreCheck validates required acceptance-test environment variables.
func PreCheck(t *testing.T) func() {
	return func() {
		t.Helper()

		for _, suffix := range envSuffix {
			if _, ok := LookupProviderEnv(suffix); !ok {
				t.Fatalf("TF_VAR_testacc_opsramp_%s not set for acceptance tests", suffix)
			}
		}
	}
}

// APIClient constructs an API client using TF_VAR naming convention
func APIClient(t *testing.T) (*client.OpsRampClient, error) {
	t.Helper()
	PreCheck(t)()

	client_id, ok := LookupProviderEnv("client_id")
	if !ok {
		return nil, fmt.Errorf("client_id not set for acceptance tests")
	}
	client_secret, ok := LookupProviderEnv("client_secret")
	if !ok {
		return nil, fmt.Errorf("client_secret not set for acceptance tests")
	}
	endpoint, ok := LookupProviderEnv("endpoint")
	if !ok {
		return nil, fmt.Errorf("endpoint not set for acceptance tests")
	}
	tenant, ok := LookupProviderEnv("tenant")
	if !ok {
		return nil, fmt.Errorf("tenant not set for acceptance tests")
	}

	clientfactory := clientfactory.NewClientFactory(
		client_id,
		client_secret,
		endpoint,
		tenant,
	)

	client, err := clientfactory.Client()
	if err != nil {
		return nil, fmt.Errorf("failed to create test client: %v", err)
	}

	return client, nil
}

// LookupProviderEnv looks up a TF_VAR provider input
func LookupProviderEnv(suffix string) (string, bool) {
	return os.LookupEnv("TF_VAR_testacc_opsramp_" + suffix)
}

// RandomName generates short, unique, test-safe names for acceptance resources.
func RandomName(prefix string) string {
	normalizedPrefix := strings.ToLower(strings.TrimSpace(prefix))
	normalizedPrefix = strings.ReplaceAll(normalizedPrefix, "_", "-")
	if normalizedPrefix == "" {
		normalizedPrefix = "opsramp"
	}

	return fmt.Sprintf("tfacc-%s-%s", normalizedPrefix, tfacctest.RandStringFromCharSet(6, tfacctest.CharSetAlphaNum))
}

// --- Scope detection and skip helpers ---
//
// Tests run with a single set of credentials (TF_VAR_testacc_opsramp_*).
// Scope is auto-detected from those credentials (MSP or CLIENT).
//
// To run all tests seamlessly:
//   Run 1: MSP credentials + target_client				→ MSP tests pass, CLIENT-scoped tests pass via override
//   Run 2: CLIENT credentials                          → CLIENT-scoped tests pass directly, MSP tests skip

const envTargetClientSuffix = "target_client"

// Scope returns the API client scope ("MSP" or "CLIENT") for the current credentials.
func Scope(t *testing.T) string {
	t.Helper()

	apiClient, err := APIClient(t)
	if err != nil {
		t.Fatalf("failed to determine scope: %v", err)
	}

	return strings.ToUpper(apiClient.Scope)
}

// SkipIfNotMSP skips the test if the credentials are not MSP-scoped.
// Use for resources that can only be managed at the MSP level (e.g., client resources).
func SkipIfNotMSP(t *testing.T) {
	t.Helper()

	if Scope(t) != "MSP" {
		t.Skip("skipping: test requires MSP-level credentials")
	}
}

// SkipIfNotClient skips the test if the credentials are not CLIENT-scoped.
// Use for resources that require direct CLIENT scope only.
func SkipIfNotClient(t *testing.T) {
	t.Helper()

	if Scope(t) != "CLIENT" {
		t.Skip("skipping: test requires CLIENT-level credentials")
	}
}

// RequireClientScope ensures the test can operate at CLIENT level.
//   - If credentials are CLIENT-scoped: returns "" (no override needed).
//   - If credentials are MSP-scoped: returns the OPSRAMP_ACC_TARGET_CLIENT value
//     (skips the test if that env var is not set).
//
// Use the returned value as the `client` attribute in resource HCL configs.
func RequireClientScope(t *testing.T) string {
	t.Helper()

	if Scope(t) == "CLIENT" {
		return ""
	}

	// MSP credentials — need a target client override
	return TargetClientID(t)
}

// TargetClientID returns the OPSRAMP_ACC_TARGET_CLIENT env var value.
// Skips the test if it is not set.
func TargetClientID(t *testing.T) string {
	t.Helper()

	id, ok := LookupProviderEnv(envTargetClientSuffix)

	if !ok {
		t.Fatalf("TF_VAR_testacc_opsramp_%s not set for acceptance tests", envTargetClientSuffix)
	}

	return id
}

const providerConfig = `
variable "testacc_opsramp_client_id" {
  default = null
}
variable "testacc_opsramp_client_secret" {
  default = null
}
variable "testacc_opsramp_endpoint" {
  default = null
}
variable "testacc_opsramp_tenant" {
  default = null
}

provider "hpe" {
  opsramp {
    client_id = var.testacc_opsramp_client_id
    client_secret = var.testacc_opsramp_client_secret
    endpoint      = var.testacc_opsramp_endpoint
    tenant        = var.testacc_opsramp_tenant
  }
}
`

// ProviderConfigHCL returns a provider block that can be used for acceptance testing.
func ProviderConfigHCL() string {
	return providerConfig
}
