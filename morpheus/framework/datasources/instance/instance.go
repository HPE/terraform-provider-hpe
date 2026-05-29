// (C) Copyright 2025-2026 Hewlett Packard Enterprise Development LP

package instance

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"

	"github.com/HewlettPackard/hpe-morpheus-go-sdk/oapigen/sdk"

	"github.com/HPE/terraform-provider-hpe/morpheus/configure"
	"github.com/HPE/terraform-provider-hpe/morpheus/utils/constants"
	errfmt "github.com/HPE/terraform-provider-hpe/morpheus/utils/errfmt"
	"github.com/HPE/terraform-provider-hpe/utils/convert"
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
	resp.TypeName = req.ProviderTypeName + "_" + constants.SubProviderName + "_instance"
}

// Schema defines the schema for the data source.
func (d *DataSource) Schema(
	ctx context.Context,
	_ datasource.SchemaRequest,
	resp *datasource.SchemaResponse,
) {
	resp.Schema = InstanceDataSourceSchema(ctx)
}

func (d *DataSource) Read(
	ctx context.Context,
	req datasource.ReadRequest,
	resp *datasource.ReadResponse,
) {
	var data InstanceModel

	// Read config
	diags := req.Config.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// API GET instance
	client, err := d.NewClient(ctx)
	if err != nil {
		resp.Diagnostics.AddError("creating client failed", err.Error())

		return
	}

	if !data.Name.IsNull() {
		name := data.Name.ValueString()
		instanceListReq := client.InstancesAPI.ListInstances(ctx).Name(name)

		instanceListResp, httpResp, err := instanceListReq.Execute()
		if instanceListResp == nil || err != nil || httpResp.StatusCode != http.StatusOK {
			resp.Diagnostics.AddError(
				fmt.Sprintf("GET failed for instance '%s'", name),
				errfmt.ErrMsg(err, httpResp),
			)

			return
		}
		var instances []sdk.ListInstances200ResponseAllOfInstancesInner

		for _, instance := range instanceListResp.Instances {
			if instance.Name == nil || *instance.Name != name {
				continue
			}

			data.Id = convert.Int64ToType(instance.Id)
			instances = append(instances, instance)
		}

		if len(instances) > 1 {
			resp.Diagnostics.AddError(
				"multiple instances were returned",
				fmt.Sprintf("datasource expected 1 instance, but %d were found", len(instances)),
			)

			return
		} else if len(instances) == 0 {
			resp.Diagnostics.AddError(
				"no instance found",
				fmt.Sprintf(
					"could not find instance by name '%s'",
					name,
				),
			)

			return
		}

		// check to make sure we found the instance
		if data.Id.IsNull() {
			resp.Diagnostics.AddError(
				"no instance found",
				fmt.Sprintf(
					"could not find instance by name '%s'",
					name,
				),
			)
		}
	}

	if !data.Id.IsNull() {
		id := data.Id.ValueInt64()

		res, httpResp, err := client.InstancesAPI.GetInstance(ctx, id).Execute()
		if err != nil || httpResp.StatusCode != http.StatusOK {
			resp.Diagnostics.AddError(
				"get instance resource",
				fmt.Sprintf("instance %d GET failed: ", id)+errfmt.ErrMsg(err, httpResp),
			)

			return
		}

		config, err := getConfigFromResponse(httpResp)
		if err != nil {
			resp.Diagnostics.AddError(
				"get instance resource",
				fmt.Sprintf("instance %d parsing config failed: %s", id, err),
			)
		}

		if res.Instance == nil {
			resp.Diagnostics.AddError(
				"get instance resource",
				fmt.Sprintf("instance %d response missing instance data", id),
			)

			return
		}

		instance := *res.Instance

		if d := parseAsData(ctx, instance, &data, config); d.HasError() {
			resp.Diagnostics.Append(d...)

			return
		}
	}

	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}

// Parse the config field from the response manually as the SDK may mangle the
// data a bit
func getConfigFromResponse(r *http.Response) (map[string]any, error) {
	dec := json.NewDecoder(r.Body)
	defer r.Body.Close()
	var resp struct {
		Instance struct {
			Config map[string]any `json:"config"`
		} `json:"instance"`
	}

	if err := dec.Decode(&resp); err != nil {
		return nil, err
	}

	return resp.Instance.Config, nil
}

func parseAsData(
	ctx context.Context,
	instance sdk.GetInstance200ResponseInstance,
	data *InstanceModel,
	configData map[string]any,
) diag.Diagnostics {
	var diags diag.Diagnostics

	// auto_scale
	data.AutoScale = convert.BoolToType(instance.AutoScale)

	// cloud
	data.Cloud = NewCloudValueNull()
	if instance.Cloud != nil {
		data.Cloud = CloudValue{
			Id:        convert.Int64ToType(instance.Cloud.Id),
			Name:      convert.StrToType(instance.Cloud.Name),
			CloudType: convert.StrToType(instance.Cloud.Type),
			state:     attr.ValueStateKnown,
		}
	}

	// cluster
	cluster, d := parseCluster(ctx, instance.Cluster)
	diags = append(diags, d...)
	data.Cluster = cluster

	// config
	config, _ := convert.MapToDynamic(ctx, configData)
	data.Config = config

	// config_group
	data.ConfigGroup = convert.StrToType(instance.ConfigGroup.Get())

	// config_id
	data.ConfigId = convert.StrToType(instance.ConfigId.Get())

	// config_role
	data.ConfigRole = convert.StrToType(instance.ConfigRole.Get())

	// connection_info
	connectionInfo, d := convert.ToSetType(
		ctx,
		instance.ConnectionInfo,
		func(
			in sdk.AddInstance200ResponseAllOfOneOfInstanceConnectionInfoInner,
		) ConnectionInfoValue {
			return ConnectionInfoValue{
				Ip:    convert.StrToType(in.Ip),
				Name:  convert.StrToType(in.Name),
				Port:  convert.Int64ToType(in.Port),
				state: attr.ValueStateKnown,
			}
		},
	)
	diags.Append(d...)
	data.ConnectionInfo = connectionInfo

	// container_details
	containerDetails, d := convert.ToSetType(
		ctx,
		instance.ContainerDetails,
		func(
			in sdk.InstanceContainer2,
		) ContainerDetailsValue {
			containerDetails, d := parseContainerDetails(ctx, in)
			if d.HasError() {
				diags = append(diags, d...)
			}

			return containerDetails
		},
	)
	diags.Append(d...)
	data.ContainerDetails = containerDetails

	// containers
	data.Containers = convert.Int64SliceToSet(instance.Containers)

	// controllers
	controllers, d := convert.ToSetType(
		ctx,
		instance.Controllers,
		func(
			in sdk.ListInstances200ResponseAllOfInstancesInnerControllersInner,
		) ControllersValue {
			controllersType := ControllersTypeValue{
				Code:  convert.StrToType(in.Type.Code),
				Id:    convert.Int64ToType(in.Type.Id),
				Name:  convert.StrToType(in.Type.Name),
				state: attr.ValueStateKnown,
			}

			controllersTypeObj, d := controllersType.ToObjectValue(ctx)
			diags.Append(d...)

			return ControllersValue{
				ControllersType:    controllersTypeObj,
				Id:                 convert.Int64ToType(in.Id),
				MaxDevices:         convert.Int64ToType(in.MaxDevices),
				Name:               convert.StrToType(in.Name),
				ReservedUnitNumber: convert.Int64ToType(in.ReservedUnitNumber),
				state:              attr.ValueStateKnown,
			}
		},
	)
	diags.Append(d...)
	data.Controllers = controllers

	// cores_per_socket
	data.CoresPerSocket = convert.Int64ToType(instance.CoresPerSocket.Get())

	// created_by
	data.CreatedBy = CreatedByValue{
		Id:       convert.Int64ToType(instance.CreatedBy.Id),
		Username: convert.StrToType(instance.CreatedBy.Username),
		state:    attr.ValueStateKnown,
	}

	// current_deploy_id
	data.CurrentDeployId = convert.StrToType(instance.CurrentDeployId.Get())

	// description
	data.Description = convert.StrToType(instance.Description.Get())

	// display_name
	data.DisplayName = convert.StrToType(instance.DisplayName)

	// domain_name
	data.DomainName = convert.StrToType(instance.DomainName.Get())

	// environment
	data.Environment = convert.StrToType(instance.Environment.Get())

	// environment_prefix
	data.EnvironmentPrefix = convert.StrToType(instance.EnvironmentPrefix.Get())

	// error_message
	data.ErrorMessage = convert.StrToType(instance.ErrorMessage.Get())

	// evars
	evars, d := convert.ToSetType(
		ctx,
		instance.Evars,
		func(
			in sdk.AddInstance200ResponseAllOfOneOfInstanceEvarsInner,
		) EvarsValue {
			value := types.StringNull()

			if in.Value.String != nil {
				value = convert.StrToType(in.Value.String)
			}

			if in.Value.Float32 != nil {
				value = types.StringValue(
					strconv.FormatFloat(float64(*in.Value.Float32), 'g', -1, 32),
				)
			}

			return EvarsValue{
				Name:   convert.StrToType(in.Name),
				Value:  value,
				Masked: convert.BoolToType(in.Masked),
				Export: convert.BoolToType(in.Export),
				state:  attr.ValueStateKnown,
			}
		},
	)
	diags.Append(d...)
	data.Evars = evars

	// expire_count
	data.ExpireCount = convert.Int64ToType(instance.ExpireCount)

	// expire_days
	data.ExpireDays = convert.Int64ToType(instance.ExpireDays)

	// expire_warning_sent
	data.ExpireWarningSent = convert.BoolToType(instance.ExpireWarningSent)

	// firewall_enabled
	data.FirewallEnabled = convert.BoolToType(instance.FirewallEnabled)

	// group
	group := instance.Group.Get()
	if group.Id == nil || group.Name == nil {
		diags.AddError(
			"missing attribute",
			"missing group id or name",
		)

		return diags
	}

	groupID := group.Id
	groupName := group.Name

	data.Group = GroupValue{
		Id:    convert.Int64ToType(groupID),
		Name:  convert.StrToType(groupName),
		state: attr.ValueStateKnown,
	}

	// host_name
	data.HostName = convert.StrToType(instance.HostName)

	// hourly_cost
	data.HourlyCost = convert.NumToType(instance.HourlyCost)

	// hourly_price
	data.HourlyPrice = convert.NumToType(instance.HourlyPrice)

	// id
	data.Id = convert.Int64ToType(instance.Id)

	// instance_context
	data.InstanceContext = convert.StrToType(instance.InstanceContext.Get())

	// instance_price
	data.InstancePrice = NewInstancePriceValueNull()
	if instance.InstancePrice != nil {
		data.InstancePrice = InstancePriceValue{
			Cost:     convert.NumToType(instance.InstancePrice.Cost),
			Currency: convert.StrToType(instance.InstancePrice.Currency),
			Price:    convert.NumToType(instance.InstancePrice.Price),
			Unit:     convert.StrToType(instance.InstancePrice.Unit),
			state:    attr.ValueStateKnown,
		}
	}

	// instance_type
	data.InstanceType = NewInstanceTypeValueNull()
	if instance.InstanceType != nil {
		data.InstanceType = InstanceTypeValue{
			Category: convert.StrToType(instance.InstanceType.Category),
			Code:     convert.StrToType(instance.InstanceType.Code),
			Id:       convert.Int64ToType(instance.InstanceType.Id),
			Image:    convert.StrToType(instance.InstanceType.Image),
			Name:     convert.StrToType(instance.InstanceType.Name),
			state:    attr.ValueStateKnown,
		}
	}

	// instance_version
	data.InstanceVersion = convert.StrToType(instance.InstanceVersion)

	// interfaces
	interfaces, d := convert.ToSetType(
		ctx,
		instance.Interfaces,
		func(
			in sdk.AddInstance200ResponseAllOfOneOfInstanceInterfacesInner,
		) InterfacesValue {
			pool := NewInterfaceNetworkPoolValueNull()
			if in.Network.Pool != nil {
				pool = InterfaceNetworkPoolValue{
					Id:    convert.Int64ToType(in.Network.Pool.Id),
					Name:  convert.StrToType(in.Network.Pool.Name),
					state: attr.ValueStateKnown,
				}
			}

			poolObj, d := pool.ToObjectValue(ctx)
			diags = append(diags, d...)

			network := NewNetworkValueNull()
			if in.Network != nil {
				network = NetworkValue{
					DhcpServer:           convert.BoolToType(in.Network.DhcpServer),
					Group:                convert.Int64ToType(in.Network.Group),
					Id:                   convert.Int64ToType(in.Network.Id),
					Name:                 convert.StrToType(in.Network.Name),
					InterfaceNetworkPool: poolObj,
					Subnet:               convert.StrToType(in.Network.Subnet),
					state:                attr.ValueStateKnown,
				}
			}

			networkObj, d := network.ToObjectValue(ctx)
			diags.Append(d...)

			networkInterfaces, d := convert.ToSetType(
				ctx,
				in.NetworkInterfaces,
				func(in sdk.InstanceInterfacesNetworkInterfacesInner1) NetworkInterfacesValue {
					network := NewInterfaceNetworkValueNull()
					if in.Network != nil {
						network = InterfaceNetworkValue{
							DhcpServer: convert.BoolToType(in.Network.DhcpServer),
							Group:      convert.Int64ToType(in.Network.Group),
							Id:         convert.Int64ToType(in.Network.Id),
							Name:       convert.StrToType(in.Network.Name),
							Pool:       poolObj,
							Subnet:     convert.StrToType(in.Network.Subnet),
							state:      attr.ValueStateKnown,
						}
					}
					networkObj, d := network.ToObjectValue(ctx)
					diags = append(diags, d...)

					return NetworkInterfacesValue{
						Id:                     convert.Int64ToType(in.Id.Int64),
						InterfaceNetwork:       networkObj,
						IpAddress:              convert.StrToType(in.IpAddress),
						IpMode:                 convert.StrToType(in.IpMode),
						NetworkInterfaceTypeId: convert.Int64ToType(in.NetworkInterfaceTypeId),
						state:                  attr.ValueStateKnown,
					}
				},
			)
			diags.Append(d...)

			return InterfacesValue{
				Id:                     convert.Int64ToType(in.Id.Int64),
				IpAddress:              convert.StrToType(in.IpAddress),
				IpMode:                 convert.StrToType(in.IpMode),
				Network:                networkObj,
				NetworkInterfaceTypeId: convert.Int64ToType(in.NetworkInterfaceTypeId),
				NetworkInterfaces:      networkInterfaces,
				state:                  attr.ValueStateKnown,
			}
		},
	)
	diags.Append(d...)
	data.Interfaces = interfaces

	// is_busy
	data.IsBusy = convert.BoolToType(instance.IsBusy)

	// is_scalable
	data.IsScalable = convert.BoolToType(instance.IsScalable)

	// labels
	data.Labels = convert.StrSliceToSet(instance.Labels)

	// layout
	data.Layout = NewLayoutValueNull()
	if instance.Layout != nil {
		data.Layout = LayoutValue{
			Id:                convert.Int64ToType(instance.Layout.Id),
			Name:              convert.StrToType(instance.Layout.Name),
			ProvisionTypeCode: convert.StrToType(instance.Layout.ProvisionTypeCode),
			ProvisionTypeId:   convert.Int64ToType(instance.Layout.ProvisionTypeId),
			state:             attr.ValueStateKnown,
		}
	}

	// locked
	data.Locked = convert.BoolToType(instance.Locked)

	// max_cores
	data.MaxCores = convert.Int64ToType(instance.MaxCores)

	// max_cpu
	data.MaxCpu = convert.Int64ToType(instance.MaxCpu.Get())

	// max_memory
	data.MaxMemory = convert.Int64ToType(instance.MaxMemory)

	// max_storage
	data.MaxStorage = convert.Int64ToType(instance.MaxStorage)

	// name
	data.Name = convert.StrToType(instance.Name)

	// network_level
	data.NetworkLevel = convert.StrToType(instance.NetworkLevel)

	// notes
	data.Notes = convert.StrToType(instance.Notes.Get())

	// owner
	data.Owner = NewOwnerValueNull()
	if instance.Owner != nil {
		data.Owner = OwnerValue{
			Id:       convert.Int64ToType(instance.Owner.Id),
			Username: convert.StrToType(instance.Owner.Username),
			state:    attr.ValueStateKnown,
		}
	}

	// plan
	data.Plan = NewPlanValueNull()
	if instance.Plan != nil {
		data.Plan = PlanValue{
			Code:  convert.StrToType(instance.Plan.Code),
			Id:    convert.Int64ToType(instance.Plan.Id),
			Name:  convert.StrToType(instance.Plan.Name),
			state: attr.ValueStateKnown,
		}
	}

	// power_schedule
	data.PowerSchedule = convert.StrToType(instance.PowerSchedule.Get())

	// renew_days
	data.RenewDays = convert.Int64ToType(instance.RenewDays)

	// servers
	data.Servers = convert.Int64SliceToSet(instance.Servers)

	// shutdown_count
	data.ShutdownCount = convert.Int64ToType(instance.ShutdownCount)

	// shutdown_days
	data.ShutdownDays = convert.Int64ToType(instance.ShutdownDays)

	// shutdown_renew_days
	data.ShutdownRenewDays = convert.Int64ToType(instance.ShutdownRenewDays)

	// shutdown_warning_sent
	data.ShutdownWarningSent = convert.BoolToType(instance.ShutdownWarningSent)

	// status
	data.Status = convert.StrToType(instance.Status)

	// status_eta
	data.StatusEta = convert.StrToType(instance.StatusEta.Get())

	// status_message
	data.StatusMessage = convert.StrToType(instance.StatusMessage.Get())

	// status_percent
	data.StatusPercent = convert.StrToType(instance.StatusPercent.Get())

	// tags
	tags, d := convert.ToSetType(
		ctx,
		instance.Tags,
		func(
			in sdk.AddInstance200ResponseAllOfOneOfInstanceTagsInner,
		) TagsValue {
			return TagsValue{
				Name:  convert.StrToType(in.Name),
				Value: convert.StrToType(in.Value),
				state: attr.ValueStateKnown,
			}
		},
	)
	diags.Append(d...)
	data.Tags = tags

	// tenant
	tenant := instance.Tenant.Get()
	if tenant.Id == nil || tenant.Name == nil {
		diags.AddError(
			"missing attribute",
			"missin tenant id or name",
		)

		return diags
	}

	tenantID := tenant.Id
	tenantName := tenant.Name

	data.Tenant = TenantValue{
		Id:    convert.Int64ToType(tenantID),
		Name:  convert.StrToType(tenantName),
		state: attr.ValueStateKnown,
	}

	// tenant_id
	data.TenantId = convert.Int64ToType(tenantID)

	// user_status
	data.UserStatus = convert.StrToType(instance.UserStatus.Get())

	// uuid
	data.Uuid = convert.StrToType(instance.Uuid)

	// volumes
	volumes, d := convert.ToSetType(
		ctx,
		instance.Volumes,
		func(
			in sdk.AddInstance200ResponseAllOfOneOfInstanceVolumesInner,
		) VolumesValue {
			return VolumesValue{
				ControllerId:         convert.Int64ToType(in.ControllerId),
				ControllerMountPoint: convert.StrToType(in.ControllerMountPoint),
				DatastoreId:          convert.StrToType(in.DatastoreId),
				DisplayOrder:         convert.Int64ToType(in.DisplayOrder),
				Id:                   convert.Int64ToType(in.Id),
				MaxIops:              convert.StrToType(in.MaxIOPS),
				MaxStorage:           convert.Int64ToType(in.MaxStorage),
				Name:                 convert.StrToType(in.Name),
				PlanResizable:        convert.BoolToType(in.PlanResizable),
				Resizeable:           convert.BoolToType(in.Resizeable),
				RootVolume:           convert.BoolToType(in.RootVolume),
				ShortName:            convert.StrToType(in.ShortName),
				Size:                 convert.Int64ToType(in.Size),
				StorageType:          convert.Int64ToType(in.StorageType),
				UnitNumber:           convert.StrToType(in.UnitNumber),
				Uuid:                 convert.StrToType(in.Uuid),
				state:                attr.ValueStateKnown,
			}
		},
	)
	diags.Append(d...)
	data.Volumes = volumes

	return diags
}

func parseCluster(
	ctx context.Context,
	cluster *sdk.GetInstance200ResponseInstanceCluster,
) (ClusterValue, diag.Diagnostics) {
	var diags diag.Diagnostics

	if cluster == nil {
		return NewClusterValueNull(), diags
	}

	clusterType := NewClusterTypeValueNull()
	if cluster.Type != nil {
		clusterType = ClusterTypeValue{
			Code:  convert.StrToType(cluster.Type.Code),
			Id:    convert.Int64ToType(cluster.Type.Id),
			Name:  convert.StrToType(cluster.Type.Name),
			state: attr.ValueStateKnown,
		}
	}
	clusterTypeObj, d := clusterType.ToObjectValue(ctx)
	diags.Append(d...)

	return ClusterValue{
		ClusterType: clusterTypeObj,
		Id:          convert.Int64ToType(cluster.Id),
		Name:        convert.StrToType(cluster.Name),
		state:       attr.ValueStateKnown,
	}, diags
}

func parseContainerDetails(
	ctx context.Context,
	containerDetails sdk.InstanceContainer2,
) (ContainerDetailsValue, diag.Diagnostics) {
	var diags diag.Diagnostics

	containerType, d := parseContainerType(ctx, containerDetails)
	diags = append(diags, d...)

	containerTypeSet := NewContainerTypeSetValueNull()
	if containerDetails.ContainerTypeSet != nil {
		containerTypeSet = ContainerTypeSetValue{
			Code:  convert.StrToType(containerDetails.ContainerTypeSet.Code),
			Id:    convert.Int64ToType(containerDetails.ContainerTypeSet.Id),
			state: attr.ValueStateKnown,
		}
	}

	containerTypeSetObj, d := containerTypeSet.ToObjectValue(ctx)
	diags.Append(d...)

	instance := NewInstanceValueNull()
	if containerDetails.Instance != nil {
		instance = InstanceValue{
			Id:    convert.Int64ToType(containerDetails.Instance.Id),
			Name:  convert.StrToType(containerDetails.Instance.Name),
			state: attr.ValueStateKnown,
		}
	}
	instanceObj, d := instance.ToObjectValue(ctx)
	diags.Append(d...)

	ports, d := convert.ToSetType(
		ctx,
		containerDetails.Ports,
		func(
			in sdk.InstanceContainer1PortsInner,
		) PortsValue {
			return PortsValue{
				DisplayName:         convert.StrToType(in.DisplayName),
				Export:              convert.BoolToType(in.Export),
				ExportName:          convert.StrToType(in.ExportName),
				External:            convert.Int64ToType(in.External),
				ExternalIp:          convert.StrToType(in.ExternalIp),
				Id:                  convert.Int64ToType(in.Id),
				Internal:            convert.Int64ToType(in.Internal),
				InternalIp:          convert.StrToType(in.InternalIp),
				Link:                convert.BoolToType(in.Link),
				LoadBalance:         convert.BoolToType(in.LoadBalance),
				LoadBalanceProtocol: convert.StrToType(in.LoadBalanceProtocol),
				PrimaryPort:         convert.BoolToType(in.PrimaryPort),
				Protocol:            convert.StrToType(in.Protocol),
				Visible:             convert.BoolToType(in.Visible),
				state:               attr.ValueStateKnown,
			}
		},
	)
	diags = append(diags, d...)

	computerServerType := NewComputeServerTypeValueNull()
	if containerDetails.Server.ComputeServerType != nil {
		computerServerType = ComputeServerTypeValue{
			Code:           convert.StrToType(containerDetails.Server.ComputeServerType.Code.Get()),
			ExternalDelete: convert.BoolToType(containerDetails.Server.ComputeServerType.ExternalDelete.Get()),
			Format:         convert.StrToType(containerDetails.Server.ComputeServerType.Format.Get()),
			Id:             convert.Int64ToType(containerDetails.Server.ComputeServerType.Id.Get()),
			Managed:        convert.BoolToType(containerDetails.Server.ComputeServerType.Managed.Get()),
			Name:           convert.StrToType(containerDetails.Server.ComputeServerType.Name.Get()),
			state:          attr.ValueStateKnown,
		}
	}

	computeServerTypeObj, d := computerServerType.ToObjectValue(ctx)
	diags.Append(d...)

	containerOwner := NewContainerOwnerValueNull()
	if containerDetails.Server.Owner != nil {
		containerOwner = ContainerOwnerValue{
			Id:       convert.Int64ToType(containerDetails.Server.Owner.Id),
			Username: convert.StrToType(containerDetails.Server.Owner.Username),
			state:    attr.ValueStateKnown,
		}
	}
	containerOwnerObj, d := containerOwner.ToObjectValue(ctx)
	diags.Append(d...)

	containerInterfaces, d := parseContainerInterfaces(ctx, containerDetails)
	if d.HasError() {
		diags = append(diags, d...)
	}

	childVolumes, d := parseChildVolumes(ctx, containerDetails)
	if d.HasError() {
		diags = append(diags, d...)
	}

	account := NewAccountValueNull()
	if containerDetails.Server.Account != nil {
		account = AccountValue{
			Id:    convert.Int64ToType(containerDetails.Server.Account.Id),
			Name:  convert.StrToType(containerDetails.Server.Account.Name),
			state: attr.ValueStateKnown,
		}
	}
	accountObj, d := account.ToObjectValue(ctx)
	diags.Append(d...)
	containerPlan := NewContainerPlanValueNull()
	if containerDetails.Server.Plan != nil {
		containerPlan = ContainerPlanValue{
			Code:  convert.StrToType(containerDetails.Server.Plan.Code),
			Id:    convert.Int64ToType(containerDetails.Server.Plan.Id),
			Name:  convert.StrToType(containerDetails.Server.Plan.Name),
			state: attr.ValueStateKnown,
		}
	}

	containerPlanObj, d := containerPlan.ToObjectValue(ctx)
	diags.Append(d...)

	serverOs := NewServerOsValueNull()
	if containerDetails.Server.ServerOs != nil {
		serverOs = ServerOsValue{
			BitCount:    convert.Int64ToType(containerDetails.Server.ServerOs.BitCount),
			Category:    convert.StrToType(containerDetails.Server.ServerOs.Category),
			Code:        convert.StrToType(containerDetails.Server.ServerOs.Code),
			Description: convert.StrToType(containerDetails.Server.ServerOs.Description.Get()),
			Id:          convert.Int64ToType(containerDetails.Server.ServerOs.Id),
			Name:        convert.StrToType(containerDetails.Server.ServerOs.Name),
			OsFamily:    convert.StrToType(containerDetails.Server.ServerOs.OsFamily),
			OsVersion:   convert.StrToType(containerDetails.Server.ServerOs.OsVersion),
			Platform:    convert.StrToType(containerDetails.Server.ServerOs.Platform),
			Vendor:      convert.StrToType(containerDetails.Server.ServerOs.Vendor),
			state:       attr.ValueStateKnown,
		}
	}

	serverOsObj, d := serverOs.ToObjectValue(ctx)
	diags.Append(d...)

	sourceImage := NewSourceImageValueNull()
	if containerDetails.Server.SourceImage != nil {
		sourceImage = SourceImageValue{
			Code:  convert.StrToType(containerDetails.Server.SourceImage.Code),
			Id:    convert.Int64ToType(containerDetails.Server.SourceImage.Id),
			Name:  convert.StrToType(containerDetails.Server.SourceImage.Name),
			state: attr.ValueStateKnown,
		}
	}
	sourceImageObj, d := sourceImage.ToObjectValue(ctx)
	diags.Append(d...)

	zone := NewZoneValueNull()
	if containerDetails.Server.Zone != nil {
		zone = ZoneValue{
			Id:    convert.Int64ToType(containerDetails.Server.Zone.Id),
			Name:  convert.StrToType(containerDetails.Server.Zone.Name),
			state: attr.ValueStateKnown,
		}
	}
	zoneObj, d := zone.ToObjectValue(ctx)
	diags.Append(d...)
	server := NewServerValueNull()
	if containerDetails.Server != nil {
		server = ServerValue{
			Account:             accountObj,
			AccountId:           convert.Int64ToType(containerDetails.Server.AccountId),
			AgentInstalled:      convert.BoolToType(containerDetails.Server.AgentInstalled),
			ChildVolumes:        childVolumes,
			ComputeServerType:   computeServerTypeObj,
			ContainerInterfaces: containerInterfaces,
			ContainerOwner:      containerOwnerObj,
			ContainerPlan:       containerPlanObj,
			Description:         convert.StrToType(containerDetails.Server.Description.Get()),
			ErrorMessage:        convert.StrToType(containerDetails.Server.ErrorMessage.Get()),
			ExternalId:          convert.StrToType(containerDetails.Server.ExternalId.Get()),
			ExternalIp:          convert.StrToType(containerDetails.Server.ExternalIp.Get()),
			HostName:            convert.StrToType(containerDetails.Server.HostName),
			Id:                  convert.Int64ToType(containerDetails.Server.Id),
			InternalId:          convert.StrToType(containerDetails.Server.InternalId.Get()),
			InternalIp:          convert.StrToType(containerDetails.Server.InternalIp.Get()),
			MaxCores:            convert.Int64ToType(containerDetails.Server.MaxCores),
			MaxMemory:           convert.Int64ToType(containerDetails.Server.MaxMemory),
			MaxStorage:          convert.Int64ToType(containerDetails.Server.MaxStorage),
			Platform:            convert.StrToType(containerDetails.Server.Platform),
			ResourcePoolId:      convert.Int64ToType(containerDetails.Server.ResourcePoolId),
			ServerOs:            serverOsObj,
			SiteId:              convert.Int64ToType(containerDetails.Server.SiteId),
			SourceImage:         sourceImageObj,
			SshHost:             convert.StrToType(containerDetails.Server.SshHost.Get()),
			SshPort:             convert.Int64ToType(containerDetails.Server.SshPort),
			Status:              convert.StrToType(containerDetails.Server.Status),
			StatusEta:           convert.Int64ToType(containerDetails.Server.StatusEta.Get()),
			StatusMessage:       convert.StrToType(containerDetails.Server.StatusMessage.Get()),
			StatusPercent:       convert.Int64ToType(containerDetails.Server.StatusPercent.Get()),
			Uuid:                convert.StrToType(containerDetails.Server.Uuid),
			Visibility:          convert.StrToType(containerDetails.Server.Visibility),
			Zone:                zoneObj,
			ZoneId:              convert.Int64ToType(containerDetails.Server.ZoneId),
			state:               attr.ValueStateKnown,
		}
	}

	serverObj, d := server.ToObjectValue(ctx)
	diags.Append(d...)

	return ContainerDetailsValue{
		AccountId:        convert.Int64ToType(containerDetails.AccountId),
		ContainerCreated: convert.BoolToType(containerDetails.ContainerCreated),
		ContainerType:    containerType,
		ContainerTypeSet: containerTypeSetObj,
		ExternalDomain:   convert.StrToType(containerDetails.ExternalDomain),
		ExternalFqdn:     convert.StrToType(containerDetails.ExternalFqdn),
		ExternalHostname: convert.StrToType(containerDetails.ExternalHostname),
		Hostname:         convert.StrToType(containerDetails.Hostname),
		Id:               convert.Int64ToType(containerDetails.Id),
		Instance:         instanceObj,
		InternalHostname: convert.StrToType(containerDetails.InternalHostname),
		InternalIp:       convert.StrToType(containerDetails.InternalIp),
		Ip:               convert.StrToType(containerDetails.Ip),
		Name:             convert.StrToType(containerDetails.Name),
		Ports:            ports,
		Server:           serverObj,
		Uuid:             convert.StrToType(containerDetails.Uuid),
		VolumeCreated:    convert.BoolToType(containerDetails.VolumeCreated),
		state:            attr.ValueStateKnown,
	}, diags
}

func parseChildVolumes(
	ctx context.Context,
	containerDetails sdk.InstanceContainer2,
) (basetypes.SetValue, diag.Diagnostics) {
	var diags diag.Diagnostics

	volumes, d := convert.ToSetType(
		ctx,
		containerDetails.Server.Volumes,
		func(in sdk.InstanceContainerServerVolume1) ChildVolumesValue {
			datastore := NewDatastoreValueNull()
			if in.Datastore != nil {
				datastore = DatastoreValue{
					Id:    convert.Int64ToType(in.Datastore.Id),
					Name:  convert.StrToType(in.Datastore.Name),
					state: attr.ValueStateKnown,
				}
			}

			datastoreObj, d := datastore.ToObjectValue(ctx)
			diags = append(diags, d...)

			volumeZone := NewVolumeZoneValueNull()
			if in.Zone != nil {
				volumeZone = VolumeZoneValue{
					Id:    convert.Int64ToType(in.Zone.Id),
					Name:  convert.StrToType(in.Zone.Name),
					state: attr.ValueStateKnown,
				}
			}

			volumeZoneObj, d := volumeZone.ToObjectValue(ctx)
			diags = append(diags, d...)

			storageServer := NewStorageServerValueNull()
			if in.StorageServer != nil {
				storageServer = StorageServerValue{
					Id:    convert.Int64ToType(in.StorageServer.Id),
					Name:  convert.StrToType(in.StorageServer.Name),
					state: attr.ValueStateKnown,
				}
			}

			storageServerObj, d := storageServer.ToObjectValue(ctx)
			diags = append(diags, d...)

			return ChildVolumesValue{
				Category:             convert.StrToType(in.Category),
				ConfigurableIops:     convert.BoolToType(in.ConfigurableIOPS),
				ControllerId:         convert.Int64ToType(in.ControllerId),
				ControllerMountPoint: convert.StrToType(in.ControllerMountPoint),
				Datastore:            datastoreObj,
				DatastoreId:          convert.Int64ToType(in.DatastoreId),
				DeviceDisplayName:    convert.StrToType(in.DeviceDisplayName),
				DeviceName:           convert.StrToType(in.DeviceName),
				DiskMode:             convert.StrToType(in.DiskMode),
				DiskType:             convert.StrToType(in.DiskType),
				ExternalId:           convert.StrToType(in.ExternalId),
				Id:                   convert.Int64ToType(in.Id),
				InternalId:           convert.StrToType(in.InternalId),
				MaxIops:              convert.StrToType(in.MaxIOPS),
				MaxStorage:           convert.Int64ToType(in.MaxStorage),
				Name:                 convert.StrToType(in.Name),
				Resizeable:           convert.BoolToType(in.Resizeable),
				RootVolume:           convert.BoolToType(in.RootVolume),
				StorageServer:        storageServerObj,
				TypeId:               convert.Int64ToType(in.TypeId),
				UniqueId:             convert.StrToType(in.UniqueId),
				UnitNumber:           convert.StrToType(in.UnitNumber),
				Uuid:                 convert.StrToType(in.Uuid),
				VolumeZone:           volumeZoneObj,
				ZoneId:               convert.Int64ToType(in.ZoneId),
				state:                attr.ValueStateKnown,
			}
		},
	)

	diags = append(diags, d...)

	return volumes, diags
}

func parseContainerType(
	ctx context.Context,
	containerDetails sdk.InstanceContainer2,
) (basetypes.ObjectValue, diag.Diagnostics) {
	var diags diag.Diagnostics

	containerType := ContainerTypeValue{
		Code:  convert.StrToType(containerDetails.ContainerType.Code),
		Id:    convert.Int64ToType(containerDetails.ContainerType.Id),
		Name:  convert.StrToType(containerDetails.ContainerType.Name),
		state: attr.ValueStateKnown,
	}

	containerTypeObj, d := containerType.ToObjectValue(ctx)
	if d.HasError() {
		diags = append(diags, d...)
	}

	return containerTypeObj, diags
}

func parseContainerInterfaces(
	ctx context.Context,
	containerDetails sdk.InstanceContainer2,
) (basetypes.SetValue, diag.Diagnostics) {
	var diags diag.Diagnostics

	interfaces, d := convert.ToSetType(
		ctx,
		containerDetails.Server.Interfaces,
		func(in sdk.InstanceContainerServerInterfacesInner1) ContainerInterfacesValue {
			childInterfaces, d := convert.ToSetType(
				ctx,
				in.Interfaces,
				func(in sdk.InstanceContainerServerInstancesInnerInner1) ChildInterfacesValue {
					return ChildInterfacesValue{
						Id:    convert.Int64ToType(in.Id),
						Name:  convert.StrToType(in.Name),
						state: attr.ValueStateKnown,
					}
				},
			)
			diags = append(diags, d...)

			childNetwork := NewChildNetworkValueNull()
			if in.Network != nil {
				childNetwork = ChildNetworkValue{
					Id:    convert.Int64ToType(in.Network.Id),
					Name:  convert.StrToType(in.Network.Name),
					state: attr.ValueStateKnown,
				}
			}
			childNetworkObj, d := childNetwork.ToObjectValue(ctx)
			diags = append(diags, d...)

			networkGroup := NewNetworkGroupValueNull()
			if in.NetworkGroup != nil {
				networkGroup = NetworkGroupValue{
					Id:    convert.Int64ToType(in.NetworkGroup.Id),
					Name:  convert.StrToType(in.NetworkGroup.Name),
					state: attr.ValueStateKnown,
				}
			}

			networkGroupObj, d := networkGroup.ToObjectValue(ctx)
			diags = append(diags, d...)

			networkPool := NewNetworkPoolValueNull()
			if in.NetworkPool != nil {
				networkPool = NetworkPoolValue{
					Id:    convert.Int64ToType(in.NetworkPool.Id),
					Name:  convert.StrToType(in.NetworkPool.Name),
					state: attr.ValueStateKnown,
				}
			}

			networkPoolObj, d := networkPool.ToObjectValue(ctx)
			diags = append(diags, d...)

			return ContainerInterfacesValue{
				Active:           convert.BoolToType(in.Active),
				ChildInterfaces:  childInterfaces,
				ChildNetwork:     childNetworkObj,
				Dhcp:             convert.BoolToType(in.Dhcp),
				Id:               convert.Int64ToType(in.Id),
				IpAddress:        convert.StrToType(in.IpAddress),
				IpMode:           convert.StrToType(in.IpMode),
				MacAddress:       convert.StrToType(in.MacAddress),
				Name:             convert.StrToType(in.Name),
				NetworkGroup:     networkGroupObj,
				NetworkPool:      networkPoolObj,
				PoolAssigned:     convert.BoolToType(in.PoolAssigned),
				PrimaryInterface: convert.BoolToType(in.PrimaryInterface),
				PublicIpAddress:  convert.StrToType(in.PublicIpAddress),
				UniqueId:         convert.StrToType(in.UniqueId),
				state:            attr.ValueStateKnown,
			}
		},
	)

	diags = append(diags, d...)

	return interfaces, diags
}
