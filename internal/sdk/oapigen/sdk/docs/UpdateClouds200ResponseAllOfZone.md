# UpdateClouds200ResponseAllOfZone

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | Pointer to **int64** |  | [optional] 
**Uuid** | Pointer to **string** |  | [optional] 
**ExternalId** | Pointer to **NullableString** |  | [optional] 
**Name** | Pointer to **string** |  | [optional] 
**Code** | Pointer to **string** |  | [optional] 
**Labels** | Pointer to **[]string** |  | [optional] 
**Location** | Pointer to **NullableString** |  | [optional] 
**Owner** | Pointer to [**NullableUpdateClouds200ResponseAllOfZoneOwner**](UpdateClouds200ResponseAllOfZoneOwner.md) |  | [optional] 
**AccountId** | Pointer to **int64** |  | [optional] 
**Account** | Pointer to [**NullableUpdateClouds200ResponseAllOfZoneAccount**](UpdateClouds200ResponseAllOfZoneAccount.md) |  | [optional] 
**Visibility** | Pointer to **string** |  | [optional] 
**Enabled** | Pointer to **bool** |  | [optional] 
**Status** | Pointer to **string** |  | [optional] 
**StatusMessage** | Pointer to **NullableString** |  | [optional] 
**StatusDate** | Pointer to **NullableTime** |  | [optional] 
**CostStatus** | Pointer to **NullableString** |  | [optional] 
**CostStatusMessage** | Pointer to **NullableString** |  | [optional] 
**CostStatusDate** | Pointer to **NullableTime** |  | [optional] 
**CostLastSyncDuration** | Pointer to **NullableInt64** |  | [optional] 
**CostLastSync** | Pointer to **NullableTime** |  | [optional] 
**ZoneType** | Pointer to [**UpdateClouds200ResponseAllOfZoneZoneType**](UpdateClouds200ResponseAllOfZoneZoneType.md) |  | [optional] 
**ZoneTypeId** | Pointer to **int64** |  | [optional] 
**GuidanceMode** | Pointer to **NullableString** |  | [optional] 
**StorageMode** | Pointer to **string** |  | [optional] 
**AgentMode** | Pointer to **string** |  | [optional] 
**UserDataLinux** | Pointer to **NullableString** |  | [optional] 
**UserDataWindows** | Pointer to **NullableString** |  | [optional] 
**ConsoleKeymap** | Pointer to **NullableString** |  | [optional] 
**ContainerMode** | Pointer to **string** |  | [optional] 
**CostingMode** | Pointer to **NullableString** |  | [optional] 
**ServiceVersion** | Pointer to **NullableString** |  | [optional] 
**SecurityMode** | Pointer to **string** |  | [optional] 
**InventoryLevel** | Pointer to **string** |  | [optional] 
**Timezone** | Pointer to **NullableString** |  | [optional] 
**ApiProxy** | Pointer to **NullableString** |  | [optional] 
**ProvisioningProxy** | Pointer to **NullableString** |  | [optional] 
**NetworkDomain** | Pointer to [**UpdateClouds200ResponseAllOfZoneNetworkDomain**](UpdateClouds200ResponseAllOfZoneNetworkDomain.md) |  | [optional] 
**DomainName** | Pointer to **string** |  | [optional] 
**RegionCode** | Pointer to **NullableString** |  | [optional] 
**AutoRecoverPowerState** | Pointer to **bool** |  | [optional] 
**ScalePriority** | Pointer to **int64** |  | [optional] 
**DefaultDatastoreSyncActive** | Pointer to **bool** |  | [optional] 
**DefaultNetworkSyncActive** | Pointer to **bool** |  | [optional] 
**DefaultFolderSyncActive** | Pointer to **bool** |  | [optional] 
**DefaultSecurityGroupSyncActive** | Pointer to **bool** |  | [optional] 
**DefaultPoolSyncActive** | Pointer to **bool** |  | [optional] 
**DefaultPlanSyncActive** | Pointer to **bool** |  | [optional] 
**Config** | Pointer to [**UpdateClouds200ResponseAllOfZoneConfig**](UpdateClouds200ResponseAllOfZoneConfig.md) |  | [optional] 
**Credential** | Pointer to [**UpdateClouds200ResponseAllOfZoneCredential**](UpdateClouds200ResponseAllOfZoneCredential.md) |  | [optional] 
**ImagePath** | Pointer to **NullableString** | Logo image URL | [optional] 
**DarkImagePath** | Pointer to **NullableString** | Dark logo image URL | [optional] 
**DateCreated** | Pointer to **time.Time** |  | [optional] 
**LastUpdated** | Pointer to **time.Time** |  | [optional] 
**LastSync** | Pointer to **NullableTime** |  | [optional] 
**LastSyncDuration** | Pointer to **NullableInt64** |  | [optional] 
**NextRunDate** | Pointer to **NullableTime** |  | [optional] 
**Groups** | Pointer to [**[]UpdateClouds200ResponseAllOfZoneGroupsInner**](UpdateClouds200ResponseAllOfZoneGroupsInner.md) |  | [optional] 
**SecurityServer** | Pointer to [**UpdateClouds200ResponseAllOfZoneSecurityServer**](UpdateClouds200ResponseAllOfZoneSecurityServer.md) |  | [optional] 
**NetworkServer** | Pointer to [**UpdateClouds200ResponseAllOfZoneNetworkServer**](UpdateClouds200ResponseAllOfZoneNetworkServer.md) |  | [optional] 
**Stats** | Pointer to [**UpdateClouds200ResponseAllOfZoneStats**](UpdateClouds200ResponseAllOfZoneStats.md) |  | [optional] 
**ServerCount** | Pointer to **int64** |  | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &UpdateClouds200ResponseAllOfZone{
    // Set fields directly
}
```

### ExternalId (Nullable)

Use the Nullable wrapper methods:
- `obj.ExternalId.IsSet()` — check if set
- `obj.ExternalId.Get()` — get the inner value (returns pointer)
- `obj.ExternalId.Set(&val)` — set the value
- `obj.ExternalId.Unset()` — clear the value
### Location (Nullable)

Use the Nullable wrapper methods:
- `obj.Location.IsSet()` — check if set
- `obj.Location.Get()` — get the inner value (returns pointer)
- `obj.Location.Set(&val)` — set the value
- `obj.Location.Unset()` — clear the value
### Owner (Nullable)

Use the Nullable wrapper methods:
- `obj.Owner.IsSet()` — check if set
- `obj.Owner.Get()` — get the inner value (returns pointer)
- `obj.Owner.Set(&val)` — set the value
- `obj.Owner.Unset()` — clear the value
### Account (Nullable)

Use the Nullable wrapper methods:
- `obj.Account.IsSet()` — check if set
- `obj.Account.Get()` — get the inner value (returns pointer)
- `obj.Account.Set(&val)` — set the value
- `obj.Account.Unset()` — clear the value
### StatusMessage (Nullable)

Use the Nullable wrapper methods:
- `obj.StatusMessage.IsSet()` — check if set
- `obj.StatusMessage.Get()` — get the inner value (returns pointer)
- `obj.StatusMessage.Set(&val)` — set the value
- `obj.StatusMessage.Unset()` — clear the value
### StatusDate (Nullable)

Use the Nullable wrapper methods:
- `obj.StatusDate.IsSet()` — check if set
- `obj.StatusDate.Get()` — get the inner value (returns pointer)
- `obj.StatusDate.Set(&val)` — set the value
- `obj.StatusDate.Unset()` — clear the value
### CostStatus (Nullable)

Use the Nullable wrapper methods:
- `obj.CostStatus.IsSet()` — check if set
- `obj.CostStatus.Get()` — get the inner value (returns pointer)
- `obj.CostStatus.Set(&val)` — set the value
- `obj.CostStatus.Unset()` — clear the value
### CostStatusMessage (Nullable)

Use the Nullable wrapper methods:
- `obj.CostStatusMessage.IsSet()` — check if set
- `obj.CostStatusMessage.Get()` — get the inner value (returns pointer)
- `obj.CostStatusMessage.Set(&val)` — set the value
- `obj.CostStatusMessage.Unset()` — clear the value
### CostStatusDate (Nullable)

Use the Nullable wrapper methods:
- `obj.CostStatusDate.IsSet()` — check if set
- `obj.CostStatusDate.Get()` — get the inner value (returns pointer)
- `obj.CostStatusDate.Set(&val)` — set the value
- `obj.CostStatusDate.Unset()` — clear the value
### CostLastSyncDuration (Nullable)

Use the Nullable wrapper methods:
- `obj.CostLastSyncDuration.IsSet()` — check if set
- `obj.CostLastSyncDuration.Get()` — get the inner value (returns pointer)
- `obj.CostLastSyncDuration.Set(&val)` — set the value
- `obj.CostLastSyncDuration.Unset()` — clear the value
### CostLastSync (Nullable)

Use the Nullable wrapper methods:
- `obj.CostLastSync.IsSet()` — check if set
- `obj.CostLastSync.Get()` — get the inner value (returns pointer)
- `obj.CostLastSync.Set(&val)` — set the value
- `obj.CostLastSync.Unset()` — clear the value
### GuidanceMode (Nullable)

Use the Nullable wrapper methods:
- `obj.GuidanceMode.IsSet()` — check if set
- `obj.GuidanceMode.Get()` — get the inner value (returns pointer)
- `obj.GuidanceMode.Set(&val)` — set the value
- `obj.GuidanceMode.Unset()` — clear the value
### UserDataLinux (Nullable)

Use the Nullable wrapper methods:
- `obj.UserDataLinux.IsSet()` — check if set
- `obj.UserDataLinux.Get()` — get the inner value (returns pointer)
- `obj.UserDataLinux.Set(&val)` — set the value
- `obj.UserDataLinux.Unset()` — clear the value
### UserDataWindows (Nullable)

Use the Nullable wrapper methods:
- `obj.UserDataWindows.IsSet()` — check if set
- `obj.UserDataWindows.Get()` — get the inner value (returns pointer)
- `obj.UserDataWindows.Set(&val)` — set the value
- `obj.UserDataWindows.Unset()` — clear the value
### ConsoleKeymap (Nullable)

Use the Nullable wrapper methods:
- `obj.ConsoleKeymap.IsSet()` — check if set
- `obj.ConsoleKeymap.Get()` — get the inner value (returns pointer)
- `obj.ConsoleKeymap.Set(&val)` — set the value
- `obj.ConsoleKeymap.Unset()` — clear the value
### CostingMode (Nullable)

Use the Nullable wrapper methods:
- `obj.CostingMode.IsSet()` — check if set
- `obj.CostingMode.Get()` — get the inner value (returns pointer)
- `obj.CostingMode.Set(&val)` — set the value
- `obj.CostingMode.Unset()` — clear the value
### ServiceVersion (Nullable)

Use the Nullable wrapper methods:
- `obj.ServiceVersion.IsSet()` — check if set
- `obj.ServiceVersion.Get()` — get the inner value (returns pointer)
- `obj.ServiceVersion.Set(&val)` — set the value
- `obj.ServiceVersion.Unset()` — clear the value
### Timezone (Nullable)

Use the Nullable wrapper methods:
- `obj.Timezone.IsSet()` — check if set
- `obj.Timezone.Get()` — get the inner value (returns pointer)
- `obj.Timezone.Set(&val)` — set the value
- `obj.Timezone.Unset()` — clear the value
### ApiProxy (Nullable)

Use the Nullable wrapper methods:
- `obj.ApiProxy.IsSet()` — check if set
- `obj.ApiProxy.Get()` — get the inner value (returns pointer)
- `obj.ApiProxy.Set(&val)` — set the value
- `obj.ApiProxy.Unset()` — clear the value
### ProvisioningProxy (Nullable)

Use the Nullable wrapper methods:
- `obj.ProvisioningProxy.IsSet()` — check if set
- `obj.ProvisioningProxy.Get()` — get the inner value (returns pointer)
- `obj.ProvisioningProxy.Set(&val)` — set the value
- `obj.ProvisioningProxy.Unset()` — clear the value
### RegionCode (Nullable)

Use the Nullable wrapper methods:
- `obj.RegionCode.IsSet()` — check if set
- `obj.RegionCode.Get()` — get the inner value (returns pointer)
- `obj.RegionCode.Set(&val)` — set the value
- `obj.RegionCode.Unset()` — clear the value
### ImagePath (Nullable)

Use the Nullable wrapper methods:
- `obj.ImagePath.IsSet()` — check if set
- `obj.ImagePath.Get()` — get the inner value (returns pointer)
- `obj.ImagePath.Set(&val)` — set the value
- `obj.ImagePath.Unset()` — clear the value
### DarkImagePath (Nullable)

Use the Nullable wrapper methods:
- `obj.DarkImagePath.IsSet()` — check if set
- `obj.DarkImagePath.Get()` — get the inner value (returns pointer)
- `obj.DarkImagePath.Set(&val)` — set the value
- `obj.DarkImagePath.Unset()` — clear the value
### LastSync (Nullable)

Use the Nullable wrapper methods:
- `obj.LastSync.IsSet()` — check if set
- `obj.LastSync.Get()` — get the inner value (returns pointer)
- `obj.LastSync.Set(&val)` — set the value
- `obj.LastSync.Unset()` — clear the value
### LastSyncDuration (Nullable)

Use the Nullable wrapper methods:
- `obj.LastSyncDuration.IsSet()` — check if set
- `obj.LastSyncDuration.Get()` — get the inner value (returns pointer)
- `obj.LastSyncDuration.Set(&val)` — set the value
- `obj.LastSyncDuration.Unset()` — clear the value
### NextRunDate (Nullable)

Use the Nullable wrapper methods:
- `obj.NextRunDate.IsSet()` — check if set
- `obj.NextRunDate.Get()` — get the inner value (returns pointer)
- `obj.NextRunDate.Set(&val)` — set the value
- `obj.NextRunDate.Unset()` — clear the value

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


