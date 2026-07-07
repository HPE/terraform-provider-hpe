// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package provisioninglicense

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"

	sdk "github.com/HPE/terraform-provider-hpe/internal/sdk/oapigen"
	"github.com/HPE/terraform-provider-hpe/morpheus/configure"
	providererrors "github.com/HPE/terraform-provider-hpe/morpheus/utils/errfmt"
	"github.com/HPE/terraform-provider-hpe/utils/convert"
)

const (
	summary                 = "read provisioning license data source"
	ErrorNoValidSearchTerms = `no valid search terms - an id or name is required`
	ErrorRunningPreApply    = `Error running pre-apply plan: exit status 1`
	ErrorNoLicenseFound     = `no provisioning license found`
	ErrorMultipleLicenses   = `multiple provisioning licenses were returned`
)

// Ensure the implementation satisfies the expected interfaces.
var _ datasource.DataSource = &DataSource{}

// NewDataSource is a helper function to simplify the provider implementation.
func NewDataSource() datasource.DataSource {
	return &DataSource{}
}

// DataSource is the data source implementation.
type DataSource struct {
	configure.DataSourceWithMorpheusConfigure
	datasource.DataSource
}

// Metadata returns the data source type name.
func (d *DataSource) Metadata(
	_ context.Context,
	req datasource.MetadataRequest,
	resp *datasource.MetadataResponse,
) {
	resp.TypeName = req.ProviderTypeName + "_" + "provisioning_license"
}

// Schema defines the schema for the data source.
func (d *DataSource) Schema(
	ctx context.Context,
	_ datasource.SchemaRequest,
	resp *datasource.SchemaResponse,
) {
	resp.Schema = ProvisioningLicenseDataSourceSchema(ctx)
}

func provisioningLicenseAsState(
	ctx context.Context,
	lic *sdk.GetProvisioningLicense200ResponseLicense,
) (ProvisioningLicenseModel, error) {
	var licenseType LicenseTypeValue
	if lic.LicenseType != nil {
		lt, diags := NewLicenseTypeValue(
			LicenseTypeValue{}.AttributeTypes(ctx),
			map[string]attr.Value{
				"id":   convert.Int64ToType(lic.LicenseType.Id),
				"code": convert.StrToType(lic.LicenseType.Code),
				"name": convert.StrToType(lic.LicenseType.Name),
			},
		)
		if diags.HasError() {
			return ProvisioningLicenseModel{}, fmt.Errorf("error creating license_type")
		}

		licenseType = lt
	} else {
		licenseType = NewLicenseTypeValueNull()
	}

	return ProvisioningLicenseModel{
		Id:               convert.Int64ToType(lic.Id),
		Name:             convert.StrToType(lic.Name),
		Description:      convert.StrToType(lic.Description),
		LicenseType:      licenseType,
		Copies:           convert.Int64ToType(lic.Copies),
		FullName:         convert.StrToType(lic.FullName),
		LicenseVersion:   convert.StrToType(lic.LicenseVersion),
		OrgName:          convert.StrToType(lic.OrgName),
		ReservationCount: convert.Int64ToType(lic.ReservationCount),
	}, nil
}

func getByID(
	ctx context.Context,
	id int64,
	apiClient *sdk.APIClient,
) (*sdk.GetProvisioningLicense200ResponseLicense, error) {
	r, hresp, err := apiClient.ProvisioningLicensesAPI.GetProvisioningLicense(ctx, id).Execute()
	if r == nil || err != nil || hresp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf(
			"GET failed for provisioning license %d: %s", id, providererrors.ErrMsg(err, hresp),
		)
	}

	if r.License == nil {
		return nil, fmt.Errorf(
			"GET failed for provisioning license %d: response missing license", id,
		)
	}

	lic := *r.License

	return &lic, nil
}

func getByName(
	ctx context.Context,
	name string,
	apiClient *sdk.APIClient,
) (*sdk.GetProvisioningLicense200ResponseLicense, error) {
	rs, hresp, err := apiClient.ProvisioningLicensesAPI.ListProvisioningLicenses(ctx).
		Name(name).Execute()
	if rs == nil || err != nil || hresp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf(
			"GET failed for provisioning license %s: %s", name, providererrors.ErrMsg(err, hresp),
		)
	}

	var matched []sdk.ListProvisioningLicenses200ResponseAllOfLicensesInner

	for _, o := range rs.Licenses {
		if o.Name != nil && *o.Name == name {
			matched = append(matched, o)
		}
	}

	if len(matched) == 0 {
		return nil, errors.New(ErrorNoLicenseFound)
	} else if len(matched) > 1 {
		return nil, errors.New(ErrorMultipleLicenses)
	}

	if matched[0].Id == nil {
		return nil, fmt.Errorf(
			"GET failed for provisioning license %s: response missing id", name,
		)
	}

	return getByID(ctx, *matched[0].Id, apiClient)
}

func getProvisioningLicense(
	ctx context.Context,
	config *ProvisioningLicenseModel,
	apiClient *sdk.APIClient,
) (*sdk.GetProvisioningLicense200ResponseLicense, error) {
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
	var config ProvisioningLicenseModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	apiClient, err := d.NewClient(ctx)
	if err != nil {
		resp.Diagnostics.AddError(summary, "could not create sdk client")

		return
	}

	lic, err := getProvisioningLicense(ctx, &config, apiClient)
	if err != nil {
		resp.Diagnostics.AddError(summary, err.Error())

		return
	}

	state, err := provisioningLicenseAsState(ctx, lic)
	if err != nil {
		resp.Diagnostics.AddError(summary, err.Error())

		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
