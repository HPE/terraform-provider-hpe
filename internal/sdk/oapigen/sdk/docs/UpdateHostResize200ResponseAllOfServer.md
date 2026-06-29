# UpdateHostResize200ResponseAllOfServer

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | Pointer to **int64** |  | [optional] 
**Uuid** | Pointer to **string** |  | [optional] 
**ExternalId** | Pointer to **NullableString** |  | [optional] 
**InternalId** | Pointer to **NullableString** |  | [optional] 
**ExternalUniqueId** | Pointer to **NullableString** |  | [optional] 
**Name** | Pointer to **string** |  | [optional] 
**ExternalName** | Pointer to **string** |  | [optional] 
**Hostname** | Pointer to **string** |  | [optional] 
**ParentServer** | Pointer to [**UpdateHostResize200ResponseAllOfServerParentServer**](UpdateHostResize200ResponseAllOfServerParentServer.md) |  | [optional] 
**AccountId** | Pointer to **int64** |  | [optional] 
**Account** | Pointer to [**UpdateHostResize200ResponseAllOfServerAccount**](UpdateHostResize200ResponseAllOfServerAccount.md) |  | [optional] 
**Owner** | Pointer to [**UpdateHostResize200ResponseAllOfServerOwner**](UpdateHostResize200ResponseAllOfServerOwner.md) |  | [optional] 
**Zone** | Pointer to [**UpdateHostResize200ResponseAllOfServerZone**](UpdateHostResize200ResponseAllOfServerZone.md) |  | [optional] 
**Plan** | Pointer to [**UpdateHostResize200ResponseAllOfServerPlan**](UpdateHostResize200ResponseAllOfServerPlan.md) |  | [optional] 
**ComputeServerType** | Pointer to [**UpdateHostResize200ResponseAllOfServerComputeServerType**](UpdateHostResize200ResponseAllOfServerComputeServerType.md) |  | [optional] 
**Visibility** | Pointer to **string** |  | [optional] 
**Description** | Pointer to **NullableString** |  | [optional] 
**ZoneId** | Pointer to **int64** |  | [optional] 
**SiteId** | Pointer to **int64** |  | [optional] 
**ResourcePoolId** | Pointer to **NullableInt64** |  | [optional] 
**FolderId** | Pointer to **NullableInt64** |  | [optional] 
**SshHost** | Pointer to **NullableString** |  | [optional] 
**SshPort** | Pointer to **int64** |  | [optional] 
**ExternalIp** | Pointer to **NullableString** |  | [optional] 
**InternalIp** | Pointer to **NullableString** |  | [optional] 
**VolumeId** | Pointer to **NullableString** |  | [optional] 
**Platform** | Pointer to **NullableString** |  | [optional] 
**PlatformVersion** | Pointer to **NullableString** |  | [optional] 
**SshUsername** | Pointer to **NullableString** |  | [optional] 
**SshPassword** | Pointer to **NullableString** |  | [optional] 
**SshPasswordHash** | Pointer to **NullableString** |  | [optional] 
**SshKeyPair** | Pointer to [**UpdateHostResize200ResponseAllOfServerSshKeyPair**](UpdateHostResize200ResponseAllOfServerSshKeyPair.md) |  | [optional] 
**OsDevice** | Pointer to **string** |  | [optional] 
**OsType** | Pointer to **string** |  | [optional] 
**DataDevice** | Pointer to **string** |  | [optional] 
**LvmEnabled** | Pointer to **bool** |  | [optional] 
**ApiKey** | Pointer to **string** |  | [optional] 
**SoftwareRaid** | Pointer to **bool** |  | [optional] 
**DateCreated** | Pointer to **time.Time** |  | [optional] 
**LastUpdated** | Pointer to **time.Time** |  | [optional] 
**Stats** | Pointer to [**UpdateHostResize200ResponseAllOfServerStats**](UpdateHostResize200ResponseAllOfServerStats.md) |  | [optional] 
**Status** | Pointer to **string** |  | [optional] 
**StatusMessage** | Pointer to **NullableString** |  | [optional] 
**ErrorMessage** | Pointer to **NullableString** |  | [optional] 
**StatusDate** | Pointer to **NullableTime** |  | [optional] 
**StatusPercent** | Pointer to **NullableString** |  | [optional] 
**StatusEta** | Pointer to **NullableString** |  | [optional] 
**PowerState** | Pointer to **string** |  | [optional] 
**AgentInstalled** | Pointer to **bool** |  | [optional] 
**LastAgentUpdate** | Pointer to **NullableString** |  | [optional] 
**AgentVersion** | Pointer to **NullableString** |  | [optional] 
**MaxCores** | Pointer to **int64** |  | [optional] 
**CoresPerSocket** | Pointer to **NullableInt64** |  | [optional] 
**MaxMemory** | Pointer to **int64** |  | [optional] 
**MaxStorage** | Pointer to **int64** |  | [optional] 
**MaxCpu** | Pointer to **NullableInt64** |  | [optional] 
**MaxGpus** | Pointer to **NullableInt64** |  | [optional] 
**ManageInternalFirewall** | Pointer to **bool** |  | [optional] 
**EnableLogs** | Pointer to **bool** |  | [optional] 
**HourlyCost** | Pointer to **float32** |  | [optional] 
**HourlyPrice** | Pointer to **float32** |  | [optional] 
**SourceImage** | Pointer to [**UpdateHostResize200ResponseAllOfServerSourceImage**](UpdateHostResize200ResponseAllOfServerSourceImage.md) |  | [optional] 
**ServerOs** | Pointer to [**UpdateHostResize200ResponseAllOfServerServerOs**](UpdateHostResize200ResponseAllOfServerServerOs.md) |  | [optional] 
**Volumes** | Pointer to [**[]UpdateHostResize200ResponseAllOfServerVolumesInner**](UpdateHostResize200ResponseAllOfServerVolumesInner.md) |  | [optional] 
**Controllers** | Pointer to [**[]UpdateHostResize200ResponseAllOfServerControllersInner**](UpdateHostResize200ResponseAllOfServerControllersInner.md) |  | [optional] 
**Interfaces** | Pointer to [**[]UpdateHostResize200ResponseAllOfServerInterfacesInner**](UpdateHostResize200ResponseAllOfServerInterfacesInner.md) |  | [optional] 
**Labels** | Pointer to **[]map[string]interface{}** |  | [optional] 
**Tags** | Pointer to **[]map[string]interface{}** |  | [optional] 
**Enabled** | Pointer to **bool** |  | [optional] 
**TagCompliant** | Pointer to **NullableBool** |  | [optional] 
**Containers** | Pointer to **[]int64** |  | [optional] 
**Config** | Pointer to [**UpdateHostResize200ResponseAllOfServerConfig**](UpdateHostResize200ResponseAllOfServerConfig.md) |  | [optional] 
**Instance** | Pointer to [**UpdateHostResize200ResponseAllOfServerInstance**](UpdateHostResize200ResponseAllOfServerInstance.md) |  | [optional] 
**GuestConsolePreferred** | Pointer to **bool** |  | [optional] 
**GuestConsoleType** | Pointer to **NullableString** |  | [optional] 
**GuestConsoleUsername** | Pointer to **NullableString** |  | [optional] 
**GuestConsolePassword** | Pointer to **NullableString** |  | [optional] 
**GuestConsolePasswordHash** | Pointer to **NullableString** |  | [optional] 
**GuestConsolePort** | Pointer to **NullableString** |  | [optional] 
**SecureMetadataDatastore** | Pointer to [**UpdateHostResize200ResponseAllOfServerSecureMetadataDatastore**](UpdateHostResize200ResponseAllOfServerSecureMetadataDatastore.md) |  | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &UpdateHostResize200ResponseAllOfServer{
    // Set fields directly
}
```

### ExternalId (Nullable)

Use the Nullable wrapper methods:
- `obj.ExternalId.IsSet()` — check if set
- `obj.ExternalId.Get()` — get the inner value (returns pointer)
- `obj.ExternalId.Set(&val)` — set the value
- `obj.ExternalId.Unset()` — clear the value
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
### Description (Nullable)

Use the Nullable wrapper methods:
- `obj.Description.IsSet()` — check if set
- `obj.Description.Get()` — get the inner value (returns pointer)
- `obj.Description.Set(&val)` — set the value
- `obj.Description.Unset()` — clear the value
### ResourcePoolId (Nullable)

Use the Nullable wrapper methods:
- `obj.ResourcePoolId.IsSet()` — check if set
- `obj.ResourcePoolId.Get()` — get the inner value (returns pointer)
- `obj.ResourcePoolId.Set(&val)` — set the value
- `obj.ResourcePoolId.Unset()` — clear the value
### FolderId (Nullable)

Use the Nullable wrapper methods:
- `obj.FolderId.IsSet()` — check if set
- `obj.FolderId.Get()` — get the inner value (returns pointer)
- `obj.FolderId.Set(&val)` — set the value
- `obj.FolderId.Unset()` — clear the value
### SshHost (Nullable)

Use the Nullable wrapper methods:
- `obj.SshHost.IsSet()` — check if set
- `obj.SshHost.Get()` — get the inner value (returns pointer)
- `obj.SshHost.Set(&val)` — set the value
- `obj.SshHost.Unset()` — clear the value
### ExternalIp (Nullable)

Use the Nullable wrapper methods:
- `obj.ExternalIp.IsSet()` — check if set
- `obj.ExternalIp.Get()` — get the inner value (returns pointer)
- `obj.ExternalIp.Set(&val)` — set the value
- `obj.ExternalIp.Unset()` — clear the value
### InternalIp (Nullable)

Use the Nullable wrapper methods:
- `obj.InternalIp.IsSet()` — check if set
- `obj.InternalIp.Get()` — get the inner value (returns pointer)
- `obj.InternalIp.Set(&val)` — set the value
- `obj.InternalIp.Unset()` — clear the value
### VolumeId (Nullable)

Use the Nullable wrapper methods:
- `obj.VolumeId.IsSet()` — check if set
- `obj.VolumeId.Get()` — get the inner value (returns pointer)
- `obj.VolumeId.Set(&val)` — set the value
- `obj.VolumeId.Unset()` — clear the value
### Platform (Nullable)

Use the Nullable wrapper methods:
- `obj.Platform.IsSet()` — check if set
- `obj.Platform.Get()` — get the inner value (returns pointer)
- `obj.Platform.Set(&val)` — set the value
- `obj.Platform.Unset()` — clear the value
### PlatformVersion (Nullable)

Use the Nullable wrapper methods:
- `obj.PlatformVersion.IsSet()` — check if set
- `obj.PlatformVersion.Get()` — get the inner value (returns pointer)
- `obj.PlatformVersion.Set(&val)` — set the value
- `obj.PlatformVersion.Unset()` — clear the value
### SshUsername (Nullable)

Use the Nullable wrapper methods:
- `obj.SshUsername.IsSet()` — check if set
- `obj.SshUsername.Get()` — get the inner value (returns pointer)
- `obj.SshUsername.Set(&val)` — set the value
- `obj.SshUsername.Unset()` — clear the value
### SshPassword (Nullable)

Use the Nullable wrapper methods:
- `obj.SshPassword.IsSet()` — check if set
- `obj.SshPassword.Get()` — get the inner value (returns pointer)
- `obj.SshPassword.Set(&val)` — set the value
- `obj.SshPassword.Unset()` — clear the value
### SshPasswordHash (Nullable)

Use the Nullable wrapper methods:
- `obj.SshPasswordHash.IsSet()` — check if set
- `obj.SshPasswordHash.Get()` — get the inner value (returns pointer)
- `obj.SshPasswordHash.Set(&val)` — set the value
- `obj.SshPasswordHash.Unset()` — clear the value
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
### LastAgentUpdate (Nullable)

Use the Nullable wrapper methods:
- `obj.LastAgentUpdate.IsSet()` — check if set
- `obj.LastAgentUpdate.Get()` — get the inner value (returns pointer)
- `obj.LastAgentUpdate.Set(&val)` — set the value
- `obj.LastAgentUpdate.Unset()` — clear the value
### AgentVersion (Nullable)

Use the Nullable wrapper methods:
- `obj.AgentVersion.IsSet()` — check if set
- `obj.AgentVersion.Get()` — get the inner value (returns pointer)
- `obj.AgentVersion.Set(&val)` — set the value
- `obj.AgentVersion.Unset()` — clear the value
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
### MaxGpus (Nullable)

Use the Nullable wrapper methods:
- `obj.MaxGpus.IsSet()` — check if set
- `obj.MaxGpus.Get()` — get the inner value (returns pointer)
- `obj.MaxGpus.Set(&val)` — set the value
- `obj.MaxGpus.Unset()` — clear the value
### Labels (Nullable)

Use the Nullable wrapper methods:
- `obj.Labels.IsSet()` — check if set
- `obj.Labels.Get()` — get the inner value (returns pointer)
- `obj.Labels.Set(&val)` — set the value
- `obj.Labels.Unset()` — clear the value
### Tags (Nullable)

Use the Nullable wrapper methods:
- `obj.Tags.IsSet()` — check if set
- `obj.Tags.Get()` — get the inner value (returns pointer)
- `obj.Tags.Set(&val)` — set the value
- `obj.Tags.Unset()` — clear the value
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


