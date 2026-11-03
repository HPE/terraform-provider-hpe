// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package testhelpers

import (
	"os"
	"testing"
)

// Morpheus 9.1.2 verifies connectivity to the external storage service when a
// storage server is created: the create POST returns 400 "Host not reachable"
// for an unreachable service URL. On 9.1.0 the record was accepted and only a
// later sync failed, which is why create-based tests passed against a
// placeholder URL before.
//
// A create therefore now needs a reachable storage service. The placeholder
// values the tests used to hard-code (e.g. https://alletra.example.com, admin)
// exist on no appliance, so every create fails with "Host not reachable" -- a
// failure that looks like a provider bug but is only a missing fixture.
// Reading the storage server details from the environment keeps the tests
// appliance-agnostic, and skipping when they are absent makes the missing
// coverage explicit instead of reporting a false failure.
const (
	EnvStorageServerType     = "TF_ACC_STORAGE_SERVER_TYPE"
	EnvStorageServerURL      = "TF_ACC_STORAGE_SERVER_URL"
	EnvStorageServerUsername = "TF_ACC_STORAGE_SERVER_USERNAME"
	// #nosec G101 -- environment variable name, not an embedded credential.
	EnvStorageServerPassword = "TF_ACC_STORAGE_SERVER_PASSWORD"
)

// StorageServerFixture identifies a reachable storage service that storage
// server acceptance tests create against. All values are rendered into HCL
// as-is.
type StorageServerFixture struct {
	// TypeCode is the storage server type code, e.g. "alletra-mp".
	TypeCode string
	// URL is the reachable service URL.
	URL string
	// Username authenticates to the storage service.
	Username string
	// Password authenticates to the storage service.
	Password string
}

// RequireStorageServerFixture returns the storage server fixture for the
// appliance under test, skipping the test when it has not been fully
// configured.
//
// It skips rather than fails: an appliance without a reachable storage service
// is a legitimate target for the rest of the suite, and turning that into a
// failure would train people to ignore red results. Since 9.1.2 verifies
// reachability on create, there is no way to exercise create against a fake
// URL.
func RequireStorageServerFixture(t *testing.T) StorageServerFixture {
	t.Helper()

	fixture := StorageServerFixture{
		TypeCode: os.Getenv(EnvStorageServerType),
		URL:      os.Getenv(EnvStorageServerURL),
		Username: os.Getenv(EnvStorageServerUsername),
		Password: os.Getenv(EnvStorageServerPassword),
	}

	var missing []string

	for _, v := range []struct {
		name  string
		value string
	}{
		{EnvStorageServerType, fixture.TypeCode},
		{EnvStorageServerURL, fixture.URL},
		{EnvStorageServerUsername, fixture.Username},
		{EnvStorageServerPassword, fixture.Password},
	} {
		if v.value == "" {
			missing = append(missing, v.name)
		}
	}

	if len(missing) > 0 {
		msg := "storage server fixture not configured; set:"
		for _, name := range missing {
			msg += " " + name
		}

		t.Skip(msg)
	}

	return fixture
}
