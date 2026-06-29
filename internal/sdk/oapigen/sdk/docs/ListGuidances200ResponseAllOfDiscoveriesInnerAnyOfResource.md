# ListGuidances200ResponseAllOfDiscoveriesInnerAnyOfResource

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
**ParentServer** | Pointer to [**NullableListGuidances200ResponseAllOfDiscoveriesInnerAnyOfResourceParentServer**](ListGuidances200ResponseAllOfDiscoveriesInnerAnyOfResourceParentServer.md) |  | [optional] 
**AccountId** | Pointer to **int64** |  | [optional] 
**Account** | Pointer to [**NullableListGuidances200ResponseAllOfDiscoveriesInnerAnyOfResourceAccount**](ListGuidances200ResponseAllOfDiscoveriesInnerAnyOfResourceAccount.md) |  | [optional] 
**Owner** | Pointer to [**ListGuidances200ResponseAllOfDiscoveriesInnerAnyOfResourceOwner**](ListGuidances200ResponseAllOfDiscoveriesInnerAnyOfResourceOwner.md) |  | [optional] 
**Zone** | Pointer to [**NullableListGuidances200ResponseAllOfDiscoveriesInnerAnyOfResourceZone**](ListGuidances200ResponseAllOfDiscoveriesInnerAnyOfResourceZone.md) |  | [optional] 
**Plan** | Pointer to [**ListGuidances200ResponseAllOfDiscoveriesInnerAnyOfResourcePlan**](ListGuidances200ResponseAllOfDiscoveriesInnerAnyOfResourcePlan.md) |  | [optional] 
**ComputeServerType** | Pointer to [**ListGuidances200ResponseAllOfDiscoveriesInnerAnyOfResourceComputeServerType**](ListGuidances200ResponseAllOfDiscoveriesInnerAnyOfResourceComputeServerType.md) |  | [optional] 
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
**Stats** | Pointer to [**ListGuidances200ResponseAllOfDiscoveriesInnerAnyOfResourceStats**](ListGuidances200ResponseAllOfDiscoveriesInnerAnyOfResourceStats.md) |  | [optional] 
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
**SourceImage** | Pointer to [**ListGuidances200ResponseAllOfDiscoveriesInnerAnyOfResourceSourceImage**](ListGuidances200ResponseAllOfDiscoveriesInnerAnyOfResourceSourceImage.md) |  | [optional] 
**ServerOs** | Pointer to [**ListGuidances200ResponseAllOfDiscoveriesInnerAnyOfResourceServerOs**](ListGuidances200ResponseAllOfDiscoveriesInnerAnyOfResourceServerOs.md) |  | [optional] 
**Volumes** | Pointer to [**[]ListGuidances200ResponseAllOfDiscoveriesInnerAnyOfResourceVolumesInner**](ListGuidances200ResponseAllOfDiscoveriesInnerAnyOfResourceVolumesInner.md) |  | [optional] 
**Controllers** | Pointer to [**[]ListGuidances200ResponseAllOfDiscoveriesInnerAnyOfResourceControllersInner**](ListGuidances200ResponseAllOfDiscoveriesInnerAnyOfResourceControllersInner.md) |  | [optional] 
**Interfaces** | Pointer to [**[]ListGuidances200ResponseAllOfDiscoveriesInnerAnyOfResourceInterfacesInner**](ListGuidances200ResponseAllOfDiscoveriesInnerAnyOfResourceInterfacesInner.md) |  | [optional] 
**Labels** | Pointer to **[]map[string]interface{}** |  | [optional] 
**Tags** | Pointer to **[]map[string]interface{}** |  | [optional] 
**Enabled** | Pointer to **bool** |  | [optional] 
**TagCompliant** | Pointer to **NullableBool** |  | [optional] 
**Containers** | Pointer to **[]int64** |  | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &ListGuidances200ResponseAllOfDiscoveriesInnerAnyOfResource{
    // Set fields directly
}
```

### ParentServer (Nullable)

Use the Nullable wrapper methods:
- `obj.ParentServer.IsSet()` — check if set
- `obj.ParentServer.Get()` — get the inner value (returns pointer)
- `obj.ParentServer.Set(&val)` — set the value
- `obj.ParentServer.Unset()` — clear the value
### Account (Nullable)

Use the Nullable wrapper methods:
- `obj.Account.IsSet()` — check if set
- `obj.Account.Get()` — get the inner value (returns pointer)
- `obj.Account.Set(&val)` — set the value
- `obj.Account.Unset()` — clear the value
### Zone (Nullable)

Use the Nullable wrapper methods:
- `obj.Zone.IsSet()` — check if set
- `obj.Zone.Get()` — get the inner value (returns pointer)
- `obj.Zone.Set(&val)` — set the value
- `obj.Zone.Unset()` — clear the value
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


