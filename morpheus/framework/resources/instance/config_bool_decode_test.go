// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package instance_test

import (
	"testing"

	sdk "github.com/HPE/terraform-provider-hpe/internal/sdk/oapigen"
)

// These cover MORPH-16966: a string in a strict-boolean instance config field
// must no longer fail the decode of the whole GET /api/instances/{id} response.
//
// The test drives the real SDK decode via GetInstance200Response.UnmarshalJSON,
// so it exercises exactly the path the provider read uses, rather than the
// decode hook in isolation (which is unexported, and lives in a package whose
// files are replaced wholesale on SDK regeneration).

// decodeInstanceConfig decodes a GET instance response whose config.createUser
// holds the given raw JSON, and returns the decoded response.
func decodeInstanceConfig(t *testing.T, field, rawJSON string) *sdk.GetInstance200Response {
	t.Helper()

	body := []byte(`{"instance":{"id":42,"name":"web-01","config":{"` +
		field + `":` + rawJSON + `}}}`)

	var resp sdk.GetInstance200Response
	if err := resp.UnmarshalJSON(body); err != nil {
		t.Fatalf("decode failed for %s=%s: %v — the whole response was lost", field, rawJSON, err)
	}

	if resp.Instance == nil {
		t.Fatalf("instance is nil for %s=%s — the whole response was lost", field, rawJSON)
	}

	// The response surviving intact is the property under test: id and name
	// must always come back, whatever the bool field held.
	if resp.Instance.Id == nil || *resp.Instance.Id != 42 {
		t.Fatalf("id not recovered for %s=%s", field, rawJSON)
	}

	if resp.Instance.Name == nil || *resp.Instance.Name != "web-01" {
		t.Fatalf("name not recovered for %s=%s", field, rawJSON)
	}

	return &resp
}

func TestInstanceConfigCreateUserDecode(t *testing.T) {
	t.Parallel()

	tests := []struct {
		raw  string
		want bool
	}{
		// Genuine JSON booleans.
		{`true`, true},
		{`false`, false},

		// Forms mapstructure's weak typing already accepted.
		{`"true"`, true},
		{`"false"`, false},
		{`"1"`, true},
		{`"0"`, false},
		{`""`, false},

		// The pre-existing on/off coercion.
		{`"on"`, true},
		{`"off"`, false},

		// The forms this change adds.
		{`"yes"`, true},
		{`"no"`, false},
		{`"y"`, true},
		{`"n"`, false},
		{`"t"`, true},
		{`"f"`, false},
		{`"enabled"`, true},
		{`"disabled"`, false},

		// Case and surrounding whitespace are tolerated.
		{`"YES"`, true},
		{`"On"`, true},
		{`" yes "`, true},
		{`"OFF"`, false},

		// The crux: an unrecognised string resolves to false rather than
		// failing the whole response. These are exactly the values that used
		// to lose the entire instance.
		{`"yes please"`, false},
		{`"banana"`, false},
		{`"maybe"`, false},
		{`"2"`, false},
	}

	for _, tt := range tests {
		t.Run(tt.raw, func(t *testing.T) {
			t.Parallel()

			resp := decodeInstanceConfig(t, "createUser", tt.raw)

			got := resp.Instance.Config.CreateUser
			if got == nil {
				t.Fatalf("createUser is nil for %s, want %v", tt.raw, tt.want)
			}

			if *got != tt.want {
				t.Errorf("createUser = %v for %s, want %v", *got, tt.raw, tt.want)
			}
		})
	}
}

// The fragility is not specific to createUser: it affects every strict-boolean
// field in the config. A stray string in any of them must be tolerated, and
// the coercion must apply per field.
func TestInstanceConfigOtherBoolFieldsDecode(t *testing.T) {
	t.Parallel()

	// boolField reads the given config field back as a *bool, so each field's
	// coerced value can be asserted rather than only its survival.
	boolField := func(cfg *sdk.GetInstance200ResponseInstanceConfig, field string) *bool {
		switch field {
		case "createUser":
			return cfg.CreateUser
		case "isVpcSelectable":
			return cfg.IsVpcSelectable
		case "createBackup":
			return cfg.CreateBackup
		default:
			return nil
		}
	}

	for _, field := range []string{"createUser", "isVpcSelectable", "createBackup"} {
		t.Run(field, func(t *testing.T) {
			t.Parallel()

			// "yes" would have failed the whole response before this change;
			// now it coerces to true.
			resp := decodeInstanceConfig(t, field, `"yes"`)
			if got := boolField(resp.Instance.Config, field); got == nil || !*got {
				t.Errorf("%s = %v for \"yes\", want true", field, got)
			}

			// An outright unrecognised value likewise, degrading to false.
			resp = decodeInstanceConfig(t, field, `"banana"`)
			if got := boolField(resp.Instance.Config, field); got == nil || *got {
				t.Errorf("%s = %v for \"banana\", want false", field, got)
			}
		})
	}
}

// Several unexpected bool values in one response must not compound into a
// failure — every one is tolerated, and the instance survives.
func TestInstanceConfigMultipleBadBoolsDecode(t *testing.T) {
	t.Parallel()

	body := []byte(`{"instance":{"id":42,"name":"web-01","config":{` +
		`"createUser":"yes","createBackup":"nope","isVpcSelectable":"maybe"}}}`)

	var resp sdk.GetInstance200Response
	if err := resp.UnmarshalJSON(body); err != nil {
		t.Fatalf("decode failed with several bad bools: %v", err)
	}

	if resp.Instance == nil || resp.Instance.Id == nil || *resp.Instance.Id != 42 {
		t.Fatal("instance not recovered with several bad bools in the config")
	}

	cfg := resp.Instance.Config
	if cfg == nil {
		t.Fatal("config not recovered")
	}

	// The recognised one is true; the unrecognised ones are false.
	if cfg.CreateUser == nil || !*cfg.CreateUser {
		t.Error("createUser: expected true from \"yes\"")
	}

	if cfg.CreateBackup == nil || *cfg.CreateBackup {
		t.Error("createBackup: expected false from an unrecognised value")
	}
}

// A stray bool value must not cost the fields around it. Everything else in the
// response still decodes.
func TestInstanceConfigBadBoolLeavesOtherFieldsIntact(t *testing.T) {
	t.Parallel()

	body := []byte(`{"instance":{"id":42,"name":"web-01","hostName":"host-a",` +
		`"config":{"createUser":"banana","customOptions":{"k":"v"}}}}`)

	var resp sdk.GetInstance200Response
	if err := resp.UnmarshalJSON(body); err != nil {
		t.Fatalf("decode failed: %v", err)
	}

	if resp.Instance == nil {
		t.Fatal("instance lost")
	}

	if resp.Instance.HostName == nil || *resp.Instance.HostName != "host-a" {
		t.Error("hostName was lost alongside the bad bool")
	}
}
