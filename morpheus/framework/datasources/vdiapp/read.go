// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package vdiapp

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/types"

	sdk "github.com/HPE/terraform-provider-hpe/internal/sdk/oapigen"

	providererrors "github.com/HPE/terraform-provider-hpe/morpheus/utils/errfmt"
	"github.com/HPE/terraform-provider-hpe/utils/convert"
)

func vdiAppAsState(
	app *sdk.GetVDIApps200ResponseVdiApp,
) VdiAppModel {
	state := VdiAppModel{
		Id:           convert.Int64ToType(app.Id),
		Name:         convert.StrToType(app.Name),
		LaunchPrefix: convert.StrToType(app.LaunchPrefix),
		DateCreated:  convert.TimeToType(app.DateCreated),
		LastUpdated:  convert.TimeToType(app.LastUpdated),
	}

	// Nullable fields
	if app.Description.IsSet() {
		state.Description = convert.StrToType(app.Description.Get())
	} else {
		state.Description = types.StringNull()
	}

	if app.IconPath.IsSet() {
		state.IconPath = convert.StrToType(app.IconPath.Get())
	} else {
		state.IconPath = types.StringNull()
	}

	if app.Logo.IsSet() {
		state.Logo = convert.StrToType(app.Logo.Get())
	} else {
		state.Logo = types.StringNull()
	}

	return state
}

func getByID(
	ctx context.Context,
	id int64,
	apiClient *sdk.APIClient,
) (*sdk.GetVDIApps200ResponseVdiApp, error) {
	r, hresp, err := apiClient.VDIAPI.GetVDIApps(ctx, id).Execute()
	if r == nil || err != nil || hresp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GET failed for vdi app %d: %s", id, providererrors.ErrMsg(err, hresp))
	}

	if r.VdiApp == nil {
		return nil, fmt.Errorf("GET failed for vdi app %d: response missing vdiApp", id)
	}

	app := *r.VdiApp

	return &app, nil
}

func getByName(
	ctx context.Context,
	name string,
	apiClient *sdk.APIClient,
) (*sdk.GetVDIApps200ResponseVdiApp, error) {
	rs, hresp, err := apiClient.VDIAPI.ListVDIApps(ctx).Name(name).Execute()
	if rs == nil || err != nil || hresp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GET failed for vdi app %s: %s", name, providererrors.ErrMsg(err, hresp))
	}

	var matched []sdk.ListVDIApps200ResponseAllOfVdiAppsInner

	for _, o := range rs.VdiApps {
		if o.Name != nil && *o.Name == name {
			matched = append(matched, o)
		}
	}

	if len(matched) == 0 {
		return nil, errors.New(ErrorNoVdiAppFound)
	} else if len(matched) > 1 {
		return nil, errors.New(ErrorMultipleVdiApps)
	}

	if matched[0].Id == nil {
		return nil, fmt.Errorf("GET failed for vdi app %s: response missing id", name)
	}

	return getByID(ctx, *matched[0].Id, apiClient)
}

func getVdiApp(
	ctx context.Context,
	config *VdiAppModel,
	apiClient *sdk.APIClient,
) (*sdk.GetVDIApps200ResponseVdiApp, error) {
	if !config.Id.IsNull() {
		return getByID(ctx, config.Id.ValueInt64(), apiClient)
	} else if !config.Name.IsNull() {
		return getByName(ctx, config.Name.ValueString(), apiClient)
	}

	return nil, errors.New(ErrorNoValidSearchTerms)
}

// Read refreshes the Terraform state with the latest data.
func (d *DataSource) Read(
	ctx context.Context,
	req datasource.ReadRequest,
	resp *datasource.ReadResponse,
) {
	var config VdiAppModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	apiClient, err := d.NewClient(ctx)
	if err != nil {
		resp.Diagnostics.AddError(summary, "could not create sdk client")

		return
	}

	app, err := getVdiApp(ctx, &config, apiClient)
	if err != nil {
		resp.Diagnostics.AddError(summary, err.Error())

		return
	}

	state := vdiAppAsState(app)

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
