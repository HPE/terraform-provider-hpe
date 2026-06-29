# GetContainer200ResponseContainer

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | Pointer to **int32** |  | [optional] 
**Uuid** | Pointer to **string** |  | [optional] 
**AccountId** | Pointer to **int32** |  | [optional] 
**Instance** | Pointer to [**GetContainer200ResponseContainerInstance**](GetContainer200ResponseContainerInstance.md) |  | [optional] 
**ContainerType** | Pointer to [**GetContainer200ResponseContainerContainerType**](GetContainer200ResponseContainerContainerType.md) |  | [optional] 
**ContainerTypeSet** | Pointer to [**GetContainer200ResponseContainerContainerTypeSet**](GetContainer200ResponseContainerContainerTypeSet.md) |  | [optional] 
**Server** | Pointer to [**GetContainer200ResponseContainerServer**](GetContainer200ResponseContainerServer.md) |  | [optional] 
**Cloud** | Pointer to [**GetContainer200ResponseContainerCloud**](GetContainer200ResponseContainerCloud.md) |  | [optional] 
**Name** | Pointer to **string** |  | [optional] 
**Ip** | Pointer to **string** |  | [optional] 
**InternalIp** | Pointer to **string** |  | [optional] 
**InternalHostname** | Pointer to **string** |  | [optional] 
**ExternalHostname** | Pointer to **string** |  | [optional] 
**ExternalDomain** | Pointer to **string** |  | [optional] 
**ExternalFqdn** | Pointer to **string** |  | [optional] 
**Ports** | Pointer to [**[]GetContainer200ResponseContainerPortsInner**](GetContainer200ResponseContainerPortsInner.md) |  | [optional] 
**Plan** | Pointer to [**GetContainer200ResponseContainerPlan**](GetContainer200ResponseContainerPlan.md) |  | [optional] 
**DateCreated** | Pointer to **time.Time** |  | [optional] 
**LastUpdated** | Pointer to **time.Time** |  | [optional] 
**StatsEnabled** | Pointer to **bool** |  | [optional] 
**Status** | Pointer to **string** |  | [optional] 
**UserStatus** | Pointer to **string** |  | [optional] 
**EnvironmentPrefix** | Pointer to **NullableString** |  | [optional] 
**Stats** | Pointer to [**GetContainer200ResponseContainerStats**](GetContainer200ResponseContainerStats.md) |  | [optional] 
**RuntimeInfo** | Pointer to **map[string]interface{}** |  | [optional] 
**ContainerVersion** | Pointer to **NullableString** |  | [optional] 
**RepositoryImage** | Pointer to **NullableString** |  | [optional] 
**PlanCategory** | Pointer to **NullableString** |  | [optional] 
**Hostname** | Pointer to **string** |  | [optional] 
**DomainName** | Pointer to **NullableString** |  | [optional] 
**VolumeCreated** | Pointer to **bool** |  | [optional] 
**ContainerCreated** | Pointer to **bool** |  | [optional] 
**MaxStorage** | Pointer to **int32** |  | [optional] 
**MaxMemory** | Pointer to **int32** |  | [optional] 
**MaxCores** | Pointer to **int32** |  | [optional] 
**MaxCpu** | Pointer to **NullableInt32** |  | [optional] 
**AvailableActions** | Pointer to **[]map[string]interface{}** |  | [optional] 
**ConfigGroup** | Pointer to **NullableString** |  | [optional] 
**ConfigId** | Pointer to **NullableString** |  | [optional] 
**ConfigRole** | Pointer to **NullableString** |  | [optional] 
**HourlyCost** | Pointer to **float64** |  | [optional] 
**HourlyPrice** | Pointer to **float64** |  | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &GetContainer200ResponseContainer{
    // Set fields directly
}
```

### Ports (Nullable)

Use the Nullable wrapper methods:
- `obj.Ports.IsSet()` — check if set
- `obj.Ports.Get()` — get the inner value (returns pointer)
- `obj.Ports.Set(&val)` — set the value
- `obj.Ports.Unset()` — clear the value
### EnvironmentPrefix (Nullable)

Use the Nullable wrapper methods:
- `obj.EnvironmentPrefix.IsSet()` — check if set
- `obj.EnvironmentPrefix.Get()` — get the inner value (returns pointer)
- `obj.EnvironmentPrefix.Set(&val)` — set the value
- `obj.EnvironmentPrefix.Unset()` — clear the value
### ContainerVersion (Nullable)

Use the Nullable wrapper methods:
- `obj.ContainerVersion.IsSet()` — check if set
- `obj.ContainerVersion.Get()` — get the inner value (returns pointer)
- `obj.ContainerVersion.Set(&val)` — set the value
- `obj.ContainerVersion.Unset()` — clear the value
### RepositoryImage (Nullable)

Use the Nullable wrapper methods:
- `obj.RepositoryImage.IsSet()` — check if set
- `obj.RepositoryImage.Get()` — get the inner value (returns pointer)
- `obj.RepositoryImage.Set(&val)` — set the value
- `obj.RepositoryImage.Unset()` — clear the value
### PlanCategory (Nullable)

Use the Nullable wrapper methods:
- `obj.PlanCategory.IsSet()` — check if set
- `obj.PlanCategory.Get()` — get the inner value (returns pointer)
- `obj.PlanCategory.Set(&val)` — set the value
- `obj.PlanCategory.Unset()` — clear the value
### DomainName (Nullable)

Use the Nullable wrapper methods:
- `obj.DomainName.IsSet()` — check if set
- `obj.DomainName.Get()` — get the inner value (returns pointer)
- `obj.DomainName.Set(&val)` — set the value
- `obj.DomainName.Unset()` — clear the value
### MaxCpu (Nullable)

Use the Nullable wrapper methods:
- `obj.MaxCpu.IsSet()` — check if set
- `obj.MaxCpu.Get()` — get the inner value (returns pointer)
- `obj.MaxCpu.Set(&val)` — set the value
- `obj.MaxCpu.Unset()` — clear the value
### AvailableActions (Nullable)

Use the Nullable wrapper methods:
- `obj.AvailableActions.IsSet()` — check if set
- `obj.AvailableActions.Get()` — get the inner value (returns pointer)
- `obj.AvailableActions.Set(&val)` — set the value
- `obj.AvailableActions.Unset()` — clear the value
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

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


