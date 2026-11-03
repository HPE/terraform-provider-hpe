// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package testhelpers

import (
	"os"
	"strings"
	"testing"
)

// Morpheus 9.1.2 verifies connectivity to the external IPAM service when a
// network pool server is created: the create POST returns 400 "Host not
// reachable" for an unreachable service URL. On 9.1.0 the record was accepted
// and only a later sync failed, which is why create-based tests passed against
// a placeholder URL before.
//
// A create therefore now needs a reachable IPAM. The placeholder values the
// tests used to hard-code (e.g. https://infoblox.example.com/wapi/v2.12, admin)
// exist on no appliance, so every create fails with "Host not reachable" -- a
// failure that looks like a provider bug but is only a missing fixture.
// Reading the IPAM details from the environment keeps the tests
// appliance-agnostic, and skipping when they are absent makes the missing
// coverage explicit instead of reporting a false failure.
const (
	EnvIpamType     = "TF_ACC_IPAM_TYPE"
	EnvIpamURL      = "TF_ACC_IPAM_URL"
	EnvIpamUsername = "TF_ACC_IPAM_USERNAME"
	// #nosec G101 -- environment variable name, not an embedded credential.
	EnvIpamPassword = "TF_ACC_IPAM_PASSWORD"
	// EnvIpamAppID is the phpIPAM application ID; required when the type is
	// phpipam (the appliance refuses the create without it) and unused by the
	// other types.
	EnvIpamAppID = "TF_ACC_IPAM_APP_ID"
)

// IpamFixture identifies a reachable IPAM integration that network pool server
// acceptance tests create against. All values are rendered into HCL as-is.
type IpamFixture struct {
	// TypeCode is the pool server type code, e.g. "infoblox", "phpipam" or
	// "bluecat".
	TypeCode string
	// URL is the reachable service URL.
	URL string
	// Username authenticates to the IPAM service.
	Username string
	// Password authenticates to the IPAM service.
	Password string
	// AppID is the phpIPAM application ID (empty for other types).
	AppID string
}

// WithAppID appends an app_id attribute to a rendered pool server resource
// block when the fixture carries one, so the same HCL serves every IPAM type.
// The block's closing brace must be the last one in hcl.
func (f IpamFixture) WithAppID(hcl string) string {
	if f.AppID == "" {
		return hcl
	}

	i := strings.LastIndex(hcl, "}")
	if i < 0 {
		return hcl
	}

	return hcl[:i] + "  app_id = \"" + f.AppID + "\"\n" + hcl[i:]
}

// RequireIpamFixture returns the IPAM fixture for the appliance under test,
// skipping the test when it has not been fully configured.
//
// It skips rather than fails: an appliance without a reachable IPAM is a
// legitimate target for the rest of the suite, and turning that into a failure
// would train people to ignore red results. Since 9.1.2 verifies reachability
// on create, there is no way to exercise create against a fake URL.
func RequireIpamFixture(t *testing.T) IpamFixture {
	t.Helper()

	fixture := IpamFixture{
		TypeCode: os.Getenv(EnvIpamType),
		URL:      os.Getenv(EnvIpamURL),
		Username: os.Getenv(EnvIpamUsername),
		Password: os.Getenv(EnvIpamPassword),
		AppID:    os.Getenv(EnvIpamAppID),
	}

	var missing []string

	for _, v := range []struct {
		name  string
		value string
	}{
		{EnvIpamType, fixture.TypeCode},
		{EnvIpamURL, fixture.URL},
		{EnvIpamUsername, fixture.Username},
		{EnvIpamPassword, fixture.Password},
	} {
		if v.value == "" {
			missing = append(missing, v.name)
		}
	}

	if fixture.TypeCode == "phpipam" && fixture.AppID == "" {
		missing = append(missing, EnvIpamAppID)
	}

	if len(missing) > 0 {
		msg := "IPAM fixture not configured; set:"
		for _, name := range missing {
			msg += " " + name
		}

		t.Skip(msg)
	}

	return fixture
}
