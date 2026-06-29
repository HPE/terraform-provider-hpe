# Instance

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | Pointer to **int64** |  | [optional] 
**Uuid** | Pointer to **string** |  | [optional] 
**AccountId** | Pointer to **int64** |  | [optional] 
**Tenant** | Pointer to [**NullableInstanceTenant**](InstanceTenant.md) |  | [optional] 
**InstanceType** | Pointer to [**InstanceInstanceType**](InstanceInstanceType.md) |  | [optional] 
**Group** | Pointer to [**NullableInstanceGroup**](InstanceGroup.md) |  | [optional] 
**Cloud** | Pointer to [**InstanceCloud**](InstanceCloud.md) |  | [optional] 
**Cluster** | Pointer to [**InstanceCluster**](InstanceCluster.md) |  | [optional] 
**Containers** | Pointer to **[]int64** |  | [optional] 
**Servers** | Pointer to **[]int64** |  | [optional] 
**ConnectionInfo** | Pointer to [**[]AddInstance200ResponseAllOfOneOfInstanceConnectionInfoInner**](AddInstance200ResponseAllOfOneOfInstanceConnectionInfoInner.md) |  | [optional] 
**Layout** | Pointer to [**InstanceLayout**](InstanceLayout.md) |  | [optional] 
**Plan** | Pointer to [**InstancePlan**](InstancePlan.md) |  | [optional] 
**Name** | Pointer to **string** |  | [optional] 
**DisplayName** | Pointer to **string** |  | [optional] 
**Description** | Pointer to **NullableString** |  | [optional] 
**Environment** | Pointer to **NullableString** |  | [optional] 
**Config** | Pointer to [**InstanceConfig**](InstanceConfig.md) |  | [optional] 
**ConfigGroup** | Pointer to **NullableString** |  | [optional] 
**ConfigId** | Pointer to **NullableString** |  | [optional] 
**ConfigRole** | Pointer to **NullableString** |  | [optional] 
**Volumes** | Pointer to [**[]AddInstance200ResponseAllOfOneOfInstanceVolumesInner**](AddInstance200ResponseAllOfOneOfInstanceVolumesInner.md) |  | [optional] 
**Controllers** | Pointer to [**[]ListInstances200ResponseAllOfInstancesInnerControllersInner**](ListInstances200ResponseAllOfInstancesInnerControllersInner.md) |  | [optional] 
**Interfaces** | Pointer to [**[]AddInstance200ResponseAllOfOneOfInstanceInterfacesInner**](AddInstance200ResponseAllOfOneOfInstanceInterfacesInner.md) |  | [optional] 
**CustomOptions** | Pointer to **map[string]interface{}** |  | [optional] 
**InstanceVersion** | Pointer to **string** |  | [optional] 
**Labels** | Pointer to **[]string** |  | [optional] 
**Tags** | Pointer to [**[]AddInstance200ResponseAllOfOneOfInstanceTagsInner**](AddInstance200ResponseAllOfOneOfInstanceTagsInner.md) |  | [optional] 
**Evars** | Pointer to [**[]AddInstance200ResponseAllOfOneOfInstanceEvarsInner**](AddInstance200ResponseAllOfOneOfInstanceEvarsInner.md) |  | [optional] 
**MaxMemory** | Pointer to **int64** |  | [optional] 
**MaxStorage** | Pointer to **int64** |  | [optional] 
**MaxCores** | Pointer to **int64** |  | [optional] 
**CoresPerSocket** | Pointer to **NullableInt64** |  | [optional] 
**MaxCpu** | Pointer to **NullableInt64** |  | [optional] 
**HourlyCost** | Pointer to **float32** |  | [optional] 
**HourlyPrice** | Pointer to **float32** |  | [optional] 
**InstancePrice** | Pointer to [**InstanceInstancePrice**](InstanceInstancePrice.md) |  | [optional] 
**DateCreated** | Pointer to **time.Time** |  | [optional] 
**LastUpdated** | Pointer to **time.Time** |  | [optional] 
**HostName** | Pointer to **string** |  | [optional] 
**DomainName** | Pointer to **NullableString** |  | [optional] 
**EnvironmentPrefix** | Pointer to **NullableString** |  | [optional] 
**FirewallEnabled** | Pointer to **bool** |  | [optional] 
**NetworkLevel** | Pointer to **string** |  | [optional] 
**AutoScale** | Pointer to **bool** |  | [optional] 
**InstanceContext** | Pointer to **NullableString** |  | [optional] 
**CurrentDeployId** | Pointer to **NullableString** |  | [optional] 
**Locked** | Pointer to **bool** |  | [optional] 
**Status** | Pointer to **string** |  | [optional] 
**StatusMessage** | Pointer to **NullableString** |  | [optional] 
**ErrorMessage** | Pointer to **NullableString** |  | [optional] 
**StatusDate** | Pointer to **time.Time** |  | [optional] 
**StatusPercent** | Pointer to **NullableString** |  | [optional] 
**StatusEta** | Pointer to **NullableString** |  | [optional] 
**UserStatus** | Pointer to **NullableString** |  | [optional] 
**ExpireDays** | Pointer to **int64** |  | [optional] 
**RenewDays** | Pointer to **int64** |  | [optional] 
**ExpireCount** | Pointer to **int64** |  | [optional] 
**ExpireDate** | Pointer to **time.Time** |  | [optional] 
**ExpireWarningDate** | Pointer to **time.Time** |  | [optional] 
**ExpireWarningSent** | Pointer to **bool** |  | [optional] 
**ShutdownDays** | Pointer to **int64** |  | [optional] 
**ShutdownRenewDays** | Pointer to **int64** |  | [optional] 
**ShutdownCount** | Pointer to **int64** |  | [optional] 
**ShutdownDate** | Pointer to **time.Time** |  | [optional] 
**ShutdownWarningDate** | Pointer to **time.Time** |  | [optional] 
**ShutdownWarningSent** | Pointer to **bool** |  | [optional] 
**RemovalDate** | Pointer to **NullableTime** |  | [optional] 
**CreatedBy** | Pointer to [**InstanceCreatedBy**](InstanceCreatedBy.md) |  | [optional] 
**Owner** | Pointer to [**InstanceOwner**](InstanceOwner.md) |  | [optional] 
**Notes** | Pointer to **NullableString** |  | [optional] 
**Stats** | Pointer to [**InstanceStats**](InstanceStats.md) |  | [optional] 
**PowerSchedule** | Pointer to **NullableString** |  | [optional] 
**IsScalable** | Pointer to **bool** |  | [optional] 
**InstanceThreshold** | Pointer to **map[string]interface{}** |  | [optional] 
**IsBusy** | Pointer to **bool** |  | [optional] 
**Apps** | Pointer to **[]map[string]interface{}** |  | [optional] 
**CurrentLoadBalancerInstances** | Pointer to **[]map[string]interface{}** |  | [optional] 
**CurrentLoadBalancerContainersIn** | Pointer to **int32** |  | [optional] 
**CurrentLoadBalancerContainersOut** | Pointer to **int32** |  | [optional] 
**LastDeploy** | Pointer to **map[string]interface{}** |  | [optional] 
**ContainerDetails** | Pointer to [**[]InstanceContainer3**](InstanceContainer3.md) |  | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &Instance{
    // Set fields directly
}
```

### Tenant (Nullable)

Use the Nullable wrapper methods:
- `obj.Tenant.IsSet()` — check if set
- `obj.Tenant.Get()` — get the inner value (returns pointer)
- `obj.Tenant.Set(&val)` — set the value
- `obj.Tenant.Unset()` — clear the value
### Group (Nullable)

Use the Nullable wrapper methods:
- `obj.Group.IsSet()` — check if set
- `obj.Group.Get()` — get the inner value (returns pointer)
- `obj.Group.Set(&val)` — set the value
- `obj.Group.Unset()` — clear the value
### Description (Nullable)

Use the Nullable wrapper methods:
- `obj.Description.IsSet()` — check if set
- `obj.Description.Get()` — get the inner value (returns pointer)
- `obj.Description.Set(&val)` — set the value
- `obj.Description.Unset()` — clear the value
### Environment (Nullable)

Use the Nullable wrapper methods:
- `obj.Environment.IsSet()` — check if set
- `obj.Environment.Get()` — get the inner value (returns pointer)
- `obj.Environment.Set(&val)` — set the value
- `obj.Environment.Unset()` — clear the value
### ConfigGroup (Nullable)

Use the Nullable wrapper methods:
- `obj.ConfigGroup.IsSet()` — check if set
- `obj.ConfigGroup.Get()` — get the inner value (returns pointer)
- `obj.ConfigGroup.Set(&val)` — set the value
- `obj.ConfigGroup.Unset()` — clear the value
### ConfigId (Nullable)

Use the Nullable wrapper methods:
- `obj.ConfigId.IsSet()` — check if set
- `obj.ConfigId.Get()` — get the inner value (returns pointer)
- `obj.ConfigId.Set(&val)` — set the value
- `obj.ConfigId.Unset()` — clear the value
### ConfigRole (Nullable)

Use the Nullable wrapper methods:
- `obj.ConfigRole.IsSet()` — check if set
- `obj.ConfigRole.Get()` — get the inner value (returns pointer)
- `obj.ConfigRole.Set(&val)` — set the value
- `obj.ConfigRole.Unset()` — clear the value
### Labels (Nullable)

Use the Nullable wrapper methods:
- `obj.Labels.IsSet()` — check if set
- `obj.Labels.Get()` — get the inner value (returns pointer)
- `obj.Labels.Set(&val)` — set the value
- `obj.Labels.Unset()` — clear the value
### CoresPerSocket (Nullable)

Use the Nullable wrapper methods:
- `obj.CoresPerSocket.IsSet()` — check if set
- `obj.CoresPerSocket.Get()` — get the inner value (returns pointer)
- `obj.CoresPerSocket.Set(&val)` — set the value
- `obj.CoresPerSocket.Unset()` — clear the value
### MaxCpu (Nullable)

Use the Nullable wrapper methods:
- `obj.MaxCpu.IsSet()` — check if set
- `obj.MaxCpu.Get()` — get the inner value (returns pointer)
- `obj.MaxCpu.Set(&val)` — set the value
- `obj.MaxCpu.Unset()` — clear the value
### DomainName (Nullable)

Use the Nullable wrapper methods:
- `obj.DomainName.IsSet()` — check if set
- `obj.DomainName.Get()` — get the inner value (returns pointer)
- `obj.DomainName.Set(&val)` — set the value
- `obj.DomainName.Unset()` — clear the value
### EnvironmentPrefix (Nullable)

Use the Nullable wrapper methods:
- `obj.EnvironmentPrefix.IsSet()` — check if set
- `obj.EnvironmentPrefix.Get()` — get the inner value (returns pointer)
- `obj.EnvironmentPrefix.Set(&val)` — set the value
- `obj.EnvironmentPrefix.Unset()` — clear the value
### InstanceContext (Nullable)

Use the Nullable wrapper methods:
- `obj.InstanceContext.IsSet()` — check if set
- `obj.InstanceContext.Get()` — get the inner value (returns pointer)
- `obj.InstanceContext.Set(&val)` — set the value
- `obj.InstanceContext.Unset()` — clear the value
### CurrentDeployId (Nullable)

Use the Nullable wrapper methods:
- `obj.CurrentDeployId.IsSet()` — check if set
- `obj.CurrentDeployId.Get()` — get the inner value (returns pointer)
- `obj.CurrentDeployId.Set(&val)` — set the value
- `obj.CurrentDeployId.Unset()` — clear the value
### StatusMessage (Nullable)

Use the Nullable wrapper methods:
- `obj.StatusMessage.IsSet()` — check if set
- `obj.StatusMessage.Get()` — get the inner value (returns pointer)
- `obj.StatusMessage.Set(&val)` — set the value
- `obj.StatusMessage.Unset()` — clear the value
### ErrorMessage (Nullable)

Use the Nullable wrapper methods:
- `obj.ErrorMessage.IsSet()` — check if set
- `obj.ErrorMessage.Get()` — get the inner value (returns pointer)
- `obj.ErrorMessage.Set(&val)` — set the value
- `obj.ErrorMessage.Unset()` — clear the value
### StatusPercent (Nullable)

Use the Nullable wrapper methods:
- `obj.StatusPercent.IsSet()` — check if set
- `obj.StatusPercent.Get()` — get the inner value (returns pointer)
- `obj.StatusPercent.Set(&val)` — set the value
- `obj.StatusPercent.Unset()` — clear the value
### StatusEta (Nullable)

Use the Nullable wrapper methods:
- `obj.StatusEta.IsSet()` — check if set
- `obj.StatusEta.Get()` — get the inner value (returns pointer)
- `obj.StatusEta.Set(&val)` — set the value
- `obj.StatusEta.Unset()` — clear the value
### UserStatus (Nullable)

Use the Nullable wrapper methods:
- `obj.UserStatus.IsSet()` — check if set
- `obj.UserStatus.Get()` — get the inner value (returns pointer)
- `obj.UserStatus.Set(&val)` — set the value
- `obj.UserStatus.Unset()` — clear the value
### RemovalDate (Nullable)

Use the Nullable wrapper methods:
- `obj.RemovalDate.IsSet()` — check if set
- `obj.RemovalDate.Get()` — get the inner value (returns pointer)
- `obj.RemovalDate.Set(&val)` — set the value
- `obj.RemovalDate.Unset()` — clear the value
### Notes (Nullable)

Use the Nullable wrapper methods:
- `obj.Notes.IsSet()` — check if set
- `obj.Notes.Get()` — get the inner value (returns pointer)
- `obj.Notes.Set(&val)` — set the value
- `obj.Notes.Unset()` — clear the value
### PowerSchedule (Nullable)

Use the Nullable wrapper methods:
- `obj.PowerSchedule.IsSet()` — check if set
- `obj.PowerSchedule.Get()` — get the inner value (returns pointer)
- `obj.PowerSchedule.Set(&val)` — set the value
- `obj.PowerSchedule.Unset()` — clear the value
### InstanceThreshold (Nullable)

Use the Nullable wrapper methods:
- `obj.InstanceThreshold.IsSet()` — check if set
- `obj.InstanceThreshold.Get()` — get the inner value (returns pointer)
- `obj.InstanceThreshold.Set(&val)` — set the value
- `obj.InstanceThreshold.Unset()` — clear the value
### Apps (Nullable)

Use the Nullable wrapper methods:
- `obj.Apps.IsSet()` — check if set
- `obj.Apps.Get()` — get the inner value (returns pointer)
- `obj.Apps.Set(&val)` — set the value
- `obj.Apps.Unset()` — clear the value
### CurrentLoadBalancerInstances (Nullable)

Use the Nullable wrapper methods:
- `obj.CurrentLoadBalancerInstances.IsSet()` — check if set
- `obj.CurrentLoadBalancerInstances.Get()` — get the inner value (returns pointer)
- `obj.CurrentLoadBalancerInstances.Set(&val)` — set the value
- `obj.CurrentLoadBalancerInstances.Unset()` — clear the value

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


