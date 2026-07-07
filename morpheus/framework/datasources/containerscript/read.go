// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package containerscript

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/types"

	sdk "github.com/HPE/terraform-provider-hpe/internal/sdk/oapigen"

	providererrors "github.com/HPE/terraform-provider-hpe/morpheus/utils/errfmt"
	"github.com/HPE/terraform-provider-hpe/utils/convert"
)

func containerScriptAsState(
	ctx context.Context,
	s *sdk.GetScript200ResponseContainerScript,
) (ContainerScriptModel, error) {
	var labels types.Set
	if s.Labels != nil {
		l, diags := types.SetValueFrom(ctx, types.StringType, s.Labels)
		if diags.HasError() {
			return ContainerScriptModel{}, fmt.Errorf("error creating labels set")
		}

		labels = l
	} else {
		labels = types.SetNull(types.StringType)
	}

	state := ContainerScriptModel{
		Id:            convert.Int64ToType(s.Id),
		Name:          convert.StrToType(s.Name),
		Code:          convert.StrToType(s.Code),
		Category:      convert.StrToType(s.Category.Get()),
		SortOrder:     convert.Int64ToType(s.SortOrder),
		ScriptVersion: convert.StrToType(s.ScriptVersion),
		ScriptPhase:   convert.StrToType(s.ScriptPhase),
		ScriptType:    convert.StrToType(s.ScriptType),
		Script:        convert.StrToType(s.Script),
		RunAsUser:     convert.StrToType(s.RunAsUser.Get()),
		SudoUser:      convert.BoolToType(s.SudoUser),
		FailOnError:   convert.BoolToType(s.FailOnError),
		DateCreated:   convert.TimeToType(s.DateCreated),
		LastUpdated:   convert.TimeToType(s.LastUpdated),
		Labels:        labels,
	}

	// ScriptMethod — nullable
	if s.ScriptMethod.IsSet() {
		state.ScriptMethod = convert.StrToType(s.ScriptMethod.Get())
	} else {
		state.ScriptMethod = types.StringNull()
	}

	// ScriptService — nullable
	if s.ScriptService.IsSet() {
		state.ScriptService = convert.StrToType(s.ScriptService.Get())
	} else {
		state.ScriptService = types.StringNull()
	}

	// Account — nested {id, name}
	if s.Account != nil {
		state.Account = NewAccountValueMust(
			AccountValue{}.AttributeTypes(ctx),
			map[string]attr.Value{
				"id":   convert.Int64ToType(s.Account.Id),
				"name": convert.StrToType(s.Account.Name),
			},
		)
	} else {
		state.Account = NewAccountValueNull()
	}

	return state, nil
}

func getByID(
	ctx context.Context,
	id int64,
	apiClient *sdk.APIClient,
) (ContainerScriptModel, error) {
	r, hresp, err := apiClient.LibraryAPI.GetScript(ctx, id).Execute()
	if err != nil || r == nil || r.ContainerScript == nil {
		return ContainerScriptModel{}, fmt.Errorf(
			"GET failed for container script %d: %s", id, providererrors.ErrMsg(err, hresp),
		)
	}

	return containerScriptAsState(ctx, r.ContainerScript)
}

func getByName(
	ctx context.Context,
	name string,
	apiClient *sdk.APIClient,
) (ContainerScriptModel, error) {
	rs, hresp, err := apiClient.LibraryAPI.ListScripts(ctx).Phrase(name).Execute()
	if err != nil || rs == nil || hresp == nil || hresp.StatusCode != http.StatusOK {
		return ContainerScriptModel{}, fmt.Errorf(
			"GET failed for container script %s: %s", name, providererrors.ErrMsg(err, hresp),
		)
	}

	var matched []sdk.ListScripts200ResponseAllOfContainerScriptsInner

	for _, o := range rs.ContainerScripts {
		if o.Name != nil && *o.Name == name {
			matched = append(matched, o)
		}
	}

	if len(matched) == 0 {
		return ContainerScriptModel{}, errors.New(ErrorNoScriptFound)
	} else if len(matched) > 1 {
		return ContainerScriptModel{}, errors.New(ErrorMultipleScripts)
	}

	if matched[0].Id == nil {
		return ContainerScriptModel{}, fmt.Errorf(
			"GET failed for container script %s: response missing id", name,
		)
	}

	// Fetch the full resource by ID for complete field data.
	return getByID(ctx, *matched[0].Id, apiClient)
}

func getContainerScript(
	ctx context.Context,
	config *ContainerScriptModel,
	apiClient *sdk.APIClient,
) (ContainerScriptModel, error) {
	if !config.Id.IsNull() {
		return getByID(ctx, config.Id.ValueInt64(), apiClient)
	} else if !config.Name.IsNull() {
		return getByName(ctx, config.Name.ValueString(), apiClient)
	}

	return ContainerScriptModel{}, errors.New(ErrorNoValidSearchTerms)
}

// Read refreshes the Terraform state with the latest data.
func (d *DataSource) Read(
	ctx context.Context,
	req datasource.ReadRequest,
	resp *datasource.ReadResponse,
) {
	var config ContainerScriptModel

	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	apiClient, err := d.NewClient(ctx)
	if err != nil {
		resp.Diagnostics.AddError(
			summary,
			"could not create sdk client",
		)

		return
	}

	state, err := getContainerScript(ctx, &config, apiClient)
	if err != nil {
		resp.Diagnostics.AddError(
			summary,
			err.Error(),
		)

		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
