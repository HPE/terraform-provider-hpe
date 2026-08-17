// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package instancestoragecontroller

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/types"

	sdk "github.com/HPE/terraform-provider-hpe/internal/sdk/oapigen"

	"github.com/HPE/terraform-provider-hpe/morpheus/configure"
	providererrors "github.com/HPE/terraform-provider-hpe/morpheus/utils/errfmt"
	"github.com/HPE/terraform-provider-hpe/utils/convert"
)

const (
	summary                                 = "read instance storage controller data source"
	ErrorNoInstanceStorageControllerFound   = `no instance storage controller found`
	ErrorMultipleInstanceStorageControllers = `multiple instance storage controllers were returned`

	// newControllerID is the id component of the mount point. Morpheus uses -1
	// to mean "a new controller": it resolves an existing controller by matching
	// busNumber + typeId, or creates one. See StorageVolume.getControllerMountPoint
	// and AbstractBoxProvisionService.assignStorageVolumeController in morpheus-ui.
	newControllerID = -1
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
	resp.TypeName = req.ProviderTypeName + "_" + "instance_storage_controller"
}

// Schema defines the schema for the data source.
func (d *DataSource) Schema(
	ctx context.Context,
	_ datasource.SchemaRequest,
	resp *datasource.SchemaResponse,
) {
	resp.Schema = InstanceStorageControllerDataSourceSchema(ctx)
}

// buildControllerMountPoint composes a controller_mount_point in the Morpheus
// format id:busNumber:typeId:unitNumber. The id is always -1, meaning "new
// controller" (Morpheus resolves an existing controller by matching busNumber +
// typeId, or creates one). This matches hpegl exactly and, being a pure function
// of its inputs, is stable across plan and apply.
func buildControllerMountPoint(busNumber, typeID, unitNumber int64) string {
	return fmt.Sprintf("%d:%d:%d:%d", newControllerID, busNumber, typeID, unitNumber)
}

// normalizeControllerName trims surrounding whitespace and lowercases, matching
// hpegl (strings.TrimSpace(strings.ToLower(...))) so ported configurations
// resolve identically.
func normalizeControllerName(s string) string {
	return strings.TrimSpace(strings.ToLower(s))
}

// matchedControllerType is the resolved instance storage controller. Fields are
// held as pointers so a value the API omits stays null in state (mapped via the
// convert helpers) rather than being flattened to a zero value.
type matchedControllerType struct {
	id           *int64
	category     *string
	maxDevices   *int64
	displayOrder *int64
	enabled      *bool
	creatable    *bool
}

// matchControllerType finds the single controller type whose name matches the
// requested name (case-insensitive, whitespace-trimmed). It errors on zero or
// more than one match, so the data source fails clearly rather than silently
// picking an arbitrary one. It is a pure function so the match/error logic is
// unit testable without an appliance.
func matchControllerType(
	controllerTypes []sdk.ListProvisionTypes200ResponseAllOfProvisionTypesInnerControllerTypesInner,
	name string,
) (matchedControllerType, error) {
	target := normalizeControllerName(name)

	var matches []matchedControllerType

	for _, ct := range controllerTypes {
		if ct.Name == nil || normalizeControllerName(*ct.Name) != target {
			continue
		}

		matches = append(matches, matchedControllerType{
			id:           ct.Id,
			category:     ct.Category,
			maxDevices:   ct.MaxDevices,
			displayOrder: ct.DisplayOrder,
			enabled:      ct.Enabled,
			creatable:    ct.Creatable,
		})
	}

	switch len(matches) {
	case 0:
		return matchedControllerType{}, errors.New(ErrorNoInstanceStorageControllerFound)
	case 1:
		return matches[0], nil
	default:
		return matchedControllerType{}, errors.New(ErrorMultipleInstanceStorageControllers)
	}
}

// getControllerTypes resolves a provision type code (e.g. "vmware") to its
// controller types, requiring exactly one matching provision type. The Morpheus
// option source for storageControllerTypes resolves layout -> provisionType and
// reads provisionType.controllerTypes, so reading via /api/provision-types by
// code yields the same list; layout scoping is incidental.
func getControllerTypes(
	ctx context.Context,
	code string,
	apiClient *sdk.APIClient,
) ([]sdk.ListProvisionTypes200ResponseAllOfProvisionTypesInnerControllerTypesInner, error) {
	pTypes, hresp, err := apiClient.ProvisioningAPI.ListProvisionTypes(ctx).Code(code).Execute()
	if pTypes == nil || err != nil || hresp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GET failed for provision type code %s: %s",
			code, providererrors.ErrMsg(err, hresp))
	}

	var matching []sdk.ListProvisionTypes200ResponseAllOfProvisionTypesInner

	for _, pt := range pTypes.ProvisionTypes {
		if pt.Code != nil && *pt.Code == code {
			matching = append(matching, pt)
		}
	}

	switch {
	case len(matching) == 0:
		return nil, fmt.Errorf("provision type with code %s not found", code)
	case len(matching) > 1:
		return nil, fmt.Errorf("multiple provision types with code %s found", code)
	default:
		return matching[0].ControllerTypes, nil
	}
}

// Read refreshes the Terraform state with the latest data.
func (d *DataSource) Read(
	ctx context.Context,
	req datasource.ReadRequest,
	resp *datasource.ReadResponse,
) {
	var config InstanceStorageControllerModel

	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	apiClient, err := d.NewClient(ctx)
	if err != nil {
		resp.Diagnostics.AddError(summary, "could not create sdk client")

		return
	}

	// Data source schemas cannot declare attribute defaults, so apply the
	// defaults (hpegl's vmware provision type, unit number 0) here.
	provisionTypeCode := defaultProvisionTypeCode
	if !config.ProvisionTypeCode.IsNull() && !config.ProvisionTypeCode.IsUnknown() {
		provisionTypeCode = config.ProvisionTypeCode.ValueString()
	}

	interfaceNumber := int64(0)
	if !config.InterfaceNumber.IsNull() && !config.InterfaceNumber.IsUnknown() {
		interfaceNumber = config.InterfaceNumber.ValueInt64()
	}

	controllerTypes, err := getControllerTypes(ctx, provisionTypeCode, apiClient)
	if err != nil {
		resp.Diagnostics.AddError(summary, err.Error())

		return
	}

	match, err := matchControllerType(controllerTypes, config.ControllerName.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(summary, err.Error())

		return
	}

	// The controller type id is required to compose the mount point; a match
	// without one cannot produce a usable value.
	if match.id == nil {
		resp.Diagnostics.AddError(summary, "matched instance storage controller has no id")

		return
	}

	mountPoint := buildControllerMountPoint(config.BusNumber.ValueInt64(), *match.id, interfaceNumber)

	state := InstanceStorageControllerModel{
		ControllerName:       config.ControllerName,
		BusNumber:            config.BusNumber,
		InterfaceNumber:      types.Int64Value(interfaceNumber),
		ProvisionTypeCode:    types.StringValue(provisionTypeCode),
		Id:                   convert.Int64ToType(match.id),
		ControllerMountPoint: types.StringValue(mountPoint),
		Category:             convert.StrToType(match.category),
		MaxDevices:           convert.Int64ToType(match.maxDevices),
		DisplayOrder:         convert.Int64ToType(match.displayOrder),
		Enabled:              convert.BoolToType(match.enabled),
		Creatable:            convert.BoolToType(match.creatable),
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
