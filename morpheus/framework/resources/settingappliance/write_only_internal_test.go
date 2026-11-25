// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package settingappliance

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"reflect"
	"sort"
	"testing"
)

// The tracker must cover exactly the schema's write-only attributes: an
// attribute left out would never be re-sent when it changes, and one tracked
// but not write-only would be re-sent for nothing. Each has a hash extractor,
// bar the Twilio auth token, which Morpheus returns masked and with no hash.
func TestUnitWriteOnlyTrackerCoversTheSchema(t *testing.T) {
	t.Parallel()

	noAPIHash := map[string]bool{"twilio_auth_token": true}

	var fromSchema, wantHashed []string

	for name, a := range SettingApplianceResourceSchema(context.Background()).Attributes {
		if a.IsWriteOnly() {
			fromSchema = append(fromSchema, name)

			if !noAPIHash[name] {
				wantHashed = append(wantHashed, name)
			}
		}
	}

	tracked := wo.Names()

	var hashed []string
	for name := range wo.APIHash {
		hashed = append(hashed, name)
	}

	sort.Strings(fromSchema)
	sort.Strings(wantHashed)
	sort.Strings(tracked)
	sort.Strings(hashed)

	if !reflect.DeepEqual(tracked, fromSchema) {
		t.Errorf("tracked %v, but the schema's write-only attributes are %v", tracked, fromSchema)
	}

	if !reflect.DeepEqual(hashed, wantHashed) {
		t.Errorf("hash extractors for %v, want them for %v", hashed, wantHashed)
	}
}

func TestUnitHashExtractors(t *testing.T) {
	t.Parallel()

	settings := readFixture(t)

	// The SMTP password: the appliance's hash, as it reports it.
	if hash, present := smtpPasswordHash(settings); !present ||
		hash != "example-smtp-hash" {
		t.Errorf("smtp_password: got (%q, %v)", hash, present)
	}

	// The proxy password is not set: its hash is rendered null. Whichever way
	// the decoder records that, it must never read as a hash.
	if hash, _ := proxyPasswordHash(settings); hash != "" {
		t.Errorf("proxy_password: a null hash must not read as one, got %q", hash)
	}

	// The currency key comes back in clear: its digest, never the key, is
	// what is recorded.
	sum := sha256.Sum256([]byte("example"))

	if hash, present := currencyKeyDigest(settings); !present || hash != hex.EncodeToString(sum[:]) {
		t.Errorf("currency_provider_api_key: got (%q, %v), want the sha256 of the key", hash, present)
	}

	// No key: a comparable "no secret".
	empty := ""
	settings.CurrencyKey.Set(&empty)

	if hash, present := currencyKeyDigest(settings); !present || hash != "" {
		t.Errorf("currency_provider_api_key, empty: got (%q, %v), want (\"\", true)", hash, present)
	}

	settings.CurrencyKey.Set(nil)

	if hash, present := currencyKeyDigest(settings); !present || hash != "" {
		t.Errorf("currency_provider_api_key, null: got (%q, %v), want (\"\", true)", hash, present)
	}

	// An object without the fields, or no settings object at all, carries no
	// information.
	extractors := map[string]func(any) (string, bool){
		"smtp_password":             smtpPasswordHash,
		"proxy_password":            proxyPasswordHash,
		"currency_provider_api_key": currencyKeyDigest,
	}

	for name, extract := range extractors {
		if _, present := extract(&applianceSettings{}); present {
			t.Errorf("%s: an omitted field must read as not present", name)
		}

		if _, present := extract("not settings"); present {
			t.Errorf("%s: another type must read as not present", name)
		}

		if _, present := extract((*applianceSettings)(nil)); present {
			t.Errorf("%s: a nil object must read as not present", name)
		}
	}
}
