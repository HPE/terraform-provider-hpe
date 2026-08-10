// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package morpheus_test

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

// The identity block examples published in the provider documentation. Users
// copy these, so a schema or validator change that makes one of them invalid is
// a documentation bug rather than only a code change.
//
// The attributes are repeated here rather than parsed out of the files, which
// would mean depending on an HCL parser. exampleFileSetsAttrs checks the two
// against each other so they cannot drift apart.
var identityExamples = map[string]struct {
	block string
	attrs map[string]string
}{
	"provider-pce-identity.tf": {
		block: "pce_identity",
		attrs: map[string]string{
			"client_id":     "client_id",
			"client_secret": "client_secret",
			"issuer_url":    "https://issuer.example.com",
			"location":      "location",
			"space":         "space",
		},
	},
	"provider-pce-identity-iamtoken.tf": {
		block: "pce_identity",
		attrs: map[string]string{
			"iam_token": "iam-token",
			"location":  "location",
			"space":     "space",
		},
	},
	"provider-pce-disconnected-identity.tf": {
		block: "pce_disconnected_identity",
		attrs: map[string]string{
			"client_id":     "client_id",
			"client_secret": "client_secret",
			"issuer_url":    "https://issuer.example.com",
			"location":      "location",
			"workspace_id":  "workspace_id",
			"broker_url":    "https://broker.example.com",
		},
	},
	"provider-pce-disconnected-identity-iamtoken.tf": {
		block: "pce_disconnected_identity",
		attrs: map[string]string{
			"iam_token":    "iam-token",
			"location":     "location",
			"workspace_id": "workspace_id",
			"broker_url":   "https://broker.example.com",
		},
	},
}

// exampleFileSetsAttrs reports whether the example file declares the block and
// every attribute this test validates on its behalf.
func exampleFileSetsAttrs(t *testing.T, name, block string, attrs map[string]string) {
	t.Helper()

	path := filepath.Join("..", "examples", "provider", "morpheus", name)

	body, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("could not read %s: %v", path, err)
	}

	content := string(body)

	if !strings.Contains(content, block+" {") {
		t.Errorf("%s does not declare a %s block", name, block)
	}

	for attr, value := range attrs {
		if !strings.Contains(content, attr+" ") {
			t.Errorf("%s does not set %q, but this test validates it", name, attr)
		}

		if !strings.Contains(content, fmt.Sprintf("%q", value)) {
			t.Errorf("%s does not contain the value %q used for %q", name, value, attr)
		}
	}
}

// Every published identity example has to be a configuration the provider
// accepts. A required attribute added to a block, or a new conflict, would
// otherwise leave the documentation telling users to write something that fails
// at "terraform validate".
func TestPublishedIdentityExamplesAreValid(t *testing.T) {
	for name, example := range identityExamples {
		t.Run(name, func(t *testing.T) {
			exampleFileSetsAttrs(t, name, example.block, example.attrs)

			diags := validateProviderConfig(t, emptyBlocks,
				func(t *testing.T, obj tftypes.Object, absent absentBlocks) map[string]tftypes.Value {
					return map[string]tftypes.Value{
						example.block: identityBlockValue(
							t, obj, absent, example.block, example.attrs,
						),
					}
				})

			if len(diags) != 0 {
				t.Errorf("example is not a valid configuration: %v", diags)
			}
		})
	}
}
