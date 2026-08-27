// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package acctest

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/HPE/terraform-provider-hpe/opsramp/client"
	"github.com/HPE/terraform-provider-hpe/opsramp/utils/clientfactory"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"

	"github.com/HPE/terraform-provider-hpe/provider"
	"github.com/HPE/terraform-provider-hpe/provider/adapter"
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
//
// A missing variable skips rather than fails. These are acceptance tests, so
// running without credentials is the normal case for anyone who is not
// targeting a live tenant, and failing would break `go test ./...` for them.
func PreCheck(t *testing.T) func() {
	return func() {
		t.Helper()

		for _, suffix := range envSuffix {
			if _, ok := LookupProviderEnv(suffix); !ok {
				t.Skipf(
					"skipping: TF_VAR_testacc_opsramp_%s not set for acceptance tests",
					suffix,
				)
			}
		}
	}
}

// APIClient constructs an API client using TF_VAR naming convention
func APIClient(t *testing.T) (*client.OpsRampClient, error) {
	t.Helper()
	PreCheck(t)()

	clientID, ok := LookupProviderEnv("client_id")
	if !ok {
		return nil, fmt.Errorf("client_id not set for acceptance tests")
	}
	clientSecret, ok := LookupProviderEnv("client_secret")
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
		clientID,
		clientSecret,
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
// Credential scope is auto-detected (MSP or CLIENT). The effective scope
// additionally considers whether target_client is set:
//
//   Credential scope  | target_client | Effective scope
//   MSP               | not set       | MSP
//   MSP               | set           | CLIENT (delegated)
//   CLIENT            | n/a           | CLIENT (direct)
//
// Three CI runs cover all paths:
//   Run 1: MSP credentials, no target_client       → MSP-level tests run
//   Run 2: MSP credentials + target_client         → Client tests via delegation (client attr set)
//   Run 3: CLIENT credentials                      → Client tests directly (no client attr)

const envTargetClientSuffix = "target_client"

// CredentialScope returns the raw API credential scope ("MSP" or "CLIENT").
func CredentialScope(t *testing.T) string {
	t.Helper()

	apiClient, err := APIClient(t)
	if err != nil {
		t.Fatalf("failed to determine scope: %v", err)
	}

	return strings.ToUpper(apiClient.Scope)
}

// EffectiveScope returns the testing scope taking target_client into account.
//   - MSP creds without target_client → "MSP"
//   - MSP creds with target_client    → "CLIENT"
//   - CLIENT creds                    → "CLIENT"
func EffectiveScope(t *testing.T) string {
	t.Helper()

	if CredentialScope(t) == "CLIENT" {
		return "CLIENT"
	}

	// MSP credentials — check if target_client makes this a client-level run
	if _, ok := LookupProviderEnv(envTargetClientSuffix); ok {
		return "CLIENT"
	}

	return "MSP"
}

// SkipIfNotMSP skips the test unless the effective scope is MSP.
// Use for resources that can only be managed at the MSP level (e.g., client resources).
func SkipIfNotMSP(t *testing.T) {
	t.Helper()

	if EffectiveScope(t) != "MSP" {
		t.Skip("skipping: test requires MSP-level scope (MSP creds without target_client)")
	}
}

// SkipIfNotClient skips the test unless the effective scope is CLIENT.
// Use for resources that require client-level scope.
func SkipIfNotClient(t *testing.T) {
	t.Helper()

	if EffectiveScope(t) != "CLIENT" {
		t.Skip("skipping: test requires CLIENT-level scope")
	}
}

// RequireClientScope ensures the test can operate at CLIENT level.
//   - CLIENT credentials: returns "" (no override needed).
//   - MSP credentials + target_client: returns the target client ID.
//   - MSP credentials without target_client: skips the test.
//
// Use the returned value as the `client` attribute in resource HCL configs.
func RequireClientScope(t *testing.T) string {
	t.Helper()

	if CredentialScope(t) == "CLIENT" {
		return ""
	}

	// MSP credentials — need a target client override
	return TargetClientID(t)
}

// TargetClientID returns the TF_VAR_testacc_opsramp_target_client env var value.
// Skips the test if it is not set.
func TargetClientID(t *testing.T) string {
	t.Helper()

	id, ok := LookupProviderEnv(envTargetClientSuffix)

	if !ok {
		t.Skipf("skipping: TF_VAR_testacc_opsramp_%s not set", envTargetClientSuffix)
	}

	return id
}

// OptionalClientOverride returns the client override for the current credentials.
// Unlike RequireClientScope it never skips:
//   - CLIENT credentials: returns "".
//   - MSP + target_client: returns the target client ID.
//   - MSP without target_client: returns "" (resource created at MSP level).
func OptionalClientOverride(t *testing.T) string {
	t.Helper()

	if CredentialScope(t) == "CLIENT" {
		return ""
	}

	if id, ok := LookupProviderEnv(envTargetClientSuffix); ok {
		return id
	}

	return ""
}

// ClientAttrHCL returns a Terraform HCL snippet for the `client` attribute.
// Returns an empty string when clientOverride is empty (no attribute emitted).
func ClientAttrHCL(clientOverride string) string {
	if clientOverride != "" {
		return fmt.Sprintf(`client = "%s"`, clientOverride)
	}

	return ""
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
