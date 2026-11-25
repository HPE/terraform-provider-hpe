// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package affinityshape

import (
	"context"
	"errors"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"

	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/HPE/terraform-provider-hpe/morpheus/utils/versioncheck"
)

func v(t *testing.T, s string) *versioncheck.Version {
	t.Helper()
	ver, err := versioncheck.Parse(s)
	if err != nil {
		t.Fatal(err)
	}

	return ver
}

func TestUnitDecideCreate(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	empty := types.SetValueMust(types.Int64Type, nil)

	cases := []struct {
		name    string
		version string
		err     error
		servers types.Set
		refuse  bool
	}{
		{"9.0 empty allowed", "9.0.2.18", nil, types.SetNull(types.Int64Type), false},
		{"9.1 null refused", "9.1.0.2", nil, types.SetNull(types.Int64Type), true},
		{"9.1 empty refused", "9.1.0.2", nil, empty, true},
		{"9.1 one server ok", "9.1.0.2", nil, SetOf(7), false},
		{"9.1 unknown set passes", "9.1.0.2", nil, types.SetUnknown(types.Int64Type), false},
		{"9.1 unknown element passes", "9.1.0.2", nil, types.SetValueMust(types.Int64Type, []attr.Value{types.Int64Unknown()}), false},
		{"9.2 refused", "9.2.0", nil, types.SetNull(types.Int64Type), true},
		{"unknown version fails open", "", errors.New("whoami failed"), types.SetNull(types.Int64Type), false},
	}
	for _, c := range cases {
		var ver *versioncheck.Version
		if c.version != "" {
			ver = v(t, c.version)
		}
		d := DecideCreate(ctx, c.servers, ver, c.err)
		if d.HasError() != c.refuse {
			t.Errorf("%s: refuse=%v, want %v (%v)", c.name, d.HasError(), c.refuse, d)
		}
	}
}

func TestUnitDecideLastMemberRemoval(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	v91 := v(t, "9.1.0.2")

	skip, d := DecideLastMemberRemoval(ctx, []int64{5}, 5, v91, nil)
	if !skip || d.WarningsCount() != 1 || d.HasError() {
		t.Error("removing the only member on 9.1 must be skipped with a warning")
	}
	if skip, _ := DecideLastMemberRemoval(ctx, []int64{5, 6}, 5, v91, nil); skip {
		t.Error("removing one of two members must proceed")
	}
	if skip, _ := DecideLastMemberRemoval(ctx, []int64{6}, 5, v91, nil); skip {
		t.Error("a server that is not the remaining member must proceed (already gone)")
	}
	if skip, _ := DecideLastMemberRemoval(ctx, []int64{5}, 5, v(t, "9.0.2.18"), nil); skip {
		t.Error("9.0 must allow emptying a group")
	}
	if skip, _ := DecideLastMemberRemoval(ctx, []int64{5}, 5, nil, errors.New("down")); skip {
		t.Error("unknown version must fail open")
	}
}
