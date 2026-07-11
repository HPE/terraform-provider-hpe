// SPDX-FileCopyrightText: Copyright Hewlett Packard Enterprise Development LP
// SPDX-License-Identifier: Apache-2.0

package acctest

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"

	"github.com/HPE/terraform-provider-opsramp/internal/client"
	providerpkg "github.com/HPE/terraform-provider-opsramp/internal/provider"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
)

var loadEnvOnce sync.Once

// ProtoV6ProviderFactories instantiates the provider for acceptance tests.
func ProtoV6ProviderFactories() map[string]func() (tfprotov6.ProviderServer, error) {
	return map[string]func() (tfprotov6.ProviderServer, error){
		"opsramp": providerserver.NewProtocol6WithError(providerpkg.New("test")()),
	}
}

// PreCheck validates required acceptance-test environment variables.
func PreCheck(t *testing.T) func() {
	return func() {
		t.Helper()

		loadEnvOnce.Do(func() {
			loadEnvFileIfPresent(filepath.Join(".env"))
			loadEnvFileIfPresent(filepath.Join("..", "..", ".env"))
		})

		required := []string{"OPSRAMP_ENDPOINT", "OPSRAMP_TENANT", "OPSRAMP_CLIENT_ID", "OPSRAMP_CLIENT_SECRET"}
		for _, key := range required {
			if strings.TrimSpace(os.Getenv(key)) == "" {
				t.Fatalf("%s must be set for acceptance tests (or defined in .env)", key)
			}
		}
	}
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

// ProviderConfigHCL returns reusable Terraform configuration for provider wiring.
func ProviderConfigHCL() string {
	return fmt.Sprintf(`
provider "opsramp" {
	endpoint      = %q
	tenant        = %q
	client_id     = %q
	client_secret = %q
}
`, os.Getenv("OPSRAMP_ENDPOINT"), os.Getenv("OPSRAMP_TENANT"), os.Getenv("OPSRAMP_CLIENT_ID"), os.Getenv("OPSRAMP_CLIENT_SECRET"))
}

func loadEnvFileIfPresent(path string) {
	content, err := os.ReadFile(path)
	if err != nil {
		return
	}

	for _, line := range strings.Split(string(content), "\n") {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])
		value = strings.Trim(value, "\"")
		if key == "" {
			continue
		}

		if strings.TrimSpace(os.Getenv(key)) == "" {
			_ = os.Setenv(key, value)
		}
	}
}

// APIClient initializes an authenticated OpsRamp API client using acceptance-test env vars.
func APIClient(t *testing.T) (*client.OpsRampClient, error) {
	t.Helper()
	PreCheck(t)()

	return client.NewOpsRampClient(
		os.Getenv("OPSRAMP_CLIENT_ID"),
		os.Getenv("OPSRAMP_CLIENT_SECRET"),
		os.Getenv("OPSRAMP_ENDPOINT"),
		os.Getenv("OPSRAMP_TENANT"),
	)
}
