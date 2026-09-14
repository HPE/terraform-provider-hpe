// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

//go:build sweep

package sweep

import (
	"strconv"
	"strings"
	"testing"

	sdk "github.com/HPE/terraform-provider-hpe/internal/sdk/oapigen"
)

func tenantItem(id int64, parent *int64) sdk.ListTenants200ResponseAllOfAccountsInner {
	item := sdk.ListTenants200ResponseAllOfAccountsInner{Id: &id}
	if parent != nil {
		item.Parent = &sdk.ListTenants200ResponseAllOfAccountsInnerParent{Id: parent}
	}

	return item
}

func ptr(v int64) *int64 { return &v }

// TestSortDeepestFirst pins the sweep order for a leaked tenant tree: every
// tenant is listed after all of its descendants, so children are deleted
// before their parents; tenants at the same depth keep their listing order.
func TestSortDeepestFirst(t *testing.T) {
	// master(1) -> 10, 20 ; 10 -> 11, 12 ; 11 -> 111 ; plus 99 with an unknown
	// (unlisted) parent and 5 with no id at all. Listed shallowest-first to
	// prove the sort actually reorders.
	accounts := []sdk.ListTenants200ResponseAllOfAccountsInner{
		tenantItem(1, nil),
		tenantItem(10, ptr(1)),
		tenantItem(20, ptr(1)),
		tenantItem(11, ptr(10)),
		tenantItem(12, ptr(10)),
		tenantItem(111, ptr(11)),
		tenantItem(99, ptr(12345)),
		{},
	}

	sortDeepestFirst(accounts)

	var got []string
	for _, a := range accounts {
		if a.Id == nil {
			got = append(got, "-")

			continue
		}
		got = append(got, strconv.FormatInt(*a.Id, 10))
	}

	// depth 3: 111; depth 2: 11, 12; depth 1: 10, 20, 99; depth 0: 1, (no id)
	want := "111,11,12,10,20,99,1,-"
	if s := strings.Join(got, ","); s != want {
		t.Fatalf("sortDeepestFirst order = %s, want %s", s, want)
	}
}

// TestSortDeepestFirstTolerantOfCycles guards the depth walk against a parent
// cycle in the listing (which a healthy appliance never returns): it must
// terminate, keep every tenant, and -- since the walk caps both depths at the
// list length, making them equal -- leave the two in their listing order.
func TestSortDeepestFirstTolerantOfCycles(t *testing.T) {
	accounts := []sdk.ListTenants200ResponseAllOfAccountsInner{
		tenantItem(1, ptr(2)),
		tenantItem(2, ptr(1)),
	}

	sortDeepestFirst(accounts)

	if len(accounts) != 2 {
		t.Fatalf("sortDeepestFirst dropped tenants: got %d, want 2", len(accounts))
	}

	got := strconv.FormatInt(*accounts[0].Id, 10) + "," + strconv.FormatInt(*accounts[1].Id, 10)
	if want := "1,2"; got != want {
		t.Fatalf("sortDeepestFirst order = %s, want %s (equal capped depths keep listing order)",
			got, want)
	}
}
