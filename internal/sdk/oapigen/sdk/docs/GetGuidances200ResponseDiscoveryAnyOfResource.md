# GetGuidances200ResponseDiscoveryAnyOfResource

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | Pointer to **int64** |  | [optional] 
**Uuid** | Pointer to **string** |  | [optional] 
**ExternalId** | Pointer to **string** |  | [optional] 
**InternalId** | Pointer to **string** |  | [optional] 
**ExternalUniqueId** | Pointer to **string** |  | [optional] 
**Name** | Pointer to **string** |  | [optional] 
**ExternalName** | Pointer to **string** |  | [optional] 
**Hostname** | Pointer to **string** |  | [optional] 
**ParentServer** | Pointer to [**GetGuidances200ResponseDiscoveryAnyOfResourceParentServer**](GetGuidances200ResponseDiscoveryAnyOfResourceParentServer.md) |  | [optional] 
**AccountId** | Pointer to **int64** |  | [optional] 
**Account** | Pointer to [**GetGuidances200ResponseDiscoveryAnyOfResourceAccount**](GetGuidances200ResponseDiscoveryAnyOfResourceAccount.md) |  | [optional] 
**Owner** | Pointer to [**GetGuidances200ResponseDiscoveryAnyOfResourceOwner**](GetGuidances200ResponseDiscoveryAnyOfResourceOwner.md) |  | [optional] 
**Zone** | Pointer to [**GetGuidances200ResponseDiscoveryAnyOfResourceZone**](GetGuidances200ResponseDiscoveryAnyOfResourceZone.md) |  | [optional] 
**Plan** | Pointer to [**GetGuidances200ResponseDiscoveryAnyOfResourcePlan**](GetGuidances200ResponseDiscoveryAnyOfResourcePlan.md) |  | [optional] 
**ComputeServerType** | Pointer to [**GetGuidances200ResponseDiscoveryAnyOfResourceComputeServerType**](GetGuidances200ResponseDiscoveryAnyOfResourceComputeServerType.md) |  | [optional] 
**Visibility** | Pointer to **string** |  | [optional] 
**Description** | Pointer to **NullableString** |  | [optional] 
**ZoneId** | Pointer to **int64** |  | [optional] 
**SiteId** | Pointer to **int64** |  | [optional] 
**ResourcePoolId** | Pointer to **int64** |  | [optional] 
**FolderId** | Pointer to **int64** |  | [optional] 
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
**Stats** | Pointer to [**GetGuidances200ResponseDiscoveryAnyOfResourceStats**](GetGuidances200ResponseDiscoveryAnyOfResourceStats.md) |  | [optional] 
**Status** | Pointer to **string** |  | [optional] 
**StatusMessage** | Pointer to **NullableString** |  | [optional] 
**ErrorMessage** | Pointer to **NullableString** |  | [optional] 
**StatusDate** | Pointer to **NullableTime** |  | [optional] 
**StatusPercent** | Pointer to **NullableString** |  | [optional] 
**StatusEta** | Pointer to **NullableString** |  | [optional] 
**PowerState** | Pointer to **string** |  | [optional] 
**AgentInstalled** | Pointer to **bool** |  | [optional] 
**LastAgentUpdate** | Pointer to **time.Time** |  | [optional] 
**AgentVersion** | Pointer to **string** |  | [optional] 
**MaxCores** | Pointer to **int64** |  | [optional] 
**CoresPerSocket** | Pointer to **int64** |  | [optional] 
**MaxMemory** | Pointer to **int64** |  | [optional] 
**MaxStorage** | Pointer to **int64** |  | [optional] 
**MaxCpu** | Pointer to **NullableString** |  | [optional] 
**HourlyPrice** | Pointer to **float32** |  | [optional] 
**SourceImage** | Pointer to [**GetGuidances200ResponseDiscoveryAnyOfResourceSourceImage**](GetGuidances200ResponseDiscoveryAnyOfResourceSourceImage.md) |  | [optional] 
**ServerOs** | Pointer to [**GetGuidances200ResponseDiscoveryAnyOfResourceServerOs**](GetGuidances200ResponseDiscoveryAnyOfResourceServerOs.md) |  | [optional] 
**Volumes** | Pointer to [**[]GetGuidances200ResponseDiscoveryAnyOfResourceVolumesInner**](GetGuidances200ResponseDiscoveryAnyOfResourceVolumesInner.md) |  | [optional] 
**Controllers** | Pointer to [**[]GetGuidances200ResponseDiscoveryAnyOfResourceControllersInner**](GetGuidances200ResponseDiscoveryAnyOfResourceControllersInner.md) |  | [optional] 
**Interfaces** | Pointer to [**[]GetGuidances200ResponseDiscoveryAnyOfResourceInterfacesInner**](GetGuidances200ResponseDiscoveryAnyOfResourceInterfacesInner.md) |  | [optional] 
**Labels** | Pointer to **[]map[string]interface{}** |  | [optional] 
**Tags** | Pointer to **[]map[string]interface{}** |  | [optional] 
**Enabled** | Pointer to **bool** |  | [optional] 
**TagCompliant** | Pointer to **NullableBool** |  | [optional] 
**Containers** | Pointer to **[]int64** |  | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &GetGuidances200ResponseDiscoveryAnyOfResource{
    // Set fields directly
}
```

### Description (Nullable)

Use the Nullable wrapper methods:
- `obj.Description.IsSet()` — check if set
- `obj.Description.Get()` — get the inner value (returns pointer)
- `obj.Description.Set(&val)` — set the value
- `obj.Description.Unset()` — clear the value
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
### MaxCpu (Nullable)

Use the Nullable wrapper methods:
- `obj.MaxCpu.IsSet()` — check if set
- `obj.MaxCpu.Get()` — get the inner value (returns pointer)
- `obj.MaxCpu.Set(&val)` — set the value
- `obj.MaxCpu.Unset()` — clear the value
### TagCompliant (Nullable)

Use the Nullable wrapper methods:
- `obj.TagCompliant.IsSet()` — check if set
- `obj.TagCompliant.Get()` — get the inner value (returns pointer)
- `obj.TagCompliant.Set(&val)` — set the value
- `obj.TagCompliant.Unset()` — clear the value

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


