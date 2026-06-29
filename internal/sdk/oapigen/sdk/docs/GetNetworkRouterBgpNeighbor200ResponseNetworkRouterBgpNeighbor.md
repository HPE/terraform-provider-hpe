# GetNetworkRouterBgpNeighbor200ResponseNetworkRouterBgpNeighbor

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | Pointer to **int64** |  | [optional] 
**IpAddress** | Pointer to **string** |  | [optional] 
**Description** | Pointer to **NullableString** |  | [optional] 
**ForwardingAddress** | Pointer to **NullableString** |  | [optional] 
**ProtocolAddress** | Pointer to **NullableString** |  | [optional] 
**RemoteAs** | Pointer to **string** |  | [optional] 
**Weight** | Pointer to **int64** |  | [optional] 
**KeepAlive** | Pointer to **int64** |  | [optional] 
**HoldDown** | Pointer to **int64** |  | [optional] 
**Password** | Pointer to **NullableString** |  | [optional] 
**PasswordHash** | Pointer to **NullableString** |  | [optional] 
**RouteFilteringType** | Pointer to **NullableString** |  | [optional] 
**RouteFilteringIn** | Pointer to **NullableString** |  | [optional] 
**RouteFilteringOut** | Pointer to **NullableString** |  | [optional] 
**BfdEnabled** | Pointer to **NullableBool** |  | [optional] 
**BfdInterval** | Pointer to **NullableInt64** |  | [optional] 
**BfdMultiple** | Pointer to **NullableInt64** |  | [optional] 
**AllowAsIn** | Pointer to **NullableBool** |  | [optional] 
**HopLimit** | Pointer to **NullableInt64** |  | [optional] 
**RestartMode** | Pointer to **NullableString** |  | [optional] 
**ProviderId** | Pointer to **NullableString** |  | [optional] 
**SyncSource** | Pointer to **NullableString** |  | [optional] 
**InternalId** | Pointer to **NullableString** |  | [optional] 
**ExternalId** | Pointer to **NullableString** |  | [optional] 
**RefType** | Pointer to **NullableString** |  | [optional] 
**RefId** | Pointer to **NullableInt64** |  | [optional] 
**Config** | Pointer to [**GetNetworkRouterBgpNeighbor200ResponseNetworkRouterBgpNeighborConfig**](GetNetworkRouterBgpNeighbor200ResponseNetworkRouterBgpNeighborConfig.md) |  | [optional] 
**DateCreated** | Pointer to **time.Time** |  | [optional] 
**LastUpdated** | Pointer to **time.Time** |  | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &GetNetworkRouterBgpNeighbor200ResponseNetworkRouterBgpNeighbor{
    // Set fields directly
}
```

### Description (Nullable)

Use the Nullable wrapper methods:
- `obj.Description.IsSet()` — check if set
- `obj.Description.Get()` — get the inner value (returns pointer)
- `obj.Description.Set(&val)` — set the value
- `obj.Description.Unset()` — clear the value
### ForwardingAddress (Nullable)

Use the Nullable wrapper methods:
- `obj.ForwardingAddress.IsSet()` — check if set
- `obj.ForwardingAddress.Get()` — get the inner value (returns pointer)
- `obj.ForwardingAddress.Set(&val)` — set the value
- `obj.ForwardingAddress.Unset()` — clear the value
### ProtocolAddress (Nullable)

Use the Nullable wrapper methods:
- `obj.ProtocolAddress.IsSet()` — check if set
- `obj.ProtocolAddress.Get()` — get the inner value (returns pointer)
- `obj.ProtocolAddress.Set(&val)` — set the value
- `obj.ProtocolAddress.Unset()` — clear the value
### Password (Nullable)

Use the Nullable wrapper methods:
- `obj.Password.IsSet()` — check if set
- `obj.Password.Get()` — get the inner value (returns pointer)
- `obj.Password.Set(&val)` — set the value
- `obj.Password.Unset()` — clear the value
### PasswordHash (Nullable)

Use the Nullable wrapper methods:
- `obj.PasswordHash.IsSet()` — check if set
- `obj.PasswordHash.Get()` — get the inner value (returns pointer)
- `obj.PasswordHash.Set(&val)` — set the value
- `obj.PasswordHash.Unset()` — clear the value
### RouteFilteringType (Nullable)

Use the Nullable wrapper methods:
- `obj.RouteFilteringType.IsSet()` — check if set
- `obj.RouteFilteringType.Get()` — get the inner value (returns pointer)
- `obj.RouteFilteringType.Set(&val)` — set the value
- `obj.RouteFilteringType.Unset()` — clear the value
### RouteFilteringIn (Nullable)

Use the Nullable wrapper methods:
- `obj.RouteFilteringIn.IsSet()` — check if set
- `obj.RouteFilteringIn.Get()` — get the inner value (returns pointer)
- `obj.RouteFilteringIn.Set(&val)` — set the value
- `obj.RouteFilteringIn.Unset()` — clear the value
### RouteFilteringOut (Nullable)

Use the Nullable wrapper methods:
- `obj.RouteFilteringOut.IsSet()` — check if set
- `obj.RouteFilteringOut.Get()` — get the inner value (returns pointer)
- `obj.RouteFilteringOut.Set(&val)` — set the value
- `obj.RouteFilteringOut.Unset()` — clear the value
### BfdEnabled (Nullable)

Use the Nullable wrapper methods:
- `obj.BfdEnabled.IsSet()` — check if set
- `obj.BfdEnabled.Get()` — get the inner value (returns pointer)
- `obj.BfdEnabled.Set(&val)` — set the value
- `obj.BfdEnabled.Unset()` — clear the value
### BfdInterval (Nullable)

Use the Nullable wrapper methods:
- `obj.BfdInterval.IsSet()` — check if set
- `obj.BfdInterval.Get()` — get the inner value (returns pointer)
- `obj.BfdInterval.Set(&val)` — set the value
- `obj.BfdInterval.Unset()` — clear the value
### BfdMultiple (Nullable)

Use the Nullable wrapper methods:
- `obj.BfdMultiple.IsSet()` — check if set
- `obj.BfdMultiple.Get()` — get the inner value (returns pointer)
- `obj.BfdMultiple.Set(&val)` — set the value
- `obj.BfdMultiple.Unset()` — clear the value
### AllowAsIn (Nullable)

Use the Nullable wrapper methods:
- `obj.AllowAsIn.IsSet()` — check if set
- `obj.AllowAsIn.Get()` — get the inner value (returns pointer)
- `obj.AllowAsIn.Set(&val)` — set the value
- `obj.AllowAsIn.Unset()` — clear the value
### HopLimit (Nullable)

Use the Nullable wrapper methods:
- `obj.HopLimit.IsSet()` — check if set
- `obj.HopLimit.Get()` — get the inner value (returns pointer)
- `obj.HopLimit.Set(&val)` — set the value
- `obj.HopLimit.Unset()` — clear the value
### RestartMode (Nullable)

Use the Nullable wrapper methods:
- `obj.RestartMode.IsSet()` — check if set
- `obj.RestartMode.Get()` — get the inner value (returns pointer)
- `obj.RestartMode.Set(&val)` — set the value
- `obj.RestartMode.Unset()` — clear the value
### ProviderId (Nullable)

Use the Nullable wrapper methods:
- `obj.ProviderId.IsSet()` — check if set
- `obj.ProviderId.Get()` — get the inner value (returns pointer)
- `obj.ProviderId.Set(&val)` — set the value
- `obj.ProviderId.Unset()` — clear the value
### SyncSource (Nullable)

Use the Nullable wrapper methods:
- `obj.SyncSource.IsSet()` — check if set
- `obj.SyncSource.Get()` — get the inner value (returns pointer)
- `obj.SyncSource.Set(&val)` — set the value
- `obj.SyncSource.Unset()` — clear the value
### InternalId (Nullable)

Use the Nullable wrapper methods:
- `obj.InternalId.IsSet()` — check if set
- `obj.InternalId.Get()` — get the inner value (returns pointer)
- `obj.InternalId.Set(&val)` — set the value
- `obj.InternalId.Unset()` — clear the value
### ExternalId (Nullable)

Use the Nullable wrapper methods:
- `obj.ExternalId.IsSet()` — check if set
- `obj.ExternalId.Get()` — get the inner value (returns pointer)
- `obj.ExternalId.Set(&val)` — set the value
- `obj.ExternalId.Unset()` — clear the value
### RefType (Nullable)

Use the Nullable wrapper methods:
- `obj.RefType.IsSet()` — check if set
- `obj.RefType.Get()` — get the inner value (returns pointer)
- `obj.RefType.Set(&val)` — set the value
- `obj.RefType.Unset()` — clear the value
### RefId (Nullable)

Use the Nullable wrapper methods:
- `obj.RefId.IsSet()` — check if set
- `obj.RefId.Get()` — get the inner value (returns pointer)
- `obj.RefId.Set(&val)` — set the value
- `obj.RefId.Unset()` — clear the value

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


