// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package plan

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
	summary                 = "read plan data source"
	ErrorNoPlanFound        = `no group found`
	ErrorNoValidSearchTerms = `no valid search terms - an id or name is required`
	ErrorRunningPreApply    = `Error running pre-apply plan: exit status 1`
	ErrorMultiplePlans      = `multiple groups were returned`
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource = &DataSource{}
)

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
	resp.TypeName = req.ProviderTypeName + "_morpheus_plan"
}

// Schema defines the schema for the data source.
func (d *DataSource) Schema(
	ctx context.Context,
	_ datasource.SchemaRequest,
	resp *datasource.SchemaResponse,
) {
	resp.Schema = PlanDataSourceSchema(ctx)
}

func GetPlanByID(
	ctx context.Context,
	id int64,
	apiClient *sdk.APIClient,
) (*sdk.GetServicePlans200ResponseServicePlan, error) {
	p, hresp, err := apiClient.ServicePlansAPI.GetServicePlans(ctx, id).Execute()
	if p == nil || err != nil || hresp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GET failed for plan %d: %s", id, internalErrors.ErrMsg(err, hresp))
	}

	plan, ok := p.GetServicePlanOk()

	if !ok {
		return nil, fmt.Errorf("plan %d is nil", id)
	}

	return plan, nil
}

func getPlanByName(
	ctx context.Context,
	name string,
	provisionType string,
	apiClient *sdk.APIClient,
) (*sdk.GetServicePlans200ResponseServicePlan, error) {
	pTypes, hresp, err := apiClient.ProvisioningAPI.ListProvisionTypes(ctx).Name(
		provisionType).Execute()
	if pTypes == nil || err != nil || hresp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GET failed for plan , provision_type %s: %s",
			provisionType, internalErrors.ErrMsg(err, hresp))
	}

	var matchingProvisionTypes []sdk.
		GetInstanceTypeProvisioning200ResponseAllOfInstanceTypeInstanceTypeLayoutsInnerProvisionType
	for _, pT := range pTypes.GetProvisionTypes() {
		if pTName, ok := pT.GetNameOk(); ok && *pTName == provisionType {
			matchingProvisionTypes = append(matchingProvisionTypes, pT)
		}
	}

	if len(matchingProvisionTypes) == 0 {
		return nil, fmt.Errorf("provision_type %s not found", provisionType)
	}

	if len(matchingProvisionTypes) > 1 {
		return nil, fmt.Errorf("multiple provision_type with name %s found", provisionType)
	}

	pTypeID, ok := matchingProvisionTypes[0].GetIdOk()
	if !ok {
		return nil, fmt.Errorf("provision_type %s id not found", provisionType)
	}

	ps, hresp, err := apiClient.ServicePlansAPI.ListServicePlans(ctx).Name(
		name).ProvisionTypeId(*pTypeID).Execute()
	if ps == nil || err != nil || hresp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GET failed for plan %s: %s", name, internalErrors.ErrMsg(err, hresp))
	}

	var matchingServicePlans []sdk.ListServicePlans200ResponseAllOfServicePlansInner

	for _, p := range ps.GetServicePlans() {
		if pName, pNameOk := p.GetNameOk(); pNameOk {
			if pProvisionType, pProvisionTypeOk := p.GetProvisionTypeOk(); pProvisionTypeOk {
				// now check name and ProvisionType match getplanByName() params
				if *pName == name && pProvisionType.GetName() == provisionType {
					matchingServicePlans = append(matchingServicePlans, p)
				}
			}
		}
	}

	if len(matchingServicePlans) == 1 {
		if pID, pIDOk := matchingServicePlans[0].GetIdOk(); pIDOk {
			// same return types as GetPlanByID
			return GetPlanByID(ctx, *pID, apiClient)
		}

		return nil, fmt.Errorf("plan %s, id not found", name)
	} else if len(matchingServicePlans) > 1 {
		return nil, errors.New(ErrorMultiplePlans)
	}

	return nil, errors.New(ErrorNoPlanFound)
}

func getPlan(
	ctx context.Context,
	data PlanModel,
	apiClient *sdk.APIClient,
) (*sdk.GetServicePlans200ResponseServicePlan, error) {
	if !data.Id.IsNull() {
		return GetPlanByID(ctx, data.Id.ValueInt64(), apiClient)
	} else if !data.Name.IsNull() && !data.ProvisionType.IsNull() {
		return getPlanByName(ctx, data.Name.ValueString(), data.ProvisionType.ValueString(), apiClient)
	}

	return nil, errors.New(ErrorNoValidSearchTerms)
}

// Read refreshes the Terraform state with the latest data.
func (d *DataSource) Read(
	ctx context.Context,
	req datasource.ReadRequest,
	resp *datasource.ReadResponse,
) {
	var data PlanModel

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

	plan, err := getPlan(ctx, data, apiClient)
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
	planProvisionType := plan.ProvisionType.GetName()
	data.ProvisionType = convert.StrToType(&planProvisionType)

	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}
