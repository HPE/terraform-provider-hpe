// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package network

import (
	"context"
	"fmt"
	"net/http"
	"sync"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	sdk "github.com/HPE/terraform-provider-hpe/internal/sdk/oapigen"
	"github.com/HPE/terraform-provider-hpe/morpheus/utils/errfmt"
	"github.com/HPE/terraform-provider-hpe/morpheus/utils/networkconfig"
)

// networkType is a network type as the appliance describes it: the id, which
// is appliance-specific; the code, which is stable and keys the catalogue; and
// its shape — option types and editability flags — normalised from the same
// response, so the provider works from what this appliance actually has (its
// release, its plugins) rather than from the release the generated table was
// harvested on.
type networkType struct {
	id   int64
	code string
	spec *networkconfig.TypeSpec
}

// Network types are immutable reference data, so a lookup is cached for the
// life of the provider process, in both directions. A provider process serves
// one appliance.
var (
	typeCacheMu sync.Mutex
	typeByID    = map[int64]networkType{}
	typeByCode  = map[string]networkType{}
)

func cacheType(t networkType) {
	typeCacheMu.Lock()
	defer typeCacheMu.Unlock()

	typeByID[t.id] = t
	typeByCode[t.code] = t
}

// normalise builds the type's spec from the API object; a response the
// normaliser cannot read leaves the spec nil and is logged, not fatal — the
// id and code are still good, and the generated table stands in.
func normalise(ctx context.Context, t networkType, apiObject any) networkType {
	spec, err := networkconfig.FromAPI(apiObject)
	if err != nil {
		tflog.Warn(ctx, "network: could not read the network type's option types from the API; using the generated table",
			map[string]any{"type_code": t.code, "error": err.Error()})

		return t
	}

	spec.Live = true
	t.spec = &spec

	tflog.Debug(ctx, "network: network type read from the appliance",
		map[string]any{"type_code": t.code, "type_id": t.id, "option_types": len(spec.Fields)})

	return t
}

// networkTypeByID returns the network type with the given id.
func networkTypeByID(ctx context.Context, client *sdk.APIClient, id int64) (networkType, error) {
	typeCacheMu.Lock()
	t, ok := typeByID[id]
	typeCacheMu.Unlock()

	if ok {
		return t, nil
	}

	r, hresp, err := client.NetworksAPI.GetNetworkType(ctx, id).Execute()
	if err != nil || hresp.StatusCode != http.StatusOK {
		return networkType{}, fmt.Errorf("network type %d: %s", id, errfmt.ErrMsg(err, hresp))
	}

	if r == nil || r.NetworkType == nil || r.NetworkType.Id == nil || r.NetworkType.Code == nil {
		return networkType{}, fmt.Errorf("network type %d: the response has no id or code", id)
	}

	t = normalise(ctx, networkType{id: *r.NetworkType.Id, code: *r.NetworkType.Code}, r.NetworkType)
	cacheType(t)

	return t, nil
}

// networkTypeByCode returns the network type with the given code. The list
// endpoint filters on code exactly.
func networkTypeByCode(ctx context.Context, client *sdk.APIClient, code string) (networkType, error) {
	typeCacheMu.Lock()
	t, ok := typeByCode[code]
	typeCacheMu.Unlock()

	if ok {
		return t, nil
	}

	r, hresp, err := client.NetworksAPI.ListNetworkTypes(ctx).Code(code).Execute()
	if err != nil || hresp.StatusCode != http.StatusOK {
		return networkType{}, fmt.Errorf("network type %q: %s", code, errfmt.ErrMsg(err, hresp))
	}

	var matches []networkType

	if r != nil {
		for i := range r.NetworkTypes {
			nt := &r.NetworkTypes[i]
			if nt.Id != nil && nt.Code != nil && *nt.Code == code {
				matches = append(matches, normalise(ctx, networkType{id: *nt.Id, code: *nt.Code}, nt))
			}
		}
	}

	switch len(matches) {
	case 0:
		return networkType{}, fmt.Errorf("network type %q: the appliance has no network type with that code", code)
	case 1:
		cacheType(matches[0])

		return matches[0], nil
	default:
		return networkType{}, fmt.Errorf("network type %q: the appliance has %d network types with that code", code, len(matches))
	}
}

// resolveType returns the network type named by type_code or type_id in
// config, whichever is set and known. It is the create-time resolution: the
// type is required, so neither being known is an error.
func resolveType(ctx context.Context, client *sdk.APIClient, config NetworkModel) (networkType, error) {
	if !config.TypeCode.IsNull() && !config.TypeCode.IsUnknown() {
		return networkTypeByCode(ctx, client, config.TypeCode.ValueString())
	}

	if !config.TypeId.IsNull() && !config.TypeId.IsUnknown() {
		return networkTypeByID(ctx, client, config.TypeId.ValueInt64())
	}

	return networkType{}, fmt.Errorf("neither type_code nor type_id is known")
}

// typeSpec returns the shape of the network type identified by id (preferred)
// or code: live from the appliance when it can be reached, otherwise from the
// generated table by code. nil when neither is possible, or the type is not
// in the table — the resource then behaves generically and Morpheus enforces
// its rules on apply.
func typeSpec(ctx context.Context, client *sdk.APIClient, id int64, code string) *networkconfig.TypeSpec {
	if client != nil {
		var (
			t   networkType
			err error
		)

		switch {
		case id > 0:
			t, err = networkTypeByID(ctx, client, id)
		case code != "":
			t, err = networkTypeByCode(ctx, client, code)
		default:
			return nil
		}

		if err == nil {
			if t.spec != nil {
				return t.spec
			}

			code = t.code
		} else {
			tflog.Debug(ctx, "network: could not read the network type from the appliance; using the generated table",
				map[string]any{"error": err.Error()})
		}
	}

	return typeSpecFor(code)
}

// planTypeSpec returns the network's type shape at plan time. The type comes
// from the configuration (type_code or type_id) when set there, otherwise from
// the prior state, where Read records both.
func (r *Resource) planTypeSpec(ctx context.Context, config, state NetworkModel) *networkconfig.TypeSpec {
	id := knownInt64(config.TypeId)
	if id == 0 {
		id = knownInt64(state.TypeId)
	}

	code := knownString(config.TypeCode)
	if code == "" {
		code = knownString(state.TypeCode)
	}

	if id == 0 && code == "" {
		return nil
	}

	client, err := r.NewClient(ctx)
	if err != nil {
		client = nil
	}

	return typeSpec(ctx, client, id, code)
}

// typeSpecFor returns the generated table's entry for a type code, or nil.
func typeSpecFor(code string) *networkconfig.TypeSpec {
	if code == "" {
		return nil
	}

	spec, ok := networkconfig.Lookup(code)
	if !ok {
		return nil
	}

	return &spec
}

// liveOrTableSpec returns a resolved type's live shape, or the table's entry
// for its code when the live one could not be read.
func liveOrTableSpec(t networkType) *networkconfig.TypeSpec {
	if t.spec != nil {
		return t.spec
	}

	return typeSpecFor(t.code)
}

func knownString(v types.String) string {
	if v.IsNull() || v.IsUnknown() {
		return ""
	}

	return v.ValueString()
}

func knownInt64(v types.Int64) int64 {
	if v.IsNull() || v.IsUnknown() {
		return 0
	}

	return v.ValueInt64()
}

// knownStringOrNull returns v, or null when v is unknown: a write-through
// attribute that was not configured must not stay unknown after apply.
func knownStringOrNull(v types.String) types.String {
	if v.IsUnknown() {
		return types.StringNull()
	}

	return v
}

// knownInt64OrNull is knownStringOrNull for an int64.
func knownInt64OrNull(v types.Int64) types.Int64 {
	if v.IsUnknown() {
		return types.Int64Null()
	}

	return v
}
