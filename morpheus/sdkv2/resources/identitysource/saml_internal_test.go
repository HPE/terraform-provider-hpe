package identitysource

import "testing"

// TestSamlProviderSettings verifies extraction of the SAML SP metadata that
// Morpheus returns under userSource.providerSettings (entity_id, acs_url,
// sp_metadata are computed outputs read from the raw response body).
func TestSamlProviderSettings(t *testing.T) {
	t.Parallel()

	body := []byte(`{
		"userSource": {
			"id": 7,
			"name": "example",
			"providerSettings": {
				"entityId": "https://morpheus.example.com/saml/abc123",
				"acsUrl": "https://morpheus.example.com/external-login/callback/abc123",
				"spMetadata": "<EntityDescriptor entityID=\"https://morpheus.example.com/saml/abc123\"/>"
			}
		}
	}`)

	entityID, acsURL, spMetadata, err := samlProviderSettings(body)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if want := "https://morpheus.example.com/saml/abc123"; entityID != want {
		t.Errorf("entityID = %q, want %q", entityID, want)
	}
	if want := "https://morpheus.example.com/external-login/callback/abc123"; acsURL != want {
		t.Errorf("acsURL = %q, want %q", acsURL, want)
	}
	if want := `<EntityDescriptor entityID="https://morpheus.example.com/saml/abc123"/>`; spMetadata != want {
		t.Errorf("spMetadata = %q, want %q", spMetadata, want)
	}
}

// TestSamlProviderSettingsMissing verifies that a response without
// providerSettings yields empty values rather than an error.
func TestSamlProviderSettingsMissing(t *testing.T) {
	t.Parallel()

	entityID, acsURL, spMetadata, err := samlProviderSettings([]byte(`{"userSource":{"id":7}}`))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if entityID != "" || acsURL != "" || spMetadata != "" {
		t.Errorf("expected empty SP settings, got %q %q %q", entityID, acsURL, spMetadata)
	}
}
