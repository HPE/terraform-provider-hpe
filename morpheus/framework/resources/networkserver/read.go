// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package networkserver

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"

	sdk "github.com/HPE/terraform-provider-hpe/internal/sdk/oapigen"

	"github.com/HPE/terraform-provider-hpe/morpheus/utils/errfmt"
	"github.com/HPE/terraform-provider-hpe/utils/convert"
	"github.com/HPE/terraform-provider-hpe/utils/schemadefaults"
)

func (r *networkServerResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	client, err := r.NewClient(ctx)
	if err != nil {
		errfmt.DiagClientError(&resp.Diagnostics, err)

		return
	}

	var state NetworkServerModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := state.Id.ValueInt64()

	result, httpResp, err := client.NetworksAPI.GetNetworkServer(ctx, id).Execute()
	if errfmt.IsNotFound(httpResp) {
		resp.State.RemoveResource(ctx)

		return
	}
	if err := errfmt.CheckResponse(err, httpResp); err != nil {
		errfmt.DiagError(&resp.Diagnostics, errfmt.OpRead, "network_server", "", err, httpResp)

		return
	}

	if result.NetworkServer == nil {
		resp.Diagnostics.AddError("API returned nil", "NetworkServer is nil in the response")

		return
	}

	newState, err := networkServerAsState(ctx, &state, result.NetworkServer)
	if err != nil {
		resp.Diagnostics.AddError("Failed to map network server", err.Error())

		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &newState)...)

	// Fill any schema-declared default the API omitted (null in state) so an
	// imported resource does not plan a change nobody made.
	resp.Diagnostics.Append(
		schemadefaults.Apply(ctx, NetworkServerResourceSchema(ctx), &resp.State)...,
	)
}

// networkServerAsState maps the API response onto the resource model. prior is
// the plan (on create/update) or the existing state (on read); it supplies the
// values the API never returns (the write-only password, the service token)
// and picks the config representation the user configured. On import prior is a
// zero model, so the type code drives the representation instead.
func networkServerAsState(
	ctx context.Context,
	prior *NetworkServerModel,
	ns *sdk.GetNetworkServer200ResponseNetworkServer,
) (NetworkServerModel, error) {
	state := NetworkServerModel{
		Id:            convert.Int64ToType(ns.Id),
		Name:          convert.StrToType(ns.Name),
		Description:   convert.StrToType(ns.Description.Get()),
		Enabled:       convert.BoolToType(ns.Enabled),
		Visibility:    convert.StrToType(ns.Visibility),
		ExternalId:    convert.StrToType(ns.ExternalId.Get()),
		InternalId:    convert.StrToType(ns.InternalId.Get()),
		Status:        convert.StrToType(ns.Status),
		StatusMessage: convert.StrToType(ns.StatusMessage.Get()),
		ServiceUrl:    serviceURLFromAPI(ns, prior),
		ServiceHost:   convert.StrToType(ns.ServiceHost.Get()),
		ServicePort:   convert.Int32ToType(ns.ServicePort.Get()),
		ServiceMode:   convert.StrToType(ns.ServiceMode.Get()),
		NetworkFilter: convert.StrToType(ns.NetworkFilter.Get()),
		TenantMatch:   convert.StrToType(ns.TenantMatch.Get()),
		CloudId:       convert.Int64ToType(ns.ZoneId.Get()),
	}

	switch {
	case ns.Type != nil && ns.Type.Code != nil:
		state.TypeCode = convert.StrToType(ns.Type.Code)
	case !prior.TypeCode.IsNull() && !prior.TypeCode.IsUnknown():
		state.TypeCode = prior.TypeCode
	default:
		// On import there is no prior value; a null type_code (required,
		// RequiresReplace) would make the next plan recreate the server.
		return NetworkServerModel{}, fmt.Errorf("network server %s: the API returned no type and none is known from state",
			idForMessage(ns.Id))
	}

	// service_username is not sensitive but the API may omit it; when it does,
	// preserve the configured value to avoid spurious drift.
	if v := ns.ServiceUsername.Get(); v != nil {
		state.ServiceUsername = types.StringValue(*v)
	} else {
		state.ServiceUsername = prior.ServiceUsername
	}

	// service_password_wo and service_token_wo are write-only (never in state);
	// their markers are carried forward from the prior model.
	state.ServicePasswordWo = prior.ServicePasswordWo
	state.ServiceTokenWo = prior.ServiceTokenWo
	state.ServicePasswordWoVersion = prior.ServicePasswordWoVersion
	state.WriteOnlyRevision = prior.WriteOnlyRevision

	state.DateCreated = convert.StrToType(timePtr(ns.DateCreated))
	state.LastUpdated = convert.StrToType(timePtr(ns.LastUpdated))

	// credential_id reflects a stored credential only. A "local" credential (or
	// none) means the server uses serviceUsername/servicePassword, so leave it
	// null.
	if ns.Credential != nil && ns.Credential.Type != nil && *ns.Credential.Type != "local" {
		state.CredentialId = convert.Int64ToType(ns.Credential.Id.Get())
	} else {
		state.CredentialId = types.Int64Null()
	}

	// tenant_ids: preserve a null prior value when the API returns no tenants,
	// so an unset attribute does not flip to an empty set.
	tenantIDs := responseTenantIDs(ns.Tenants)
	switch {
	case len(tenantIDs) > 0:
		set, d := types.SetValueFrom(ctx, types.Int64Type, tenantIDs)
		if d.HasError() {
			return state, fmt.Errorf("failed to convert tenant_ids: %s", d.Errors())
		}

		state.TenantIds = set
	case prior.TenantIds.IsNull() || prior.TenantIds.IsUnknown():
		state.TenantIds = types.SetNull(types.Int64Type)
	default:
		state.TenantIds = types.SetValueMust(types.Int64Type, []attr.Value{})
	}

	if err := setConfigRepresentation(ctx, &state, prior, ns); err != nil {
		return state, err
	}

	return state, nil
}

// setConfigRepresentation records the type-specific configuration in the same
// shape the user configured it: the typed config_juniper_sdn block for the SDN
// type, the free-form config map otherwise. On import (no prior representation)
// the reported type code decides. The unused representation is set to its null.
func setConfigRepresentation(
	ctx context.Context,
	state, prior *NetworkServerModel,
	ns *sdk.GetNetworkServer200ResponseNetworkServer,
) error {
	priorHasSdn := !prior.ConfigJuniperSdn.IsNull() && !prior.ConfigJuniperSdn.IsUnknown()
	priorHasConfig := !prior.Config.IsNull() && !prior.Config.IsUnknown()
	isSdn := ns.Type != nil && ns.Type.Code != nil && *ns.Type.Code == typeCodeJuniperSdn

	switch {
	case priorHasSdn || (!priorHasConfig && isSdn):
		priorAddress := types.StringNull()
		if priorHasSdn {
			priorAddress = prior.ConfigJuniperSdn.ControllerAddress
		}

		sdnCfg, err := parseJuniperSdnConfig(ctx, ns.Config, priorAddress)
		if err != nil {
			return err
		}

		state.ConfigJuniperSdn = sdnCfg
		state.Config = types.DynamicNull()
	case priorHasConfig:
		// The user owns the generic representation; echo the configured value
		// back, since the API enriches the stored map.
		state.Config = prior.Config
		state.ConfigJuniperSdn = NewConfigJuniperSdnValueNull()
	default:
		// Import of a non-SDN type: surface the stored config map.
		state.ConfigJuniperSdn = NewConfigJuniperSdnValueNull()

		if ns.Config != nil {
			dyn, err := convert.MapToDynamic(ctx, ns.Config)
			if err != nil {
				return fmt.Errorf("failed to convert config: %w", err)
			}

			state.Config = dyn
		} else {
			state.Config = types.DynamicNull()
		}
	}

	return nil
}

func responseTenantIDs(tenants []sdk.GetNetworkServer200ResponseNetworkServerTenantsInner) []int64 {
	var out []int64
	for _, t := range tenants {
		if t.Id != nil {
			out = append(out, *t.Id)
		}
	}

	return out
}

// timePtr formats a time pointer as an RFC 3339 string pointer for the date
// attributes, or nil when the API omitted it.
func timePtr(t *time.Time) *string {
	if t == nil {
		return nil
	}

	s := t.Format(time.RFC3339)

	return &s
}

// serviceURLFromAPI reads service_url back from the API, except for the SDN
// type: that plugin takes no URL and writes the controller's API version into
// the field after its first connection, which would otherwise plan as drift.
func serviceURLFromAPI(ns *sdk.GetNetworkServer200ResponseNetworkServer, prior *NetworkServerModel) types.String {
	if ns.Type != nil && ns.Type.Code != nil && *ns.Type.Code == typeCodeJuniperSdn {
		return prior.ServiceUrl
	}

	return convert.StrToType(ns.ServiceUrl.Get())
}

// idForMessage renders an optional id for a diagnostic without dereferencing
// a nil pointer.
func idForMessage(id *int64) string {
	if id == nil {
		return "(no id)"
	}

	return fmt.Sprintf("%d", *id)
}
