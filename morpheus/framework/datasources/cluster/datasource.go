// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package cluster

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/HewlettPackard/hpe-morpheus-go-sdk/oapigen/sdk"

	"github.com/HPE/terraform-provider-hpe/morpheus/configure"
	"github.com/HPE/terraform-provider-hpe/morpheus/utils/errfmt"
	"github.com/HPE/terraform-provider-hpe/utils/convert"
)

const (
	summary = "read cluster data source"

	ErrorNoValidSearchTerms = "no valid search terms - an id or name is required"
	ErrorRunningPreApply    = "Error running pre-apply plan: exit status 1"
	ErrorNoClusterFound     = "no cluster found"
	ErrorMultipleClusters   = "multiple clusters were returned"
)

var _ datasource.DataSource = &DataSource{}

func NewDataSource() datasource.DataSource {
	return &DataSource{}
}

type DataSource struct {
	configure.DataSourceWithMorpheusConfigure
	datasource.DataSource
}

func (d *DataSource) Metadata(
	_ context.Context,
	req datasource.MetadataRequest,
	resp *datasource.MetadataResponse,
) {
	resp.TypeName = req.ProviderTypeName + "_morpheus_cluster"
}

func (d *DataSource) Schema(
	ctx context.Context,
	_ datasource.SchemaRequest,
	resp *datasource.SchemaResponse,
) {
	resp.Schema = ClusterDataSourceSchema(ctx)
}

func populateClusterData(
	ctx context.Context,
	data *ClusterModel,
	cluster *sdk.GetCluster200ResponseCluster,
) error {
	var diags diag.Diagnostics

	data.Id = convert.Int64ToType(cluster.Id)
	data.Uuid = convert.StrToType(cluster.Uuid)
	data.Name = convert.StrToType(cluster.Name)
	data.Category = convert.StrToType(cluster.Category.Get())
	data.Visibility = convert.StrToType(cluster.Visibility)
	data.Description = convert.StrToType(cluster.Description.Get())
	data.Location = convert.StrToType(cluster.Location.Get())
	data.Enabled = convert.BoolToType(cluster.Enabled)
	data.ServiceUrl = convert.StrToType(cluster.ServiceUrl.Get())
	data.ServiceHost = convert.StrToType(cluster.ServiceHost.Get())
	data.ServicePath = convert.StrToType(cluster.ServicePath.Get())
	data.ServiceHostname = convert.StrToType(cluster.ServiceHostname.Get())
	data.ServicePort = convert.Int64ToType(cluster.ServicePort)
	data.ServiceUsername = convert.StrToType(cluster.ServiceUsername.Get())
	data.ServicePassword = convert.StrToType(cluster.ServicePassword.Get())
	data.ServicePasswordHash = convert.StrToType(cluster.ServicePasswordHash.Get())
	data.ServiceToken = convert.StrToType(cluster.ServiceToken.Get())
	data.ServiceTokenHash = convert.StrToType(cluster.ServiceTokenHash.Get())
	data.ServiceAccess = convert.StrToType(cluster.ServiceAccess.Get())
	data.ServiceAccessHash = convert.StrToType(cluster.ServiceAccessHash.Get())
	data.ServiceCert = convert.StrToType(cluster.ServiceCert.Get())
	data.ServiceCertHash = convert.StrToType(cluster.ServiceCertHash.Get())
	data.ServiceVersion = convert.StrToType(cluster.ServiceVersion.Get())
	data.SearchDomains = convert.StrToType(cluster.SearchDomains.Get())
	data.EnableInternalDns = convert.BoolToType(cluster.EnableInternalDns)
	data.InternalId = convert.StrToType(cluster.InternalId.Get())
	data.ExternalId = convert.StrToType(cluster.ExternalId.Get())
	data.DatacenterId = convert.StrToType(cluster.DatacenterId.Get())
	data.Status = convert.StrToType(cluster.Status)
	data.StatusDate = timeToType(cluster.StatusDate.Get())
	data.StatusMessage = convert.StrToType(cluster.StatusMessage.Get())
	data.InventoryLevel = convert.StrToType(cluster.InventoryLevel)
	data.LastSync = timeToType(cluster.LastSync.Get())
	data.NextRunDate = timeToType(cluster.NextRunDate.Get())
	data.LastSyncDuration = convert.Int64ToType(cluster.LastSyncDuration.Get())
	data.DateCreated = timeToType(cluster.DateCreated)
	data.LastUpdated = timeToType(cluster.LastUpdated)
	data.Managed = convert.BoolToType(cluster.Managed)
	data.Labels = convert.StrSliceToSet(cluster.GetLabels())
	data.AutoRecoverPowerState = convert.BoolToType(cluster.AutoRecoverPowerState)
	data.UseAgent = convert.StrToType(cluster.UseAgent.Get())
	data.ProvisionComplete = convert.BoolToType(cluster.ProvisionComplete)
	data.ServiceEntry = convert.StrToType(cluster.ServiceEntry.Get())
	data.UserGroup = convert.StrToType(cluster.UserGroup.Get())
	data.ContainersCount = convert.Int64ToType(cluster.ContainersCount)
	data.DeploymentsCount = convert.Int64ToType(cluster.DeploymentsCount)
	data.PodsCount = convert.Int64ToType(cluster.PodsCount)
	data.JobsCount = convert.Int64ToType(cluster.JobsCount)
	data.VolumesCount = convert.Int64ToType(cluster.VolumesCount)
	data.NamespacesCount = convert.Int64ToType(cluster.NamespacesCount)
	data.WorkersCount = convert.Int64ToType(cluster.WorkersCount)
	data.ServicesCount = convert.Int64ToType(cluster.ServicesCount)

	config, err := convert.MapToDynamic(ctx, cluster.Config)
	if err != nil {
		return fmt.Errorf("cluster config mapping failed: %w", err)
	}
	data.Config = config

	data.CreatedBy = NewCreatedByValueNull()
	if cluster.CreatedBy != nil {
		data.CreatedBy = CreatedByValue{
			Id:       convert.Int64ToType(cluster.CreatedBy.Id),
			Username: convert.StrToType(cluster.CreatedBy.Username),
			state:    attr.ValueStateKnown,
		}
	}

	data.Group = NewGroupValueNull()
	if cluster.Site != nil {
		data.Group = GroupValue{
			Id:    convert.Int64ToType(cluster.Site.Id),
			Name:  convert.StrToType(cluster.Site.Name),
			state: attr.ValueStateKnown,
		}
	}

	data.Layout = NewLayoutValueNull()
	if cluster.Layout != nil {
		data.Layout = LayoutValue{
			Id:                convert.Int64ToType(cluster.Layout.Id),
			Name:              convert.StrToType(cluster.Layout.Name),
			ProvisionTypeCode: convert.StrToType(cluster.Layout.ProvisionTypeCode),
			state:             attr.ValueStateKnown,
		}
	}

	data.Owner = NewOwnerValueNull()
	if cluster.Owner != nil {
		data.Owner = OwnerValue{
			Id:    convert.Int64ToType(cluster.Owner.Id),
			Name:  convert.StrToType(cluster.Owner.Name),
			state: attr.ValueStateKnown,
		}
	}

	data.Type = NewTypeValueNull()
	if cluster.Type != nil {
		data.Type = TypeValue{
			Id:    convert.Int64ToType(cluster.Type.Id),
			Name:  convert.StrToType(cluster.Type.Name),
			state: attr.ValueStateKnown,
		}
	}

	data.WorkerStats = NewWorkerStatsValueNull()
	if cluster.WorkerStats != nil {
		data.WorkerStats = WorkerStatsValue{
			CpuUsage:     convert.NumToType(cluster.WorkerStats.CpuUsage),
			CpuUsageAvg:  convert.NumToType(cluster.WorkerStats.CpuUsageAvg),
			CpuUsagePeak: convert.NumToType(cluster.WorkerStats.CpuUsagePeak),
			MaxMemory:    convert.Int64ToType(cluster.WorkerStats.MaxMemory),
			MaxStorage:   convert.Int64ToType(cluster.WorkerStats.MaxStorage),
			UsedCpu:      convert.Int64ToType(cluster.WorkerStats.UsedCpu),
			UsedMemory:   convert.Int64ToType(cluster.WorkerStats.UsedMemory),
			UsedStorage:  convert.Int64ToType(cluster.WorkerStats.UsedStorage),
			state:        attr.ValueStateKnown,
		}
	}

	cloudValue, cloudDiags := buildCloudValue(ctx, cluster.Zone)
	diags.Append(cloudDiags...)
	data.Cloud = cloudValue

	permissionsValue, permissionsDiags := buildPermissionsValue(ctx, cluster.Permissions)
	diags.Append(permissionsDiags...)
	data.Permissions = permissionsValue

	servers, serversDiags := buildServersSet(ctx, cluster.Servers)
	diags.Append(serversDiags...)
	data.Servers = servers

	integrations, integrationDiags := convert.ToSetType(
		ctx,
		cluster.Integrations,
		func(_ map[string]any) IntegrationsValue {
			return IntegrationsValue{state: attr.ValueStateKnown}
		},
	)
	diags.Append(integrationDiags...)
	data.Integrations = integrations

	tenants, tenantDiags := convert.ToSetType(
		ctx,
		cluster.Accounts,
		func(_ map[string]any) TenantsValue {
			return TenantsValue{state: attr.ValueStateKnown}
		},
	)
	diags.Append(tenantDiags...)
	data.Tenants = tenants

	if diags.HasError() {
		return fmt.Errorf("cluster mapping failed: %v", diags)
	}

	return nil
}

func buildCloudValue(
	ctx context.Context,
	zone *sdk.GetCluster200ResponseClusterZone,
) (CloudValue, diag.Diagnostics) {
	diags := diag.Diagnostics{}
	if zone == nil {
		return NewCloudValueNull(), diags
	}

	cloudTypeObj, cloudTypeDiags := NewCloudTypeValueNull().ToObjectValue(ctx)
	diags.Append(cloudTypeDiags...)
	if zone.ZoneType != nil {
		cloudType := CloudTypeValue{
			Id:    convert.Int64ToType(zone.ZoneType.Id),
			state: attr.ValueStateKnown,
		}

		cloudTypeObj, cloudTypeDiags = cloudType.ToObjectValue(ctx)
		diags.Append(cloudTypeDiags...)
	}

	return CloudValue{
		CloudType: cloudTypeObj,
		Id:        convert.Int64ToType(zone.Id),
		Name:      convert.StrToType(zone.Name),
		state:     attr.ValueStateKnown,
	}, diags
}

func buildPermissionsValue(
	ctx context.Context,
	permissions *sdk.GetCluster200ResponseClusterPermissions,
) (PermissionsValue, diag.Diagnostics) {
	diags := diag.Diagnostics{}
	if permissions == nil {
		return NewPermissionsValueNull(), diags
	}

	resourcePermissionsObj, resourcePermissionsDiags := NewResourcePermissionsValueNull().ToObjectValue(ctx)
	diags.Append(resourcePermissionsDiags...)
	if permissions.ResourcePermissions != nil {
		tenantObj, tenantDiags := NewTenantValueNull().ToObjectValue(ctx)
		diags.Append(tenantDiags...)
		if permissions.ResourcePermissions.Account != nil {
			tenant := TenantValue{
				Id:    convert.Int64ToType(permissions.ResourcePermissions.Account.Id),
				state: attr.ValueStateKnown,
			}

			tenantObj, tenantDiags = tenant.ToObjectValue(ctx)
			diags.Append(tenantDiags...)
		}

		groups, groupsDiags := convert.ToSetType(
			ctx,
			permissions.ResourcePermissions.Sites,
			func(_ map[string]any) GroupsValue {
				return GroupsValue{state: attr.ValueStateKnown}
			},
		)
		diags.Append(groupsDiags...)

		plans, plansDiags := convert.ToSetType(
			ctx,
			permissions.ResourcePermissions.Plans,
			func(_ map[string]any) PlansValue {
				return PlansValue{state: attr.ValueStateKnown}
			},
		)
		diags.Append(plansDiags...)

		resourcePermissions := ResourcePermissionsValue{
			All:           convert.BoolToType(permissions.ResourcePermissions.All),
			AllPlans:      convert.BoolToType(permissions.ResourcePermissions.AllPlans),
			CanManage:     convert.BoolToType(permissions.ResourcePermissions.CanManage),
			DefaultStore:  convert.BoolToType(permissions.ResourcePermissions.DefaultStore),
			DefaultTarget: convert.BoolToType(permissions.ResourcePermissions.DefaultTarget),
			Groups:        groups,
			Plans:         plans,
			Tenant:        tenantObj,
			state:         attr.ValueStateKnown,
		}

		resourcePermissionsObj, resourcePermissionsDiags = resourcePermissions.ToObjectValue(ctx)
		diags.Append(resourcePermissionsDiags...)
	}

	resourcePoolObj, resourcePoolDiags := NewResourcePoolValueNull().ToObjectValue(ctx)
	diags.Append(resourcePoolDiags...)
	if permissions.ResourcePool != nil {
		resourcePool := ResourcePoolValue{
			Id:         convert.Int64ToType(permissions.ResourcePool.Id),
			Name:       convert.StrToType(permissions.ResourcePool.Name),
			Visibility: convert.StrToType(permissions.ResourcePool.Visibility),
			state:      attr.ValueStateKnown,
		}

		resourcePoolObj, resourcePoolDiags = resourcePool.ToObjectValue(ctx)
		diags.Append(resourcePoolDiags...)
	}

	return PermissionsValue{
		ResourcePermissions: resourcePermissionsObj,
		ResourcePool:        resourcePoolObj,
		state:               attr.ValueStateKnown,
	}, diags
}

func buildServersSet(
	ctx context.Context,
	servers []sdk.GetCluster200ResponseClusterServersInner,
) (types.Set, diag.Diagnostics) {
	var serverDiags diag.Diagnostics

	result, diags := convert.ToSetType(
		ctx,
		servers,
		func(in sdk.GetCluster200ResponseClusterServersInner) ServersValue {
			computeServerTypeObj, objectDiags := NewComputeServerTypeValueNull().ToObjectValue(ctx)
			serverDiags.Append(objectDiags...)
			if in.ComputeServerType != nil {
				computeServerType := ComputeServerTypeValue{
					Code:     convert.StrToType(in.ComputeServerType.Code),
					Id:       convert.Int64ToType(in.ComputeServerType.Id),
					NodeType: convert.StrToType(in.ComputeServerType.NodeType),
					state:    attr.ValueStateKnown,
				}

				computeServerTypeObj, objectDiags = computeServerType.ToObjectValue(ctx)
				serverDiags.Append(objectDiags...)
			}

			typeSetObj, objectDiags := NewTypeSetValueNull().ToObjectValue(ctx)
			serverDiags.Append(objectDiags...)
			if in.TypeSet != nil {
				typeSet := TypeSetValue{
					Code:  convert.StrToType(in.TypeSet.Code),
					Id:    convert.Int64ToType(in.TypeSet.Id),
					Name:  convert.StrToType(in.TypeSet.Name),
					state: attr.ValueStateKnown,
				}

				typeSetObj, objectDiags = typeSet.ToObjectValue(ctx)
				serverDiags.Append(objectDiags...)
			}

			return ServersValue{
				ComputeServerType: computeServerTypeObj,
				Id:                convert.Int64ToType(in.Id),
				Name:              convert.StrToType(in.Name),
				TypeSet:           typeSetObj,
				state:             attr.ValueStateKnown,
			}
		},
	)
	diags.Append(serverDiags...)

	return result, diags
}

func timeToType(t *time.Time) types.String {
	if t == nil {
		return types.StringNull()
	}

	return types.StringValue(t.Format(time.RFC3339))
}

func getClusterByID(
	ctx context.Context,
	id int64,
	data *ClusterModel,
	apiClient *sdk.APIClient,
) error {
	clusterResp, hresp, err := apiClient.ClustersAPI.GetCluster(ctx, id).Execute()
	if err != nil || hresp.StatusCode != http.StatusOK {
		return fmt.Errorf("cluster %d GET failed: %s", id, errfmt.ErrMsg(err, hresp))
	}

	cluster, ok := clusterResp.GetClusterOk()
	if !ok || cluster == nil {
		return errors.New(ErrorNoClusterFound)
	}

	if err := populateClusterData(ctx, data, cluster); err != nil {
		return err
	}

	return nil
}

func getClusterByName(
	ctx context.Context,
	name string,
	data *ClusterModel,
	apiClient *sdk.APIClient,
) error {
	clustersResp, hresp, err := apiClient.ClustersAPI.ListClusters(ctx).Name(name).Execute()
	if err != nil || hresp.StatusCode != http.StatusOK {
		return fmt.Errorf("cluster %s GET failed: %s", name, errfmt.ErrMsg(err, hresp))
	}

	var matchingClusters []sdk.ListClusters200ResponseAllOfClustersInner
	for _, c := range clustersResp.GetClusters() {
		if c.GetName() == name {
			matchingClusters = append(matchingClusters, c)
		}
	}

	if len(matchingClusters) == 0 {
		return errors.New(ErrorNoClusterFound)
	}

	if len(matchingClusters) > 1 {
		return errors.New(ErrorMultipleClusters)
	}

	id, ok := matchingClusters[0].GetIdOk()
	if !ok || id == nil {
		return errors.New(ErrorNoClusterFound)
	}

	return getClusterByID(ctx, *id, data, apiClient)
}

func getCluster(
	ctx context.Context,
	data *ClusterModel,
	apiClient *sdk.APIClient,
) error {
	if !data.Id.IsNull() {
		return getClusterByID(ctx, data.Id.ValueInt64(), data, apiClient)
	}

	if !data.Name.IsNull() {
		return getClusterByName(ctx, data.Name.ValueString(), data, apiClient)
	}

	return errors.New(ErrorNoValidSearchTerms)
}

func (d *DataSource) Read(
	ctx context.Context,
	req datasource.ReadRequest,
	resp *datasource.ReadResponse,
) {
	var data ClusterModel

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

	if err := getCluster(ctx, &data, apiClient); err != nil {
		resp.Diagnostics.AddError(
			summary,
			err.Error(),
		)

		return
	}

	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}
