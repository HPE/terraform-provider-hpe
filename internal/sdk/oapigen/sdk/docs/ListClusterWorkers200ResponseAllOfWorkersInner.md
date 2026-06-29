# ListClusterWorkers200ResponseAllOfWorkersInner

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | Pointer to **int64** |  | [optional] 
**Uuid** | Pointer to **string** |  | [optional] 
**ExternalId** | Pointer to **string** |  | [optional] 
**InternalId** | Pointer to **NullableString** |  | [optional] 
**ExternalUniqueId** | Pointer to **NullableString** |  | [optional] 
**Name** | Pointer to **string** |  | [optional] 
**ExternalName** | Pointer to **string** |  | [optional] 
**Hostname** | Pointer to **string** |  | [optional] 
**AccountId** | Pointer to **int64** |  | [optional] 
**Account** | Pointer to [**ListClusterWorkers200ResponseAllOfWorkersInnerAccount**](ListClusterWorkers200ResponseAllOfWorkersInnerAccount.md) |  | [optional] 
**Owner** | Pointer to [**ListClusterWorkers200ResponseAllOfWorkersInnerOwner**](ListClusterWorkers200ResponseAllOfWorkersInnerOwner.md) |  | [optional] 
**Zone** | Pointer to [**ListClusterWorkers200ResponseAllOfWorkersInnerZone**](ListClusterWorkers200ResponseAllOfWorkersInnerZone.md) |  | [optional] 
**Plan** | Pointer to [**ListClusterWorkers200ResponseAllOfWorkersInnerPlan**](ListClusterWorkers200ResponseAllOfWorkersInnerPlan.md) |  | [optional] 
**ComputeServerType** | Pointer to [**ListClusterWorkers200ResponseAllOfWorkersInnerComputeServerType**](ListClusterWorkers200ResponseAllOfWorkersInnerComputeServerType.md) |  | [optional] 
**Visibility** | Pointer to **string** |  | [optional] 
**Description** | Pointer to **string** |  | [optional] 
**ZoneId** | Pointer to **int64** |  | [optional] 
**SiteId** | Pointer to **int64** |  | [optional] 
**ResourcePoolId** | Pointer to **int64** |  | [optional] 
**FolderId** | Pointer to **NullableString** |  | [optional] 
**SshHost** | Pointer to **string** |  | [optional] 
**SshPort** | Pointer to **int64** |  | [optional] 
**ExternalIp** | Pointer to **string** |  | [optional] 
**InternalIp** | Pointer to **string** |  | [optional] 
**VolumeId** | Pointer to **NullableString** |  | [optional] 
**Platform** | Pointer to **string** |  | [optional] 
**PlatformVersion** | Pointer to **string** |  | [optional] 
**SshUsername** | Pointer to **string** |  | [optional] 
**SshPassword** | Pointer to **string** |  | [optional] 
**SshPasswordHash** | Pointer to **string** |  | [optional] 
**OsDevice** | Pointer to **string** |  | [optional] 
**OsType** | Pointer to **string** |  | [optional] 
**DataDevice** | Pointer to **string** |  | [optional] 
**LvmEnabled** | Pointer to **bool** |  | [optional] 
**ApiKey** | Pointer to **string** |  | [optional] 
**SoftwareRaid** | Pointer to **bool** |  | [optional] 
**DateCreated** | Pointer to **time.Time** |  | [optional] 
**LastUpdated** | Pointer to **time.Time** |  | [optional] 
**Stats** | Pointer to [**ListClusterWorkers200ResponseAllOfWorkersInnerStats**](ListClusterWorkers200ResponseAllOfWorkersInnerStats.md) |  | [optional] 
**Status** | Pointer to **string** |  | [optional] 
**StatusMessage** | Pointer to **NullableString** |  | [optional] 
**ErrorMessage** | Pointer to **NullableString** |  | [optional] 
**StatusDate** | Pointer to **NullableString** |  | [optional] 
**StatusPercent** | Pointer to **NullableString** |  | [optional] 
**StatusEta** | Pointer to **NullableString** |  | [optional] 
**PowerState** | Pointer to **string** |  | [optional] 
**AgentInstalled** | Pointer to **bool** |  | [optional] 
**LastAgentUpdate** | Pointer to **time.Time** |  | [optional] 
**AgentVersion** | Pointer to **string** |  | [optional] 
**MaxCores** | Pointer to **int64** |  | [optional] 
**CoresPerSocket** | Pointer to **NullableString** |  | [optional] 
**MaxMemory** | Pointer to **int64** |  | [optional] 
**MaxStorage** | Pointer to **int64** |  | [optional] 
**MaxCpu** | Pointer to **int64** |  | [optional] 
**HourlyPrice** | Pointer to **float32** |  | [optional] 
**SourceImage** | Pointer to [**ListClusterWorkers200ResponseAllOfWorkersInnerSourceImage**](ListClusterWorkers200ResponseAllOfWorkersInnerSourceImage.md) |  | [optional] 
**ServerOs** | Pointer to **NullableString** |  | [optional] 
**Volumes** | Pointer to [**[]ListClusterWorkers200ResponseAllOfWorkersInnerVolumesInner**](ListClusterWorkers200ResponseAllOfWorkersInnerVolumesInner.md) |  | [optional] 
**Controllers** | Pointer to **[]map[string]interface{}** |  | [optional] 
**Interfaces** | Pointer to [**[]ListClusterWorkers200ResponseAllOfWorkersInnerInterfacesInner**](ListClusterWorkers200ResponseAllOfWorkersInnerInterfacesInner.md) |  | [optional] 
**Labels** | Pointer to **[]string** |  | [optional] 
**Tags** | Pointer to [**[]ListClusterWorkers200ResponseAllOfWorkersInnerTagsInner**](ListClusterWorkers200ResponseAllOfWorkersInnerTagsInner.md) |  | [optional] 
**Enabled** | Pointer to **bool** |  | [optional] 
**TagCompliant** | Pointer to **NullableBool** |  | [optional] 
**Containers** | Pointer to **[]int64** |  | [optional] 
**GuestConsolePreferred** | Pointer to **bool** |  | [optional] 
**GuestConsoleType** | Pointer to **NullableString** |  | [optional] 
**GuestConsoleUsername** | Pointer to **NullableString** |  | [optional] 
**GuestConsolePassword** | Pointer to **NullableString** |  | [optional] 
**GuestConsolePasswordHash** | Pointer to **NullableString** |  | [optional] 
**GuestConsolePort** | Pointer to **NullableString** |  | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &ListClusterWorkers200ResponseAllOfWorkersInner{
    // Set fields directly
}
```

### InternalId (Nullable)

Use the Nullable wrapper methods:
- `obj.InternalId.IsSet()` — check if set
- `obj.InternalId.Get()` — get the inner value (returns pointer)
- `obj.InternalId.Set(&val)` — set the value
- `obj.InternalId.Unset()` — clear the value
### ExternalUniqueId (Nullable)

Use the Nullable wrapper methods:
- `obj.ExternalUniqueId.IsSet()` — check if set
- `obj.ExternalUniqueId.Get()` — get the inner value (returns pointer)
- `obj.ExternalUniqueId.Set(&val)` — set the value
- `obj.ExternalUniqueId.Unset()` — clear the value
### FolderId (Nullable)

Use the Nullable wrapper methods:
- `obj.FolderId.IsSet()` — check if set
- `obj.FolderId.Get()` — get the inner value (returns pointer)
- `obj.FolderId.Set(&val)` — set the value
- `obj.FolderId.Unset()` — clear the value
### VolumeId (Nullable)

Use the Nullable wrapper methods:
- `obj.VolumeId.IsSet()` — check if set
- `obj.VolumeId.Get()` — get the inner value (returns pointer)
- `obj.VolumeId.Set(&val)` — set the value
- `obj.VolumeId.Unset()` — clear the value
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
### StatusDate (Nullable)

Use the Nullable wrapper methods:
- `obj.StatusDate.IsSet()` — check if set
- `obj.StatusDate.Get()` — get the inner value (returns pointer)
- `obj.StatusDate.Set(&val)` — set the value
- `obj.StatusDate.Unset()` — clear the value
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
### CoresPerSocket (Nullable)

Use the Nullable wrapper methods:
- `obj.CoresPerSocket.IsSet()` — check if set
- `obj.CoresPerSocket.Get()` — get the inner value (returns pointer)
- `obj.CoresPerSocket.Set(&val)` — set the value
- `obj.CoresPerSocket.Unset()` — clear the value
### ServerOs (Nullable)

Use the Nullable wrapper methods:
- `obj.ServerOs.IsSet()` — check if set
- `obj.ServerOs.Get()` — get the inner value (returns pointer)
- `obj.ServerOs.Set(&val)` — set the value
- `obj.ServerOs.Unset()` — clear the value
### TagCompliant (Nullable)

Use the Nullable wrapper methods:
- `obj.TagCompliant.IsSet()` — check if set
- `obj.TagCompliant.Get()` — get the inner value (returns pointer)
- `obj.TagCompliant.Set(&val)` — set the value
- `obj.TagCompliant.Unset()` — clear the value
### GuestConsoleType (Nullable)

Use the Nullable wrapper methods:
- `obj.GuestConsoleType.IsSet()` — check if set
- `obj.GuestConsoleType.Get()` — get the inner value (returns pointer)
- `obj.GuestConsoleType.Set(&val)` — set the value
- `obj.GuestConsoleType.Unset()` — clear the value
### GuestConsoleUsername (Nullable)

Use the Nullable wrapper methods:
- `obj.GuestConsoleUsername.IsSet()` — check if set
- `obj.GuestConsoleUsername.Get()` — get the inner value (returns pointer)
- `obj.GuestConsoleUsername.Set(&val)` — set the value
- `obj.GuestConsoleUsername.Unset()` — clear the value
### GuestConsolePassword (Nullable)

Use the Nullable wrapper methods:
- `obj.GuestConsolePassword.IsSet()` — check if set
- `obj.GuestConsolePassword.Get()` — get the inner value (returns pointer)
- `obj.GuestConsolePassword.Set(&val)` — set the value
- `obj.GuestConsolePassword.Unset()` — clear the value
### GuestConsolePasswordHash (Nullable)

Use the Nullable wrapper methods:
- `obj.GuestConsolePasswordHash.IsSet()` — check if set
- `obj.GuestConsolePasswordHash.Get()` — get the inner value (returns pointer)
- `obj.GuestConsolePasswordHash.Set(&val)` — set the value
- `obj.GuestConsolePasswordHash.Unset()` — clear the value
### GuestConsolePort (Nullable)

Use the Nullable wrapper methods:
- `obj.GuestConsolePort.IsSet()` — check if set
- `obj.GuestConsolePort.Get()` — get the inner value (returns pointer)
- `obj.GuestConsolePort.Set(&val)` — set the value
- `obj.GuestConsolePort.Unset()` — clear the value

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


