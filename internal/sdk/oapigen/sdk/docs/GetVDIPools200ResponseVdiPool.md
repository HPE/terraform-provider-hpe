# GetVDIPools200ResponseVdiPool

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | Pointer to **int64** |  | [optional] 
**Name** | Pointer to **string** |  | [optional] 
**Description** | Pointer to **NullableString** |  | [optional] 
**MinIdle** | Pointer to **int64** |  | [optional] 
**MaxIdle** | Pointer to **int64** |  | [optional] 
**InitialPoolSize** | Pointer to **int64** |  | [optional] 
**MaxPoolSize** | Pointer to **int64** |  | [optional] 
**AllocationTimeoutMinutes** | Pointer to **int64** |  | [optional] 
**PersistentUser** | Pointer to **NullableBool** |  | [optional] 
**Recyclable** | Pointer to **NullableBool** |  | [optional] 
**Enabled** | Pointer to **bool** |  | [optional] 
**AutoCreateLocalUserOnReservation** | Pointer to **bool** |  | [optional] 
**AllowHypervisorConsole** | Pointer to **NullableBool** |  | [optional] 
**AllowCopy** | Pointer to **NullableBool** |  | [optional] 
**AllowPrinter** | Pointer to **NullableBool** |  | [optional] 
**AllowFileshare** | Pointer to **NullableBool** |  | [optional] 
**GuestConsoleJumpHost** | Pointer to **NullableString** |  | [optional] 
**GuestConsoleJumpPort** | Pointer to **NullableString** |  | [optional] 
**GuestConsoleJumpUsername** | Pointer to **NullableString** |  | [optional] 
**GuestConsoleJumpPassword** | Pointer to **NullableString** |  | [optional] 
**GuestConsoleJumpKeypair** | Pointer to **NullableString** |  | [optional] 
**Gateway** | Pointer to [**GetVDIPools200ResponseVdiPoolGateway**](GetVDIPools200ResponseVdiPoolGateway.md) |  | [optional] 
**IconPath** | Pointer to **string** |  | [optional] 
**Logo** | Pointer to **string** |  | [optional] 
**Apps** | Pointer to [**[]GetVDIPools200ResponseVdiPoolAppsInner**](GetVDIPools200ResponseVdiPoolAppsInner.md) |  | [optional] 
**Owner** | Pointer to [**GetVDIPools200ResponseVdiPoolOwner**](GetVDIPools200ResponseVdiPoolOwner.md) |  | [optional] 
**Config** | Pointer to [**GetVDIPools200ResponseVdiPoolConfig**](GetVDIPools200ResponseVdiPoolConfig.md) |  | [optional] 
**Group** | Pointer to [**GetVDIPools200ResponseVdiPoolGroup**](GetVDIPools200ResponseVdiPoolGroup.md) |  | [optional] 
**Cloud** | Pointer to [**GetVDIPools200ResponseVdiPoolCloud**](GetVDIPools200ResponseVdiPoolCloud.md) |  | [optional] 
**UsedCount** | Pointer to **int64** |  | [optional] 
**ReservedCount** | Pointer to **int64** |  | [optional] 
**PreparingCount** | Pointer to **int64** |  | [optional] 
**IdleCount** | Pointer to **int64** |  | [optional] 
**Status** | Pointer to **string** |  | [optional] 
**DateCreated** | Pointer to **time.Time** |  | [optional] 
**LastUpdated** | Pointer to **time.Time** |  | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &GetVDIPools200ResponseVdiPool{
    // Set fields directly
}
```

### Description (Nullable)

Use the Nullable wrapper methods:
- `obj.Description.IsSet()` — check if set
- `obj.Description.Get()` — get the inner value (returns pointer)
- `obj.Description.Set(&val)` — set the value
- `obj.Description.Unset()` — clear the value
### PersistentUser (Nullable)

Use the Nullable wrapper methods:
- `obj.PersistentUser.IsSet()` — check if set
- `obj.PersistentUser.Get()` — get the inner value (returns pointer)
- `obj.PersistentUser.Set(&val)` — set the value
- `obj.PersistentUser.Unset()` — clear the value
### Recyclable (Nullable)

Use the Nullable wrapper methods:
- `obj.Recyclable.IsSet()` — check if set
- `obj.Recyclable.Get()` — get the inner value (returns pointer)
- `obj.Recyclable.Set(&val)` — set the value
- `obj.Recyclable.Unset()` — clear the value
### AllowHypervisorConsole (Nullable)

Use the Nullable wrapper methods:
- `obj.AllowHypervisorConsole.IsSet()` — check if set
- `obj.AllowHypervisorConsole.Get()` — get the inner value (returns pointer)
- `obj.AllowHypervisorConsole.Set(&val)` — set the value
- `obj.AllowHypervisorConsole.Unset()` — clear the value
### AllowCopy (Nullable)

Use the Nullable wrapper methods:
- `obj.AllowCopy.IsSet()` — check if set
- `obj.AllowCopy.Get()` — get the inner value (returns pointer)
- `obj.AllowCopy.Set(&val)` — set the value
- `obj.AllowCopy.Unset()` — clear the value
### AllowPrinter (Nullable)

Use the Nullable wrapper methods:
- `obj.AllowPrinter.IsSet()` — check if set
- `obj.AllowPrinter.Get()` — get the inner value (returns pointer)
- `obj.AllowPrinter.Set(&val)` — set the value
- `obj.AllowPrinter.Unset()` — clear the value
### AllowFileshare (Nullable)

Use the Nullable wrapper methods:
- `obj.AllowFileshare.IsSet()` — check if set
- `obj.AllowFileshare.Get()` — get the inner value (returns pointer)
- `obj.AllowFileshare.Set(&val)` — set the value
- `obj.AllowFileshare.Unset()` — clear the value
### GuestConsoleJumpHost (Nullable)

Use the Nullable wrapper methods:
- `obj.GuestConsoleJumpHost.IsSet()` — check if set
- `obj.GuestConsoleJumpHost.Get()` — get the inner value (returns pointer)
- `obj.GuestConsoleJumpHost.Set(&val)` — set the value
- `obj.GuestConsoleJumpHost.Unset()` — clear the value
### GuestConsoleJumpPort (Nullable)

Use the Nullable wrapper methods:
- `obj.GuestConsoleJumpPort.IsSet()` — check if set
- `obj.GuestConsoleJumpPort.Get()` — get the inner value (returns pointer)
- `obj.GuestConsoleJumpPort.Set(&val)` — set the value
- `obj.GuestConsoleJumpPort.Unset()` — clear the value
### GuestConsoleJumpUsername (Nullable)

Use the Nullable wrapper methods:
- `obj.GuestConsoleJumpUsername.IsSet()` — check if set
- `obj.GuestConsoleJumpUsername.Get()` — get the inner value (returns pointer)
- `obj.GuestConsoleJumpUsername.Set(&val)` — set the value
- `obj.GuestConsoleJumpUsername.Unset()` — clear the value
### GuestConsoleJumpPassword (Nullable)

Use the Nullable wrapper methods:
- `obj.GuestConsoleJumpPassword.IsSet()` — check if set
- `obj.GuestConsoleJumpPassword.Get()` — get the inner value (returns pointer)
- `obj.GuestConsoleJumpPassword.Set(&val)` — set the value
- `obj.GuestConsoleJumpPassword.Unset()` — clear the value
### GuestConsoleJumpKeypair (Nullable)

Use the Nullable wrapper methods:
- `obj.GuestConsoleJumpKeypair.IsSet()` — check if set
- `obj.GuestConsoleJumpKeypair.Get()` — get the inner value (returns pointer)
- `obj.GuestConsoleJumpKeypair.Set(&val)` — set the value
- `obj.GuestConsoleJumpKeypair.Unset()` — clear the value

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


