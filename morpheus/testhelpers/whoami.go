// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package testhelpers

import (
	"context"
	"sync"
	"testing"
)

// WhoamiBlock renders the whoami data source used by fixtures to resolve the
// caller's own tenant and user id at plan time, instead of hardcoding the
// master tenant. Append it once to a test's dependenciesConfig.
func WhoamiBlock() string {
	return "data \"hpe_morpheus_whoami\" \"current\" {}\n"
}

const (
	// WhoamiTenantIDRef is the HCL expression for the caller's tenant id.
	WhoamiTenantIDRef = "data.hpe_morpheus_whoami.current.tenant_id"
	// WhoamiUserIDRef is the HCL expression for the caller's user id.
	WhoamiUserIDRef = "data.hpe_morpheus_whoami.current.id"
)

var (
	tenancyOnce     sync.Once
	tenancyIsMaster bool
)

// IsMasterTenant reports whether the configured test caller is the master
// tenant, observed once per test binary via whoami. It fails the test on error.
func IsMasterTenant(t *testing.T) bool {
	t.Helper()

	tenancyOnce.Do(func() {
		ctx := context.Background()
		client, err := NewClientForServer(ctx, "")
		if err != nil {
			t.Fatalf("failed to create test client for tenancy check: %v", err)
		}

		who, hresp, err := client.AuthenticationAPI.Whoami(ctx).Execute()
		if err != nil {
			t.Fatalf("whoami failed during tenancy check: %v", err)
		}
		if hresp != nil {
			_ = hresp.Body.Close()
		}
		if who == nil {
			t.Fatalf("whoami returned an empty response during tenancy check")
		}

		tenancyIsMaster = who.IsMasterAccount != nil && *who.IsMasterAccount

		if tenancyIsMaster {
			t.Logf("observed tenancy: MASTER tenant")
		} else {
			t.Logf("observed tenancy: SUB-tenant (non-master)")
		}
	})

	return tenancyIsMaster
}

// TenantVisibility returns the visibility value a fixture should use given the
// caller's tenancy: "public" is only accepted for the master tenant, so a
// sub-tenant caller falls back to "private".
func TenantVisibility(t *testing.T) string {
	t.Helper()

	if IsMasterTenant(t) {
		return "public"
	}

	return "private"
}
