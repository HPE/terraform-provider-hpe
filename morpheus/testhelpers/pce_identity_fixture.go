// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package testhelpers

import (
	"os"
	"testing"
)

// The PCE identity acceptance test authenticates through GreenLake IAM rather
// than the appliance credentials the rest of the suite uses, so it needs an API
// client, the PCE instance it belongs to, and a cloud to read back through the
// resulting token. None of that can be created from Terraform, so it must
// already exist in the GreenLake tenant under test.
//
// The naming splits on who consumes the value. The six credentials must carry
// the TF_VAR_ prefix: the provider block returned by ProviderBlockPceIdentity
// declares matching Terraform variables, and Terraform reads them from the
// environment itself. CloudName has no Terraform variable behind it -- the test
// renders it straight into HCL -- so it takes the TF_ACC_ prefix used by the
// NSX-T fixture.
const (
	EnvPceIdentityClientID = "TF_VAR_testacc_pce_identity_client_id"
	// #nosec G101 -- environment variable name, not an embedded credential.
	EnvPceIdentityClientSecret = "TF_VAR_testacc_pce_identity_client_secret"
	EnvPceIdentityIssuerURL    = "TF_VAR_testacc_pce_identity_issuer_url"
	EnvPceIdentityLocation     = "TF_VAR_testacc_pce_identity_location"
	EnvPceIdentitySpace        = "TF_VAR_testacc_pce_identity_space"
	EnvPceIdentityBrokerURL    = "TF_VAR_testacc_pce_identity_broker_url"
	EnvPceIdentityCloudName    = "TF_ACC_PCE_IDENTITY_CLOUD_NAME"
)

// PceIdentityFixture identifies the GreenLake API client and PCE instance that
// the PCE identity auth-flow test authenticates with.
//
// Only CloudName is rendered into HCL by the test. The credentials are consumed
// by Terraform, which reads them from the TF_VAR_ environment directly, so the
// remaining fields are captured to define the fixture's surface and to drive
// the presence check in RequirePceIdentityFixture -- they are deliberately not
// dead weight to be removed, and having them to hand means a future assertion
// on, say, Location needs no new plumbing.
type PceIdentityFixture struct {
	// ClientID is the GreenLake API client used to mint an IAM token.
	ClientID string
	// ClientSecret is the secret for ClientID.
	ClientSecret string
	// IssuerURL is the GreenLake IAM issuer that mints the access token. It is
	// the "Issuer" URL shown for the API client.
	IssuerURL string
	// Location is the PCE instance's location.
	Location string
	// Space is the GreenLake space the PCE instance sits in.
	Space string
	// BrokerURL overrides the PCE broker. It is optional and therefore not
	// required by RequirePceIdentityFixture: leaving it empty selects the
	// HPE-hosted cloud broker, which is the case worth exercising by default.
	BrokerURL string
	// CloudName is the name of a cloud visible through the token exchange. The
	// test renders it into HCL and reads it via the hpe_morpheus_cloud data
	// source to prove the exchange produced a usable Morpheus session.
	CloudName string
}

// RequirePceIdentityFixture returns the PCE identity fixture for the GreenLake
// tenant under test, skipping the test when it has not been configured.
//
// It skips rather than fails: a runner without GreenLake access is a legitimate
// target for the rest of the suite, and turning that into a failure would train
// people to ignore red results.
func RequirePceIdentityFixture(t *testing.T) PceIdentityFixture {
	t.Helper()

	fixture := PceIdentityFixture{
		ClientID:     os.Getenv(EnvPceIdentityClientID),
		ClientSecret: os.Getenv(EnvPceIdentityClientSecret),
		IssuerURL:    os.Getenv(EnvPceIdentityIssuerURL),
		Location:     os.Getenv(EnvPceIdentityLocation),
		Space:        os.Getenv(EnvPceIdentitySpace),
		BrokerURL:    os.Getenv(EnvPceIdentityBrokerURL),
		CloudName:    os.Getenv(EnvPceIdentityCloudName),
	}

	var missing []string

	// BrokerURL is absent by design: it is optional, and an empty value is the
	// meaningful default rather than a misconfiguration.
	for _, v := range []struct {
		name  string
		value string
	}{
		{EnvPceIdentityClientID, fixture.ClientID},
		{EnvPceIdentityClientSecret, fixture.ClientSecret},
		{EnvPceIdentityIssuerURL, fixture.IssuerURL},
		{EnvPceIdentityLocation, fixture.Location},
		{EnvPceIdentitySpace, fixture.Space},
		{EnvPceIdentityCloudName, fixture.CloudName},
	} {
		if v.value == "" {
			missing = append(missing, v.name)
		}
	}

	if len(missing) > 0 {
		msg := "PCE identity fixture not configured; set:"
		for _, name := range missing {
			msg += " " + name
		}

		t.Skip(msg)
	}

	return fixture
}
