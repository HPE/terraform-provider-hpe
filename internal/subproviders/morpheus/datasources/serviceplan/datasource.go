// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package serviceplan

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/HewlettPackard/hpe-morpheus-go-sdk/sdk"
	"github.com/hashicorp/terraform-plugin-framework/datasource"

	"github.com/HPE/terraform-provider-hpe/internal/subproviders/morpheus/configure"
	"github.com/HPE/terraform-provider-hpe/internal/subproviders/morpheus/convert"
	internalErrors "github.com/HPE/terraform-provider-hpe/internal/subproviders/morpheus/errors"
)

const (
	summary                 = "read service plan data source"
	ErrorNoServicePlanFound = `no service plan found`
	ErrorNoValidSearchTerms = "no valid search terms - an id or (name and provision_type_code) " +
		"is required"
	ErrorRunningPreApply      = `Error running pre-apply plan: exit status 1`
	ErrorMultipleServicePlans = `multiple service plans were returned`
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
	resp.TypeName = req.ProviderTypeName + "_morpheus_service_plan"
}

// Schema defines the schema for the data source.
func (d *DataSource) Schema(
	ctx context.Context,
	_ datasource.SchemaRequest,
	resp *datasource.SchemaResponse,
) {
	resp.Schema = ServicePlanDataSourceSchema(ctx)
}

func getServicePlanByID(
	ctx context.Context,
	id int64,
	apiClient *sdk.APIClient,
) (*sdk.GetServicePlans200ResponseServicePlan, error) {
	sp, hresp, err := apiClient.ServicePlansAPI.GetServicePlans(ctx, id).Execute()
	if sp == nil || err != nil || hresp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf(
			"GET failed for service plan %d: %s", id, internalErrors.ErrMsg(err, hresp))
	}

	servicePlan, ok := sp.GetServicePlanOk()

	if !ok {
		return nil, fmt.Errorf("service plan %d is nil", id)
	}

	return servicePlan, nil
}

func getServicePlanByName(
	ctx context.Context,
	name string,
	provisionTypeCode string,
	apiClient *sdk.APIClient,
) (*sdk.GetServicePlans200ResponseServicePlan, error) {
	pTypes, hresp, err := apiClient.ProvisioningAPI.ListProvisionTypes(ctx).Code(
		provisionTypeCode).Execute()
	if pTypes == nil || err != nil || hresp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GET failed for service plan , provision type code %s: %s",
			provisionTypeCode, internalErrors.ErrMsg(err, hresp))
	}

	var matchingProvisionTypes []sdk.
		GetInstanceTypeProvisioning200ResponseAllOfInstanceTypeInstanceTypeLayoutsInnerProvisionType
	for _, pt := range pTypes.GetProvisionTypes() {
		if ptCode, ok := pt.GetCodeOk(); ok && *ptCode == provisionTypeCode {
			matchingProvisionTypes = append(matchingProvisionTypes, pt)
		}
	}

	if len(matchingProvisionTypes) == 0 {
		return nil, fmt.Errorf("provision type with code %s not found", provisionTypeCode)
	}

	if len(matchingProvisionTypes) > 1 {
		return nil, fmt.Errorf("multiple provision types with code %s found", provisionTypeCode)
	}

	pTypeID, ok := matchingProvisionTypes[0].GetIdOk()
	if !ok {
		return nil, fmt.Errorf("id not found for provision type with code %s", provisionTypeCode)
	}

	ps, hresp, err := apiClient.ServicePlansAPI.ListServicePlans(ctx).Name(
		name).ProvisionTypeId(*pTypeID).Execute()
	if ps == nil || err != nil || hresp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf(
			"GET failed for service_plan %s: %s", name, internalErrors.ErrMsg(err, hresp))
	}

	var matchingServicePlans []sdk.ListServicePlans200ResponseAllOfServicePlansInner

	for _, sp := range ps.GetServicePlans() {
		if pName, pNameOk := sp.GetNameOk(); pNameOk {
			if pProvisionType, pProvisionTypeOk := sp.GetProvisionTypeOk(); pProvisionTypeOk {
				// now check name and ProvisionType match getplanByName() params
				if *pName == name && pProvisionType.GetCode() == provisionTypeCode {
					matchingServicePlans = append(matchingServicePlans, sp)
				}
			}
		}
	}
	if len(matchingServicePlans) == 1 {
		if pID, pIDOk := matchingServicePlans[0].GetIdOk(); pIDOk {
			// same return types as GetPlanByID
			return getServicePlanByID(ctx, *pID, apiClient)
		}

		return nil, fmt.Errorf("service plan %s, id not found", name)
	} else if len(matchingServicePlans) > 1 {
		return nil, errors.New(ErrorMultipleServicePlans)
	}

	return nil, errors.New(ErrorNoServicePlanFound)
}

func getServicePlan(
	ctx context.Context,
	data ServicePlanModel,
	apiClient *sdk.APIClient,
) (*sdk.GetServicePlans200ResponseServicePlan, error) {
	if !data.Id.IsNull() {
		return getServicePlanByID(ctx, data.Id.ValueInt64(), apiClient)
	} else if !data.Name.IsNull() && !data.ProvisionTypeCode.IsNull() {
		return getServicePlanByName(
			ctx, data.Name.ValueString(), data.ProvisionTypeCode.ValueString(), apiClient)
	}

	return nil, errors.New(ErrorNoValidSearchTerms)
}

// Read refreshes the Terraform state with the latest data.
func (d *DataSource) Read(
	ctx context.Context,
	req datasource.ReadRequest,
	resp *datasource.ReadResponse,
) {
	var data ServicePlanModel

	// Read config
	diags := req.Config.Get(ctx, &data)
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

	plan, err := getServicePlan(ctx, data, apiClient)
	if err != nil {
		resp.Diagnostics.AddError(
			summary,
			err.Error(),
		)

		return
	}

	data.Id = convert.Int64ToType(plan.Id)
	data.Name = convert.StrToType(plan.Name)
	data.Code = convert.StrToType(plan.Code)
	data.Description = convert.StrToType(plan.Description)
	planProvisionType := plan.ProvisionType.GetCode()
	data.ProvisionTypeCode = convert.StrToType(&planProvisionType)

	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}
