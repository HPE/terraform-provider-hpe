// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package datastore

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	sdk "github.com/HPE/terraform-provider-hpe/internal/sdk/oapigen"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/HPE/terraform-provider-hpe/morpheus/configure"
	"github.com/HPE/terraform-provider-hpe/morpheus/utils/errfmt"
	"github.com/HPE/terraform-provider-hpe/utils/compare"
	"github.com/HPE/terraform-provider-hpe/utils/convert"
)

const (
	// listMax bounds a by-name lookup. name is an exact server-side filter, so
	// this is only reached when many datastores share one name; the default
	// page of 25 would silently hide the rest.
	listMax = 250

	nfsDatastoreCode   = "libvirt-netfs-nfs"
	alletraMPHVMCode   = "hpedatastore-alletra-mp"
	alletraMPBmaasCode = "hpedatastore-alletra-mp-bmaas"
	gfs2DatastoreCode  = "libvirt-dir-gfs2"

	cloudRefType   = "ComputeZone"
	clusterRefType = "ComputeServerGroup"

	associatedResourceTypeCluster = "Cluster"
	associatedResourceTypeCloud   = "Cloud"
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
	resp.TypeName = req.ProviderTypeName + "_" + "datastore"
}

// Schema defines the schema for the data source.
func (d *DataSource) Schema(
	ctx context.Context,
	_ datasource.SchemaRequest,
	resp *datasource.SchemaResponse,
) {
	resp.Schema = DatastoreDataSourceSchema(ctx)
}

func getDatastoreById(
	ctx context.Context,
	id int64,
	client *sdk.APIClient,
) (*DatastoreModel, error) {
	response, hresp, err := client.DatastoresAPI.GetDatastores(ctx, id).Execute()
	if err != nil || hresp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("datastore %d GET failed: %s", id, errfmt.ErrMsg(err, hresp))
	}

	datastore := response.Datastore
	if datastore == nil {
		return nil, fmt.Errorf("datastore %d is nil", id)
	}

	return datastoreToModel(ctx, datastore, id, client)
}

// datastoreToModel maps a datastore onto the data source model.
//
// Split out from getDatastoreById so that the by-name path can reuse it. That
// path already holds the datastore from the listing and does not fetch it
// again; see getDatastoreByName.
func datastoreToModel(
	ctx context.Context,
	datastore *sdk.GetDatastores200ResponseAllOfDatastore,
	id int64,
	client *sdk.APIClient,
) (*DatastoreModel, error) {
	state := &DatastoreModel{}

	state.Id = types.Int64Value(id)
	state.Name = convert.StrToType(&datastore.Name)
	state.Type = convert.StrToType(&datastore.Type)
	state.Code = convert.StrToType(datastore.Code.Get())
	state.Visibility = convert.StrToType(datastore.Visibility)
	state.Status = convert.StrToType(&datastore.Status)
	state.Active = convert.BoolToType(datastore.Active)
	state.ExternalId = convert.StrToType(datastore.ExternalId)
	state.ExternalType = convert.StrToType(datastore.ExternalType)
	state.StorageSize = convert.Int64ToType(datastore.StorageSize.Get())
	state.ExternalPath = convert.StrToType(datastore.ExternalPath)
	state.FreeSpace = convert.Int64ToType(datastore.FreeSpace.Get())
	state.HeartBeatTarget = convert.BoolToType(datastore.HeartBeatTarget)
	state.Online = convert.BoolToType(datastore.Online)
	state.StatusMessage = convert.StrToType(datastore.StatusMessage)
	state.AllowProvision = convert.BoolToType(datastore.AllowProvision)
	state.AllowWrite = convert.BoolToType(datastore.AllowWrite)
	state.AllowRead = convert.BoolToType(datastore.AllowRead)
	state.DrsEnabled = convert.BoolToType(datastore.DrsEnabled)
	state.DefaultStore = convert.BoolToType(datastore.DefaultStore)

	// Set DatastoreType, we get the code and id back from the API
	datastoreType := datastore.DatastoreType
	datastoreTypeValue := DatastoreTypeValue{}
	datastoreTypeValue.Id = convert.Int64ToType(&datastoreType.Id)
	datastoreTypeValue.Code = convert.StrToType(&datastoreType.Code)
	datastoreTypeValue.state = attr.ValueStateKnown
	state.DatastoreType = datastoreTypeValue

	// associated_resource_type, associated_resource_id, tenants and
	// resource_permissions all derive from refType and refId.
	//
	// Not every appliance reports them. Some answer with the datastore's own
	// fields and nothing about what it belongs to, and a datastore that exists
	// and can be named is still worth returning: the common use of this data
	// source is to resolve a name to an id. So a missing association leaves
	// those four attributes null rather than failing the read, which is how the
	// plural hpe_morpheus_datastores data source has always treated it.
	//
	// An unrecognised refType is treated the same way. It means the API knows a
	// kind of association this provider does not, which is not a reason to
	// refuse the datastore.
	refType := datastore.RefType
	refId := datastore.RefId

	// Typed nulls, not zero values. TenantsType{} on its own carries no
	// attribute types, so the set it produces is Object[] rather than the
	// object the schema declares, and Terraform rejects that as a provider bug
	// instead of reading it as an absent value. Deriving the type from the
	// value keeps it in step with the schema.
	state.Tenants = types.SetNull(TenantsValue{}.Type(ctx))
	state.ResourcePermissions = NewResourcePermissionsValueNull()

	if refType != nil && refId != nil {
		state.AssociatedResourceId = convert.Int64ToType(refId)

		switch *refType {
		case cloudRefType:
			state.AssociatedResourceType = types.StringValue(associatedResourceTypeCloud)
			tenants, resourcePermissions, err := populateCloudDatastoreInformation(
				ctx, id, *refId, client)
			if err != nil {
				return nil, err
			}
			state.Tenants = tenants
			state.ResourcePermissions = resourcePermissions

		case clusterRefType:
			state.AssociatedResourceType = types.StringValue(associatedResourceTypeCluster)
			tenants, resourcePermissions, err := populateClusterDatastoreInformation(
				ctx, id, *refId, client)
			if err != nil {
				return nil, err
			}
			state.Tenants = tenants
			state.ResourcePermissions = resourcePermissions

		default:
			tflog.Debug(ctx, "datastore has an unrecognised ref type", map[string]any{
				"datastore_id": id,
				"ref_type":     *refType,
			})
		}
	} else {
		tflog.Debug(ctx, "datastore reports no association; "+
			"associated_resource_type, associated_resource_id, tenants and "+
			"resource_permissions will be null", map[string]any{"datastore_id": id})
	}

	// Populate Config
	switch datastoreType.Code {
	case nfsDatastoreCode:
		var configNfsValue ConfigNfsValue
		for k, v := range datastore.Config {
			switch k {
			case "sourceDirPath":
				str := v.(string)
				configNfsValue.SourceDirPath = convert.StrToType(&str)
			case "sourceHostname":
				str := v.(string)
				configNfsValue.SourceHostname = convert.StrToType(&str)
			}
		}
		configNfsValue.state = attr.ValueStateKnown
		state.ConfigNfs = configNfsValue

	case alletraMPHVMCode:
		var configAlletraMPHVM ConfigAlletrampHvmValue
		for k, v := range datastore.Config {
			switch k {
			case "enableransomware":
				switch t := v.(type) {
				case bool:
					configAlletraMPHVM.EnableRansomware = convert.BoolToType(&t)
				case string:
					configAlletraMPHVM.EnableRansomware = convert.StringToBool(ctx, t)
				}
			case "protocolType":
				str := v.(string)
				configAlletraMPHVM.ProtocolType = convert.StrToType(&str)
			}
		}
		configAlletraMPHVM.state = attr.ValueStateKnown
		state.ConfigAlletrampHvm = configAlletraMPHVM

		// The Alletra MP HVM plugin returns extra values in config other than those set in the API POST
		// We put these extra values into "config_from_api"
		keysMap := map[string]string{
			"enable_ransomware": "enableransomware",
			"protocol_type":     "protocolType",
		}
		keysFromMap, _ := compare.CheckPlanAttributeAgainstAPIAttribute(
			ctx, configAlletraMPHVM, datastore.Config, keysMap)

		configFromAPI, _ := createConfigFromApiDynamic(ctx, datastore.Config, keysFromMap)
		state.ConfigFromApi = configFromAPI

	case alletraMPBmaasCode:
		var configAlletraMPBmaas ConfigAlletrampBmaasValue
		for k, v := range datastore.Config {
			if k == "protocolType" {
				if str, ok := v.(string); ok {
					configAlletraMPBmaas.ProtocolType = convert.StrToType(&str)
				}
			}
		}
		configAlletraMPBmaas.state = attr.ValueStateKnown
		state.ConfigAlletrampBmaas = configAlletraMPBmaas

		// The Alletra MP BMaaS plugin returns extra values in config other than those set in the API POST
		// We put these extra values into "config_from_api"
		keysMap := map[string]string{
			"protocol_type": "protocolType",
		}
		keysFromMap, _ := compare.CheckPlanAttributeAgainstAPIAttribute(
			ctx, configAlletraMPBmaas, datastore.Config, keysMap)

		configFromAPI, _ := createConfigFromApiDynamic(ctx, datastore.Config, keysFromMap)
		state.ConfigFromApi = configFromAPI

	case gfs2DatastoreCode:
		var configGfs2 ConfigGfs2Value
		for k, v := range datastore.Config {
			switch k {
			case "blockDevice":
				str := v.(string)
				configGfs2.BlockDevice = convert.StrToType(&str)
			}
		}
		configGfs2.state = attr.ValueStateKnown
		state.ConfigGfs2 = configGfs2

	default:
		// Generic config
		config, _ := convert.MapToDynamic(ctx, datastore.Config)
		state.Config = config
	}

	// Set StorageServer to that returned by the API
	if server := datastore.StorageServer; server != nil {
		storageServer := StorageServerValue{}
		storageServer.Id = convert.Int64ToType(server.Id)
		storageServer.state = attr.ValueStateKnown
		state.StorageServer = storageServer
	}

	state.Cloud = NewCloudValueNull()
	if datastore.Zone != nil && datastore.Zone.Id != nil {
		cloud := CloudValue{}
		cloud.Id = convert.Int64ToType(datastore.Zone.Id)
		cloud.state = attr.ValueStateKnown
		state.Cloud = cloud
	}

	state.ResourcePool = NewResourcePoolValueNull()
	if datastore.ZonePool != nil && datastore.ZonePool.Id != nil {
		resourcePool := ResourcePoolValue{}
		resourcePool.Id = convert.Int64ToType(datastore.ZonePool.Id)
		resourcePool.state = attr.ValueStateKnown
		state.ResourcePool = resourcePool
	}

	state.Owner = NewOwnerValueNull()
	if datastore.Owner != nil && datastore.Owner.Id != nil {
		owner := OwnerValue{}
		owner.Id = convert.Int64ToType(datastore.Owner.Id)
		owner.state = attr.ValueStateKnown
		state.Owner = owner
	}

	state.Locations = basetypes.NewSetNull(LocationsValue{}.Type(ctx))
	if locations := datastore.Locations; locations != nil {
		locationsSet, d := convert.ToSetType(
			ctx,
			locations,
			func(
				in sdk.GetDatastores200ResponseAllOfDatastoreLocationsInner,
			) LocationsValue {
				return LocationsValue{
					RefId:         convert.Int64ToType(in.RefId),
					RefType:       convert.StrToType(in.RefType),
					Status:        convert.StrToType(in.Status),
					StatusMessage: convert.StrToType(in.StatusMessage),
					state:         attr.ValueStateKnown,
				}
			},
		)
		if d.HasError() {
			return nil, fmt.Errorf("datastore %d error in creating locations set", id)
		}

		state.Locations = locationsSet
	}

	state.Datastores = basetypes.NewSetNull(DatastoresValue{}.Type(ctx))
	if datastores := datastore.Datastores; datastores != nil {
		datastoresSet, d := convert.ToSetType(
			ctx,
			datastores,
			func(
				in sdk.GetDatastores200ResponseAllOfDatastoreDatastoresInner,
			) DatastoresValue {
				var datastoreID *int64
				if in.Id != nil {
					id64 := int64(*in.Id)
					datastoreID = &id64
				}

				return DatastoresValue{
					Id:    convert.Int64ToType(datastoreID),
					Name:  convert.StrToType(in.Name),
					state: attr.ValueStateKnown,
				}
			},
		)
		if d.HasError() {
			return nil, fmt.Errorf("datastore %d error in creating datastores set", id)
		}
		state.Datastores = datastoresSet
	}

	return state, nil
}

// populateCloudDatastoreInformation gets the cloud datastore information for a datastore
// and populates the Tenants and ResourcePermissions fields of the datastore resource model
func populateCloudDatastoreInformation(
	ctx context.Context,
	id, cloudId int64,
	client *sdk.APIClient,
) (types.Set, ResourcePermissionsValue, error) {
	// Get cloud datastore information
	cdResp, hresp, err := client.CloudsAPI.GetCloudDatastores(ctx, cloudId, id).Execute()
	if err != nil {
		retErr := fmt.Errorf("datastore %d cloud datastore GET failed: %s", id, errfmt.ErrMsg(err, hresp))

		return types.SetNull(TenantsType{}), NewResourcePermissionsValueNull(), retErr
	}

	cloudDatastore := cdResp.Datastore
	if cloudDatastore == nil {
		retErr := fmt.Errorf("datastore %d missing cloud datastore in response", id)

		return types.SetNull(TenantsType{}), NewResourcePermissionsValueNull(), retErr
	}

	// Populate Tenants
	var tenantsSet types.Set
	tenants := cloudDatastore.Tenants
	if tenants != nil {

		var tdiag diag.Diagnostics
		tenantsSet, tdiag = convert.ToSetType(
			ctx,
			tenants,
			func(
				in sdk.GetCloudDatastores200ResponseAllOfDatastoreTenantsInner,
			) TenantsValue {
				return TenantsValue{
					Id:    convert.Int64ToType(in.Id),
					state: attr.ValueStateKnown,
				}
			},
		)

		if tdiag.HasError() {
			retErr := fmt.Errorf("datastore %d error in creating tenants set", id)

			return types.SetNull(TenantsType{}), NewResourcePermissionsValueNull(), retErr
		}
	}

	// Populate ResourcePermissions, we'll only do Groups for now
	var resourcePermissions ResourcePermissionsValue
	rp := cloudDatastore.ResourcePermission
	if rp != nil {
		sites := rp.Sites
		if sites == nil {
			tflog.Debug(ctx, fmt.Sprintf("datastore %d missing resource permission sites in response", id))

			return tenantsSet, NewResourcePermissionsValueNull(), nil
		}
		groupsSet, gdiag := convert.ToSetType(
			ctx,
			sites,
			func(
				in sdk.GetCloudDatastores200ResponseAllOfDatastoreResourcePermissionSitesInner,
			) GroupsValue {
				return GroupsValue{
					Id:    convert.Int64ToType(in.Id),
					state: attr.ValueStateKnown,
				}
			},
		)

		if gdiag.HasError() {
			retErr := fmt.Errorf("datastore %d error in creating groups set", id)

			return tenantsSet, NewResourcePermissionsValueNull(), retErr
		}

		resourcePermissions, err = populateResourcePermissionsFromApi(ctx, id, groupsSet)
	}

	return tenantsSet, resourcePermissions, err
}

// populateClusterDatastoreInformation gets the cluster datastore information for a datastore
// and populates the Tenants and ResourcePermissions fields of the datastore resource model
func populateClusterDatastoreInformation(
	ctx context.Context,
	id, clusterId int64,
	client *sdk.APIClient,
) (types.Set, ResourcePermissionsValue, error) {
	// Get cluster datastore information
	cdResp, hresp, err := client.ClustersAPI.GetClusterDatastore(ctx, clusterId, id).Execute()
	if err != nil {
		retErr := fmt.Errorf("datastore %d cluster datastore GET failed: %s", id, errfmt.ErrMsg(err, hresp))

		return types.SetNull(TenantsType{}), NewResourcePermissionsValueNull(), retErr
	}
	clusterDatastore := cdResp.Datastore
	if clusterDatastore == nil {
		retErr := fmt.Errorf("datastore %d missing cluster datastore in response", id)

		return types.SetNull(TenantsType{}), NewResourcePermissionsValueNull(), retErr
	}

	// Populate Tenants
	var tenantsSet types.Set
	tenants := clusterDatastore.Tenants
	if tenants != nil {

		var tdiag diag.Diagnostics
		tenantsSet, tdiag = convert.ToSetType(
			ctx,
			tenants,
			func(
				in sdk.GetClusterDatastore200ResponseDatastoreTenantsInner,
			) TenantsValue {
				return TenantsValue{
					Id:    convert.Int64ToType(in.Id),
					state: attr.ValueStateKnown,
				}
			},
		)
		if tdiag.HasError() {
			retErr := fmt.Errorf("datastore %d error in creating tenants set", id)

			return tenantsSet, NewResourcePermissionsValueNull(), retErr
		}
	}

	// Populate ResourcePermissions, we'll only do Groups for now
	var resourcePermissions ResourcePermissionsValue
	rp := clusterDatastore.ResourcePermissions
	if rp != nil {

		sites := rp.Sites
		if sites == nil {
			tflog.Debug(ctx, fmt.Sprintf("datastore %d missing resource permission sites in response", id))

			return tenantsSet, NewResourcePermissionsValueNull(), nil
		}

		groupsSet, gdiag := convert.ToSetType(
			ctx,
			sites,
			func(
				in sdk.GetClusterDatastore200ResponseDatastoreResourcePermissionsSitesInner,
			) GroupsValue {
				return GroupsValue{
					Id:    convert.Int64ToType(in.Id),
					state: attr.ValueStateKnown,
				}
			},
		)
		if gdiag.HasError() {
			retErr := fmt.Errorf("datastore %d error in creating groups set", id)

			return tenantsSet, NewResourcePermissionsValueNull(), retErr
		}

		resourcePermissions, err = populateResourcePermissionsFromApi(ctx, id, groupsSet)
	}

	return tenantsSet, resourcePermissions, err
}

func createConfigFromApiDynamic(
	ctx context.Context,
	config, keysFromMap map[string]any,
) (types.Dynamic, error) {
	// We get more information back from the API, we put this in config_from_api
	apiConfigForConfigApi := make(map[string]any)
	for k := range keysFromMap {
		if _, ok := config[k]; ok {
			apiConfigForConfigApi[k] = config[k]
		}
	}

	configFromApi, err := convert.MapToDynamic(ctx, apiConfigForConfigApi)
	if err != nil {
		return types.DynamicNull(), err
	}

	return configFromApi, nil
}

// populateResourcePermissionsFromApi populates the ResourcePermissionsValue
// from the API returned groups set
// we only populate Groups for now
func populateResourcePermissionsFromApi(
	ctx context.Context,
	id int64,
	groupsSet basetypes.SetValue,
) (ResourcePermissionsValue, error) {
	var err error

	if !groupsSet.IsNull() {
		attrTypes := make(map[string]attr.Type)
		attrValues := make(map[string]attr.Value)

		attrTypes["groups"] = groupsSet.Type(ctx)
		attrValues["groups"] = groupsSet

		resourcePermissions, rdiags := NewResourcePermissionsValue(attrTypes, attrValues)
		if rdiags.HasError() {
			err = fmt.Errorf("datastore %d error in creating resource permissions", id)
		}

		return resourcePermissions, err
	}

	return NewResourcePermissionsValueNull(), err
}

func getDatastoreByName(
	ctx context.Context,
	name string,
	client *sdk.APIClient,
) (*DatastoreModel, error) {
	datastores, hresp, err := client.DatastoresAPI.ListDatastores(ctx).
		Name(name).
		Max(listMax).
		Execute()
	if err != nil || hresp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("datastore %s list failed: %s", name, errfmt.ErrMsg(err, hresp))
	}

	matchingDatastores := datastores.Datastores
	if len(matchingDatastores) == 0 {
		return nil, fmt.Errorf("datastore %s not found", name)
	}

	// name is an exact server-side filter, so meta.total counts the matches
	// rather than the whole collection. A total larger than the page means more
	// datastores share this name than were fetched, and the ones not fetched
	// cannot be reported. Paging them in would not help: the practitioner still
	// has to say which one they meant.
	if datastores.Meta != nil && datastores.Meta.Total != nil &&
		*datastores.Meta.Total > int64(len(matchingDatastores)) {
		return nil, fmt.Errorf(
			"%d datastores are named %s, more than the %d fetched. "+
				"Specify an ID instead",
			*datastores.Meta.Total, name, listMax,
		)
	}

	if len(matchingDatastores) > 1 {
		var datastoreIDs []string
		for _, n := range matchingDatastores {
			datastoreIDs = append(datastoreIDs, fmt.Sprintf("%d", n.Id))
		}

		return nil, fmt.Errorf(
			"multiple datastores found with name %s. Datastore IDs: %s. "+
				"Please specify an ID instead",
			name,
			strings.Join(datastoreIDs, ", "),
		)
	}

	entry := matchingDatastores[0]

	// Newer appliances answer this listing with the whole datastore, which is
	// worth using: it saves a request, and it avoids GET /api/data-stores/{id},
	// which on some appliances answers 404 for a cloud-associated datastore the
	// listing returns quite happily.
	//
	// Older ones answer with nothing but an id and a name. There the fetch is
	// worth attempting, because it carries fields the entry does not -- but it
	// is an improvement, not a requirement. If it fails, the entry the listing
	// already gave us is still a datastore, and resolving a name to an id is
	// what this data source is mostly asked for. Refusing to answer because a
	// second request failed would leave the data source unusable on exactly the
	// appliances that need it most.
	//
	// refType marks a full entry: it is absent from a thin one.
	if entry.RefType == nil {
		model, err := getDatastoreById(ctx, entry.Id, client)
		if err == nil {
			return model, nil
		}

		tflog.Debug(ctx, "datastore could not be fetched by id; "+
			"falling back to the listing entry", map[string]any{
			"datastore_id": entry.Id,
			"error":        err.Error(),
		})
	}

	datastore, err := datastoreFromListEntry(&entry)
	if err != nil {
		return nil, err
	}

	return datastoreToModel(ctx, datastore, entry.Id, client)
}

// datastoreFromListEntry converts a listing entry into the single-item shape.
//
// The two are generated from the same API object and carry the same fields;
// only the Go type names differ, nested types included. Re-encoding is used in
// preference to copying thirty-odd fields and seven nested structures by hand,
// which would be silently wrong the first time the SDK gained a field and
// nobody updated the copy.
func datastoreFromListEntry(
	in *sdk.ListDatastores200ResponseAllOfDatastoresInner,
) (*sdk.GetDatastores200ResponseAllOfDatastore, error) {
	encoded, err := json.Marshal(in)
	if err != nil {
		return nil, fmt.Errorf("datastore %d could not be encoded: %w", in.Id, err)
	}

	var out sdk.GetDatastores200ResponseAllOfDatastore
	if err := json.Unmarshal(encoded, &out); err != nil {
		return nil, fmt.Errorf("datastore %d could not be decoded: %w", in.Id, err)
	}

	return &out, nil
}

func getDatastore(
	ctx context.Context,
	config DatastoreModel,
	client *sdk.APIClient,
) (*DatastoreModel, error) {
	if !config.Id.IsNull() {
		return getDatastoreById(ctx, config.Id.ValueInt64(), client)
	}
	if !config.Name.IsNull() {
		return getDatastoreByName(ctx, config.Name.ValueString(), client)
	}

	return nil, fmt.Errorf("either id or name must be specified")
}

func (d *DataSource) Read(
	ctx context.Context,
	req datasource.ReadRequest,
	resp *datasource.ReadResponse,
) {
	var config DatastoreModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	client, err := d.NewClient(ctx)
	if err != nil {
		resp.Diagnostics.AddError(
			"read datastore data source",
			fmt.Sprintf("failed to create client: %s", err.Error()),
		)

		return
	}

	state, err := getDatastore(ctx, config, client)
	if err != nil {
		resp.Diagnostics.AddError(
			"read datastore data source",
			err.Error(),
		)

		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}
