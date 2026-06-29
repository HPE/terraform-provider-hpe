# ClusterContainers

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | Pointer to **int64** |  | [optional] 
**Uuid** | Pointer to **string** |  | [optional] 
**AccountId** | Pointer to **int64** |  | [optional] 
**Instance** | Pointer to **NullableString** |  | [optional] 
**ContainerType** | Pointer to [**ClusterContainersContainerType**](ClusterContainersContainerType.md) |  | [optional] 
**ContainerTypeSet** | Pointer to [**ClusterContainersContainerTypeSet**](ClusterContainersContainerTypeSet.md) |  | [optional] 
**Server** | Pointer to [**ClusterContainersServer**](ClusterContainersServer.md) |  | [optional] 
**Cloud** | Pointer to [**ClusterContainersCloud**](ClusterContainersCloud.md) |  | [optional] 
**Name** | Pointer to **string** |  | [optional] 
**Ip** | Pointer to **string** |  | [optional] 
**InternalIp** | Pointer to **string** |  | [optional] 
**InternalHostname** | Pointer to **string** |  | [optional] 
**ExternalHostname** | Pointer to **string** |  | [optional] 
**ExternalDomain** | Pointer to **string** |  | [optional] 
**ExternalFqdn** | Pointer to **string** |  | [optional] 
**Ports** | Pointer to **[]map[string]interface{}** |  | [optional] 
**Plan** | Pointer to [**ClusterContainersPlan**](ClusterContainersPlan.md) |  | [optional] 
**DateCreated** | Pointer to **NullableTime** |  | [optional] 
**LastUpdated** | Pointer to **time.Time** |  | [optional] 
**StatsEnabled** | Pointer to **bool** |  | [optional] 
**Status** | Pointer to **string** |  | [optional] 
**UserStatus** | Pointer to **string** |  | [optional] 
**EnvironmentPrefix** | Pointer to **NullableString** |  | [optional] 
**ConfigGroup** | Pointer to **NullableString** |  | [optional] 
**ConfigId** | Pointer to **NullableString** |  | [optional] 
**ConfigRole** | Pointer to **NullableString** |  | [optional] 
**Stats** | Pointer to [**ClusterContainersStats**](ClusterContainersStats.md) |  | [optional] 
**RuntimeInfo** | Pointer to **map[string]interface{}** |  | [optional] 
**ContainerVersion** | Pointer to **NullableString** |  | [optional] 
**RepositoryImage** | Pointer to **NullableString** |  | [optional] 
**PlanCategory** | Pointer to **NullableString** |  | [optional] 
**Hostname** | Pointer to **NullableString** |  | [optional] 
**DomainName** | Pointer to **NullableString** |  | [optional] 
**VolumeCreated** | Pointer to **bool** |  | [optional] 
**ContainerCreated** | Pointer to **bool** |  | [optional] 
**MaxStorage** | Pointer to **NullableString** |  | [optional] 
**MaxMemory** | Pointer to **NullableString** |  | [optional] 
**MaxCores** | Pointer to **NullableString** |  | [optional] 
**MaxCpu** | Pointer to **NullableString** |  | [optional] 
**HourlyPrice** | Pointer to **float32** |  | [optional] 
**AvailableActions** | Pointer to [**[]ClusterContainersAvailableActionsInner**](ClusterContainersAvailableActionsInner.md) |  | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &ClusterContainers{
    // Set fields directly
}
```

### Instance (Nullable)

Use the Nullable wrapper methods:
- `obj.Instance.IsSet()` — check if set
- `obj.Instance.Get()` — get the inner value (returns pointer)
- `obj.Instance.Set(&val)` — set the value
- `obj.Instance.Unset()` — clear the value
### DateCreated (Nullable)

Use the Nullable wrapper methods:
- `obj.DateCreated.IsSet()` — check if set
- `obj.DateCreated.Get()` — get the inner value (returns pointer)
- `obj.DateCreated.Set(&val)` — set the value
- `obj.DateCreated.Unset()` — clear the value
### EnvironmentPrefix (Nullable)

Use the Nullable wrapper methods:
- `obj.EnvironmentPrefix.IsSet()` — check if set
- `obj.EnvironmentPrefix.Get()` — get the inner value (returns pointer)
- `obj.EnvironmentPrefix.Set(&val)` — set the value
- `obj.EnvironmentPrefix.Unset()` — clear the value
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
### Hostname (Nullable)

Use the Nullable wrapper methods:
- `obj.Hostname.IsSet()` — check if set
- `obj.Hostname.Get()` — get the inner value (returns pointer)
- `obj.Hostname.Set(&val)` — set the value
- `obj.Hostname.Unset()` — clear the value
### DomainName (Nullable)

Use the Nullable wrapper methods:
- `obj.DomainName.IsSet()` — check if set
- `obj.DomainName.Get()` — get the inner value (returns pointer)
- `obj.DomainName.Set(&val)` — set the value
- `obj.DomainName.Unset()` — clear the value
### MaxStorage (Nullable)

Use the Nullable wrapper methods:
- `obj.MaxStorage.IsSet()` — check if set
- `obj.MaxStorage.Get()` — get the inner value (returns pointer)
- `obj.MaxStorage.Set(&val)` — set the value
- `obj.MaxStorage.Unset()` — clear the value
### MaxMemory (Nullable)

Use the Nullable wrapper methods:
- `obj.MaxMemory.IsSet()` — check if set
- `obj.MaxMemory.Get()` — get the inner value (returns pointer)
- `obj.MaxMemory.Set(&val)` — set the value
- `obj.MaxMemory.Unset()` — clear the value
### MaxCores (Nullable)

Use the Nullable wrapper methods:
- `obj.MaxCores.IsSet()` — check if set
- `obj.MaxCores.Get()` — get the inner value (returns pointer)
- `obj.MaxCores.Set(&val)` — set the value
- `obj.MaxCores.Unset()` — clear the value
### MaxCpu (Nullable)

Use the Nullable wrapper methods:
- `obj.MaxCpu.IsSet()` — check if set
- `obj.MaxCpu.Get()` — get the inner value (returns pointer)
- `obj.MaxCpu.Set(&val)` — set the value
- `obj.MaxCpu.Unset()` — clear the value

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


